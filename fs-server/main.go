package main

import (
	"fmt"
	"net/http"

	fileserver "github.com/prdpx7/go-fileserver"
	utils "github.com/prdpx7/go-fileserver/utils"
	qrcode "github.com/skip2/go-qrcode"
)

//
const (
	PORT = 8000
)

func main() {
	dirpath := utils.ParseArgs()
	localIP := utils.GetLocalIP()
	fmt.Printf("Currently Serving `%s` on:\n", dirpath)
	fmt.Printf("http://localhost:%d\n", PORT)

	if localIP != nil {
		url := fmt.Sprintf("http://%s:%d\n", localIP, PORT)
		fmt.Printf(url)
		qr, _ := qrcode.New(url, qrcode.High)
		// qr.DisableBorder = true
		fmt.Println(qr.ToSmallString(false))
	}

	fs := fileserver.CustomFileServer(http.Dir(dirpath))
	portNumber := fmt.Sprintf(":%d", PORT)
	err := http.ListenAndServe(portNumber, fileserver.RequestLogger(fs))
	if err != nil {
		fmt.Println(err)
	}
}