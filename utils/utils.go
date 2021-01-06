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
    addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
    /*
    // Not properly working where multiple network interfaces are enabled

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4
		}
	}
	*/
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
	helpMessage := `fs-server - A simple HTTP Server to share files on a network via QRCode.
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

//DirListTemplateHTML to be used as index.html for redering DirectoryList
var DirListTemplateHTML = `
<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
	<title>Index of /</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/file-icon-vectors@1.0.0/dist/file-icon-classic.min.css" />
	<style type="text/css">td.icon-parent { height: 16px; width: 16px; }

	td.file-size { text-align: right; padding-left: 1em; white-space:nowrap;}
	td.display-name { padding-left: 1.5em; }

	</style>	
	</head>
  <body>
  <h1>Index of {{.DirName}}</h1>
  <table>
	<!-- <tr><th>Type</th><th>Size</th><th>Name</th></tr> -->
  
	  {{range .Files}}
	  <tr>
		  <td class="icon-parent">
			  <i class="fiv-cla fiv-icon-{{.Extension}}"></i>
		  </td>
		  <td class="file-size"><code>{{.Size}}</code></td>
		  <td class="display-name"><a href="{{.URL}}">{{.Name}}</a></td>
	  </tr>
	  {{end}}
  </table>
  
  <br><address style="font-size: 1.2em;"><a href="https://github.com/prdpx7/go-fileserver"><strong>fs-server</strong></a> running @ {{.IPAddr}}</address>
  </body></html>
`
