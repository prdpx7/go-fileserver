# go-fileserver
> A simple HTTP server to share files over WiFi via QRCode

# Installation
* You can download compressed version from [releases](https://github.com/prdpx7/go-fileserver/releases)
    ```
    wget https://github.com/prdpx7/go-fileserver/releases/download/v0.1/fs-server-2020.07.25.tar.gz
    tar -xzf fs-server-2020.07.25.tar.gz
	chmod +x fs-server && sudo cp fs-server /usr/local/bin/fs-server
    ```
* Or download the binary directly
	```
	wget https://github.com/prdpx7/go-fileserver/releases/download/v0.1/fs-server
	chmod +x fs-server && sudo cp fs-server /usr/local/bin/fs-server
	```

* Or you can clone from GitHub and build the binary yourself
	```
	git clone https://github.com/prdpx7/go-fileserver --depth=1
	cd go-fileserver/fs-server
	# requires go 1.14
	go build
	# make binary executable
	chmod +x ./fs-server
	# may require root permission
	cp fs-server /usr/local/bin/fs-server
	```
# Usage
```
fs-server - A simple HTTP Server to share files on a network.
Usage: fs-server [OPTIONS] <dir-path>
Options:
	-h | --help - show this message
Example:
fs-server - serve files from current directory
fs-server /home/user/documents/ - serve files from given directory
```
# Demo

### Step 1 - Run in terminal
<img src ="https://i.imgur.com/ywUaM08.gif" width=800 height=450>

### Step 2 - Scan QRCode on Phone
<img src="https://i.imgur.com/pIlaFol.gif" width=350 height=700>

# Inspiration
* Inspired from [http-server](https://github.com/http-party/http-server) project

# License
* MIT