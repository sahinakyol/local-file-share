package fileshare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var mu sync.Mutex

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("./uploads/*")
	for i := range files {
		files[i] = filepath.Base(files[i])
	}
	json.NewEncoder(w).Encode(files)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	filePath := "./uploads/" + fileName

	http.ServeFile(w, r, filePath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("fileName")
	chunkIndex, _ := strconv.Atoi(r.FormValue("chunkIndex"))
	totalChunks, _ := strconv.Atoi(r.FormValue("totalChunks"))

	mu.Lock()
	defer mu.Unlock()

	file, _, err := r.FormFile("chunk")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	os.MkdirAll("./uploads/temp/"+fileName, os.ModePerm)
	chunkPath := fmt.Sprintf("./uploads/temp/%s/%d.chunk", fileName, chunkIndex)

	out, err := os.Create(chunkPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	io.Copy(out, file)

	if chunkIndex == totalChunks-1 {
		combineChunks(fileName, totalChunks)
	}
}

func combineChunks(fileName string, totalChunks int) {
	combinedFilePath := "./uploads/" + fileName
	combinedFile, _ := os.Create(combinedFilePath)
	defer combinedFile.Close()

	for i := 0; i < totalChunks; i++ {
		chunkPath := fmt.Sprintf("./uploads/temp/%s/%d.chunk", fileName, i)
		chunkFile, _ := os.Open(chunkPath)
		io.Copy(combinedFile, chunkFile)
		chunkFile.Close()
		os.Remove(chunkPath)
	}

	os.RemoveAll("./uploads/temp/" + fileName)
}
