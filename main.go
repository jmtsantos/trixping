package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/matrix-org/gomatrix"
)

const (
	// Description ..
	Description = "A simple command line shortcut to send messages trough the Matrix.org protocol"
)

var config Config

// Config matrix client configuration
type Config struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Server   string `json:"server"`
	Room     string `json:"room"`
}

// RoomMsg message structure
type RoomMsg struct {
	MsgType       string `json:"msgtype"`
	Format        string `json:"format"`
	Body          string `json:"body"`
	FormattedBody string `json:"formatted_body"`
}

// Generic command structure
var cmdArg struct {
	ConfigPath string `help:"Full path to the config file. Default paths are:\n  ~/.config/trixping.json\n  /etc/trixping.json" short:"c" type:"path"`
	Message    string `help:"Message to be sent. If empty STDIN will be used as input" short:"m"`
	Sender     string `help:"Set the full name of the sender." short:"F" sep:"none"`
}

// A sendmail compatible emulation `/usr/sbin/sendmail -i -FCronDaemon -B8BITMIME -oem root`
var sendMailArg struct {
	Dots   bool   `help:"Ignore dots alone on lines by themselves in incoming messages. This should be set if you are reading data from a file." short:"i" default:"false" required:"false"`
	Type   string `help:"Set the body type to type. Current legal values are 7BIT or 8BITMIME." short:"B" default:"false" sep:"none" optional`
	Oe     string `help:"Unknown" short:"o" default:"false" sep:"none" optional`
	Sender string `help:"Set the full name of the sender." short:"F" sep:"none"`

	DestinationUser []string `name:"dest-user" help:"Destination user" type:"path" arg optional`
}

func main() {
	var (
		cli        *gomatrix.Client
		hostname   string
		configJSON []byte
		err        error
	)

	// Command line parse
	if strings.HasSuffix(os.Args[0], "/sendmail") {
		kong.Parse(&sendMailArg,
			kong.Name("trixping"),
			kong.Description(Description),
			kong.UsageOnError(),
			kong.ConfigureHelp(kong.HelpOptions{
				Compact: true,
				Summary: true,
			}))

	} else {
		kong.Parse(&cmdArg,
			kong.Name("trixping"),
			kong.Description(Description),
			kong.UsageOnError(),
			kong.ConfigureHelp(kong.HelpOptions{
				Compact: true,
				Summary: true,
			}))
	}

	// Get the configuration file
	if cmdArg.ConfigPath == "" {
		homePath := fmt.Sprintf("%s/.config/trixping.json", os.Getenv("HOME"))
		if _, err = os.Stat(homePath); !os.IsNotExist(err) {
			cmdArg.ConfigPath = homePath
		} else if _, err = os.Stat("/etc/trixping.json"); !os.IsNotExist(err) {
			cmdArg.ConfigPath = "/etc/trixping.json"
		} else {
			fmt.Println("error configuration file does not exist")
			os.Exit(1)
		}
	}

	if configJSON, err = ioutil.ReadFile(cmdArg.ConfigPath); err != nil {
		fmt.Println("error reading file from disk", err)
		os.Exit(1)
	}

	if err = json.Unmarshal(configJSON, &config); err != nil {
		fmt.Println("error unmarshalling json", err)
		os.Exit(1)
	}

	if cli, err = gomatrix.NewClient(config.Server, config.Username, config.Token); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Message composition starts here
	var messageHeader string
	if hostname, err = os.Hostname(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	messageHeader = fmt.Sprintf("<h4>host: %s</h4>", hostname)

	// Set the sender field
	if cmdArg.Sender != "" {
		messageHeader = fmt.Sprintf("%s<h4>sender: %s</h4>", messageHeader, cmdArg.Sender)
	} else {
		messageHeader = fmt.Sprintf("%s<h4>sender: %s</h4>", messageHeader, "undefined")
	}

	if cmdArg.Message == "" {

		cmdArg.Message = cmdArg.Message + messageHeader

		cmdArg.Message = cmdArg.Message + "<code>"
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			cmdArg.Message = cmdArg.Message + strings.ReplaceAll(strings.ReplaceAll(scanner.Text(), "<", "&lt;"), ">", "&gt;") + "<br/>"
		}
		cmdArg.Message = cmdArg.Message + "</code>"
	} else {
		cmdArg.Message = fmt.Sprintf("%s<code>%s</code>", messageHeader, cmdArg.Message)
	}

	msg := RoomMsg{
		MsgType:       "m.text",
		Format:        "org.matrix.custom.html",
		Body:          cmdArg.Message,
		FormattedBody: cmdArg.Message,
	}

	if _, err = cli.SendMessageEvent(config.Room, "m.room.message", msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
