package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// GetLocalIP i.e. 192.168.X.Y ~ your local private IP
func GetLocalIP() net.IP {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4
		}
	}
	return nil
}

func isDirectoryExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Directory `%s` does not exists\n", path)
		return false
	}
	return true
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

//ParseArgs ...
func ParseArgs() string {
	homeDir, _ := os.UserHomeDir()
	currentDir, _ := os.Getwd()
	if len(os.Args) > 1 {
		opt := os.Args[1]
		if strings.HasPrefix(opt, "-h") || strings.HasPrefix(opt, "--h") {
			showUsage()
			os.Exit(0)
		} else {
			if opt == "~" {
				return homeDir
			} else if strings.HasPrefix(opt, "~/") {
				dirpath := filepath.Join(homeDir, opt[2:])
				if isDirectoryExists(dirpath) {
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

//HTMLReplacer ...
var HTMLReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)


//GetHumanReadableSize ...
func GetHumanReadableSize(f os.FileInfo) string{
	if f.IsDir() {
		return "--"
	}
	bytes := f.Size()
	mb := float32(bytes)/(1024.0*1024.0)
	return fmt.Sprintf("%.2f MB",mb)
}