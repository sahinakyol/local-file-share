package main

import (
	"fmt"
	"local-file-share/fileshare"
	"local-file-share/linkshare"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/upload", fileshare.UploadHandler)
	http.HandleFunc("/files", fileshare.ListFilesHandler)
	http.HandleFunc("/download/", fileshare.DownloadHandler)
	http.HandleFunc("/share-link", linkshare.ShareLinkHandler)
	http.HandleFunc("/list-links", linkshare.ListLinksHandler)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
