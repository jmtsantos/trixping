# trixping
A simple command line shortcut to send messages trough the Matrix.org protocol

## Usage

```
  -h, --help                  Show context-sensitive help.
  -c, --config-path=STRING    Full path to the config file. Default paths are:\n ~/.config/trixping.json\n /etc/trixping.json
  -m, --message=STRING        Message to be sent. If empty STDIN will be used as input
  -F, --sender=STRING         Set the full name of the sender.
```

A simple use case for this app can be to notify the user everytime someone logs in to a server via SSH. This can be done by editing the `/etc/ssh/sshrc` as so:
```
echo "<h3>ssh login - `cat /etc/hostname`</h3><p>User $USER just logged in `echo $SSH_CONNECTION | cut -d " " -f 1`</p>" | trixping "`cat /etc/hostname` - SSH Login" &
```

Added support to be used as in place replacement for `sendmail`, this way you can edit the simlink to `/usr/bin/sendmail` to point to `trixping` and have your system mail redirected to your matrix server.