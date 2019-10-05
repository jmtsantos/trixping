# trixping
A simple command line shortcut to send messages trough the Matrix.org protocol

## Usage

```
Usage of ./trixping:
  -c string
        Full path to the config file. Default paths are:
          ~/.config/trixping.json
          /etc/trixping.json
  -m string
        Message to be sent. Use "-" to use STDIN as input
```

A simple use case for this app can be to notify the user everytime someone logs in to a server via SSH. This can be done by editing the `/etc/ssh/sshrc` as so:
```
echo "<h3>ssh login - `cat /etc/hostname`</h3><p>User $USER just logged in `echo $SSH_CONNECTION | cut -d " " -f 1`</p>" | trixping "`cat /etc/hostname` - SSH Login" &
```
