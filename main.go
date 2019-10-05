package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/matrix-org/gomatrix"
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

func main() {
	var (
		filePath, message       *string
		cli                     *gomatrix.Client
		configJSON, messageByte []byte
		err                     error
	)

	// Command line parse
	filePath = flag.String("c", "", "Full path to the config file. Default paths are:\n  ~/.config/trixping.json\n  /etc/trixping.json")
	message = flag.String("m", "", "HTML message to be sent. Use \"-\" to use STDIN as input")
	flag.Parse()

	if *message == "" {
		fmt.Println("error message not set")
		os.Exit(1)
	}

	// Get the configuration file
	if *filePath == "" {
		homePath := fmt.Sprintf("%s/.config/trixping.json", os.Getenv("HOME"))
		if _, err = os.Stat(homePath); !os.IsNotExist(err) {
			*filePath = homePath
		} else if _, err = os.Stat("/etc/trixping.json"); !os.IsNotExist(err) {
			*filePath = "/etc/trixping.json"
		} else {
			fmt.Println("error configuration file does not exist")
			os.Exit(1)
		}
	}

	if configJSON, err = ioutil.ReadFile(*filePath); err != nil {
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

	if *message == "-" {
		reader := bufio.NewReader(os.Stdin)
		messageByte, _, err = reader.ReadLine()
		*message = string(messageByte)
	}

	msg := RoomMsg{
		MsgType:       "m.text",
		Format:        "org.matrix.custom.html",
		Body:          *message,
		FormattedBody: *message,
	}

	if _, err = cli.SendMessageEvent(config.Room, "m.room.message", msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
