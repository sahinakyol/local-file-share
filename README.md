# Local File and Link Share

Local File and Link Share is a simple web application designed to facilitate the sharing of files and links within your local network without the need for third-party applications such as WhatsApp, AirDrop, Messenger, etc.

## Features

- **File Upload:** Easily upload files from your device and share them within your local network.
- **Link Sharing:** Share links and access them from any device connected to your local network.
- **Chunked File Uploads:** Supports uploading large files in chunks to ensure smooth and reliable file transfers.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) installed on your system

### Installation

1. **Clone the repository:**

```shell
    git clone https://github.com/sahinakyol/local-file-link-share.git
    cd local-file-link-share
    go build -o local-file-link-share
    ./local-file-link-share
    curl "http://$(ifconfig en0 | grep -o 'inet [0-9.]\+' | sed 's/inet //' | tr -d '\n'):8080"
```

## TODO
- Add functionality to delete uploaded files and shared links.

## Personal Use Note
I use it with my old Raspberry Pi, and frankly, it works quite well :)