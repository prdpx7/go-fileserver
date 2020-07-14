package fileserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	utils "github.com/prdpx7/go-fileserver/utils"
)


type customFileHandler struct {
	root http.FileSystem
}

func CustomFileServer(root http.FileSystem) http.Handler {
	return &customFileHandler{root}
}

func (cf *customFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	ServeFile(w, r, cf.root, path.Clean(upath), true)
}

// ServeFile ...
func ServeFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool) {
	f, err := fs.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	if d.IsDir() {
		ListDirectory(w, r, f, "index")
		return
	}
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

//ListDirectory render directory content in templateName.html
func ListDirectory(w http.ResponseWriter, r *http.Request, f http.File, templateName string) {
	RootDir, err := f.Stat()
	if err != nil {
		panic(err)
	}
	var dirContents DirectoryContent
	dirContents.DirName = RootDir.Name()
	dirContents.Files = make([]FileContent, 0)
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Printf("http: error reading directory: %v", err)
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, d := range dirs {
		name := d.Name()
		fileExtension := "page"
		if d.IsDir() {
			name += "/"
			fileExtension = "directory"
		} else if len(filepath.Ext(name)) > 1 {
			fileExtension = filepath.Ext(name)[1:]
		}

		url := url.URL{Path: name}
		fileContent := FileContent{Name: name, Size: utils.GetHumanReadableSize(d), URL: url, Extension: fileExtension}
		dirContents.Files = append(dirContents.Files, fileContent)
	}
	dirContents.IPAddr = r.Host
	renderTemplate(w, templateName, dirContents)
}

//DirectoryContent to be used in rendering Index Page
type DirectoryContent struct {
	DirName string
	Files   []FileContent
	IPAddr  string
}

//FileContent ...
type FileContent struct {
	Name      string
	Size      string
	URL       url.URL
	Extension string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		fmt.Println("Error in parsing template ",err)
		panic(err)
	}
	t.Execute(w, data)
}

//RequestLogger
func RequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}


