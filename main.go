package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		files := []string{"one.zip", "two.zip", "three.zip", "four.zip"}

		for _, fileName := range files {
			go downloadFile(fileName)
		}
		fmt.Fprintf(w, "Downloding started")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func downloadFile(fileName string) {

	time.Sleep(8 * time.Second)

	url := "http://213.136.80.59:8001/zip_files/" + fileName

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()


	// copy resp body to bytes
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()
	
	bytes, err := out.Write(bodyBytes)
	sbytes:= strconv.FormatUint(uint64(bytes), 10)

	fmt.Println("file", fileName, "Bytes", sbytes)
}

