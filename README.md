# nurgling
simple http server written in golang as a learning project

## current features
- reads index.html at start time from current directory
- responds to any connection with the HTTP 200 OK and the index.html
- can parse HTTO requests, but is not used yet

## planned features
- ability to parse and serve specific GET requests
- command line flags for configuration
- ability to read confoguration file
- some sort of api for web toolkits
- in the far distant future: HTTPS
