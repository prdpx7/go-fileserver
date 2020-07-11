package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//
const (
	PORT = 8000
)
func getLocalIP() net.IP{
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4
		}
	}
	return nil
}

func showUsage() {
	helpMessage := `fs-server - A simple HTTP Server to share files on a network.
Usage: fs-server [OPTIONS] <dir-path>
Options:
	-h | --help - show this message
Example:
fs-server - serve files from current directory
fs-server /home/user/documents/ - serve files from given directory
`
	fmt.Println(helpMessage)
}

func isDirectoryExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err){
		fmt.Printf("Directory `%s` does not exists\n", path)
		return false
	}
	return true
}

func parseArgs() string {
	homeDir, _ := os.UserHomeDir()
	currentDir, _ := os.Getwd()
	if len(os.Args) > 1 {
		opt := os.Args[1]
		if strings.HasPrefix(opt, "-h") || strings.HasPrefix(opt, "--h"){
			showUsage()
			os.Exit(0)
		} else {
			if opt == "~" {
				return homeDir
			} else if strings.HasPrefix(opt, "~/") {
				dirpath := filepath.Join(homeDir, opt[2:])
				if isDirectoryExists(dirpath){
					return dirpath
				}
			} else if isDirectoryExists(opt) {
				return opt
			}
		}
	}
	//fallback to current directory
	return currentDir
}

func main() {
	dirpath := parseArgs()
	localIP := getLocalIP()
	fmt.Printf("Currently Serving `%s` on:\n", dirpath)
	fmt.Printf("http://localhost:%d\n", PORT)
	if localIP != nil {
		fmt.Printf("http://%s:%d\n", localIP, PORT)
	}
	fileServer := http.FileServer(http.Dir(dirpath))
	err	:=	http.ListenAndServe(":8000",	fileServer)
	if err != nil {
		fmt.Println(err)
	}

}

