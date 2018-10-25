# nurgling
simple http server written in golang as a learning project

## Current features
- parses and serve HTTP and HTTPS request
- parse configuration files 
- supports CGI
- privilege dropping after opening the sockets

## Planned features
- command line flags for configuration (currently only `-f config_file`)

## Usage
### Compilation
- create a go working directory
`mkdir -p ~/go/src`
`cd ~/go/src`
- copy repository
`git clone https://github.com/karlyan17/nurgling.git`
- adjust config file for your needs
- compile nurgling webserver
`go install nurgling`
- run the webserver
` ~/go/bin/nurgling -f ~/go/src/nurgling.config`
