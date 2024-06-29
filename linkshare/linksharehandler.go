package linkshare

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ShareLinkHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Contains(requestData.Data, "\n") {
		http.Error(w, "Multi-line input is not allowed", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile("./file.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return
	}
	fileSize := fileInfo.Size()

	buffer := make([]byte, 1)
	var lastLine []byte
	var pos = fileSize - 1

	for pos >= 0 {
		_, err := f.ReadAt(buffer, pos)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		char := buffer[0]

		if char == '\n' && len(lastLine) > 0 {
			break
		}
		lastLine = append([]byte{char}, lastLine...)
		pos--
	}

	var indexNo string
	var index int
	for _, r := range string(lastLine) {
		if r >= '0' && r <= '9' {
			indexNo += string(r)
		} else {
			break
		}
	}
	if indexNo == "" {
		indexNo = "0"
	} else {
		index, err = strconv.Atoi(indexNo)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return
		}
		index += 1
	}

	_, err = f.WriteString(fmt.Sprintf("%d %s\n", index, requestData.Data))

	if err != nil {
		fmt.Println("Error appending record:", err)
	}

	w.WriteHeader(http.StatusOK)
}

func ListLinksHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile("./file.txt")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var dataList []Data
	for _, line := range lines {
		if line != "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				data := Data{
					Key:   parts[0],
					Value: parts[1],
				}
				dataList = append(dataList, data)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(dataList); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
