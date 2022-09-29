package main

import (
	"archive/zip"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		//files := []string{"one.zip", "two.zip", "three.zip", "four.zip"}

		files := []string{"three.zip"}

		for _, fileName := range files {
			go downloadFile(fileName)
		}
		fmt.Fprintf(w, "Downloding started")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}

func downloadFile(fileName string) {
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
	if err != nil {
		fmt.Println(err)
	}

	sbytes := strconv.FormatUint(uint64(bytes), 10)

	fmt.Println("file", fileName, "Bytes", sbytes)

	erra := unzipSource(fileName, "testivan")
	if erra != nil {
		log.Fatal(err)
	}

}

func unzipSource(source, destination string) error {

	fmt.Println("started", source)
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {

		
		// added by ivan
		if !f.FileInfo().IsDir() {
			if filepath.Ext(f.FileInfo().Name()) != ".css" {
				continue
			}
		}
		// added by ivan EOF

		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	fmt.Println("done", source)

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
