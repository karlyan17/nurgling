# nurgling
simple http server written in golang as a learning project

## Current features
- parses HTTP request
- on a GET request, returns the requested resource
- constructs valid HTTP header based on requested resource
- parse configuration files 

## Planned features
- command line flags for configuration (currently only `-f config_file`)
- some sort of api for web toolkits
- in the far distant future: HTTPS

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
