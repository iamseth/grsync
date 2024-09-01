package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Directory struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

type PhotoListReponse struct {
	StatusCode   int         `json:"errCode"`
	ErrorMessage string      `json:"errMsg"`
	Directories  []Directory `json:"dirs"`
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DownloadPhoto(src, dest string) error {
	if FileExists(dest) {
		return nil
	}

	// Ensure the directory exists
	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create the destination file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download the file
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the file to the destination
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func ListPhotos(endpoint string) PhotoListReponse {
	resp, err := http.Get(endpoint + "/photos")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var photoList PhotoListReponse
	if err := json.NewDecoder(resp.Body).Decode(&photoList); err != nil {
		log.Fatal(err)
	}

	if photoList.StatusCode != 200 {
		log.Fatalf("Error: %s", photoList.ErrorMessage)
	}

	return photoList
}

func DownloadAllPhotos(endpoint string, photoList PhotoListReponse, path string) {

	var wg sync.WaitGroup
	for _, dir := range photoList.Directories {
		for _, file := range dir.Files {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				src := endpoint + "/photos/" + dir.Name + "/" + file
				dest := path + "/" + dir.Name + "/" + file
				log.Printf("Downloading %s -> %s", src, dest)
				err := DownloadPhoto(src, dest)
				if err != nil {
					log.Printf("failed to download: %v", err)
				}

			}(file)
		}
	}
	wg.Wait()
	log.Println("All files downloaded")
}

func main() {
	// Get the destination directory from a command line argument
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <destination>", os.Args[0])
	}

	endpoint := "http://192.168.0.1/v1"
	photoList := ListPhotos(endpoint)
	DownloadAllPhotos(endpoint, photoList, os.Args[1])

}
