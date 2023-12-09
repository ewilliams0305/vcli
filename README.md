[![Go Report Card](https://goreportcard.com/badge/github.com/ewilliams0305/offshoot?style=flat-square)](https://goreportcard.com/report/github.com/ewilliams0305/VC4-CLI)
```go
__      _______ _  _      _____ _      _____ 
\ \    / / ____| || |    / ____| |    |_   _|
 \ \  / / |    | || |_  | |    | |      | |  
  \ \/ /| |    |__   _| | |    | |      | |  
   \  / | |____   | |   | |____| |____ _| |_ 
    \/   \_____|  |_|    \_____|______|_____|
```

# VC4-CLI
Ever found yourself connected to a VC4 appliance troublshooting the OS, working on the file system,
or restarting services, only to find you can't perform actions on the actual VC4 service? Well now you can. 
The VC4 CLI provides full control over the VC4 service from within the linux terminal allowing operators to:


Command line interface to operate a Crestron Virtual Control server application from the Linux Shell
The VC4 CLI leverages the Crestron Virtual control REST API with a loopback IP address 
to provide a localized CLI for the VC4 service. This CLI 
uses the `BubbleTea` TUI framework for navigation and command line workflows. 

# Building 
To compile the cli for your VC4 appliance you will beed to install
the GO sdk. Download the SDK https://go.dev/dl/
 
Once download you can build the /cmd/cli directory as an executable. 
### Windows

### Linux

# Crestron REST API Reference 
https://www.crestron.com/getmedia/29921c49-86df-488c-a63b-ab88620d7175/mg_pg_rest-api-crestron-virtual-control


