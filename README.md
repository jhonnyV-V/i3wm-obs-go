# What is this?
is a tool for people that use i3wm and obs, where you can set a number of excluded workspaces
where a obs source or scene should be disabled

# How to use?
open obs, run the tool once to create the config file and edit the config file created
now you can just run the program without problems

the config file will look like this
```yaml
password: YourSecurePassword #you can get it in Tools -> Websocker server configs
url: localhost:4455 #this is the default port
sourceName: "Screen Capture (XSHM)" #Source Name (Leave the string empty if you want to block the scene)
excludedWorkspaces: ["10", "9", "8", "7"] #Workspaces names (emojis are also part of the names)
isScene: false # set if you want to block the scene
```

# How to install?

### From Releases Page
binaries for linux and mac are [here](https://github.com/jhonnyV-V/i3wm-obs-go/releases)

### With Go
```bash
go install github.com/jhonnyV-V/i3wm-obs-go@latest
```

### Compile it yourself

```bash
git clone git@github.com:jhonnyV-V/i3wm-obs-go && cd i3wm-obs-go
```
```bash
go mod tidy
```
```bash
go build -o i3wm-obs-go 
```
```bash
sudo mv ./i3wm-obs-go/usr/local/bin
```

