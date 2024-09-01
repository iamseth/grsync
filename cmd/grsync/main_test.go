package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	if FileExists("main.go") == false {
		t.Error("FileExists() failed")
	}
}

func TestListPhotos(t *testing.T) {
	reponse := PhotoListReponse{
		StatusCode:   200,
		ErrorMessage: "OK",
		Directories: []Directory{
			{
				Name: "dir1",
				Files: []string{
					"file1.jpg",
					"file2.jpg",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		_ = json.NewEncoder(rw).Encode(reponse)

	}))

	defer server.Close()

	resp := ListPhotos(server.URL)
	if resp.StatusCode != 200 {
		t.Error("StatusCode not 200")
	}
	if len(resp.Directories) != 1 {
		t.Error("Directories length not 1")
	}
}

func TestDownloadPhoto(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "image/jpeg")
		rw.WriteHeader(200)
		_, _ = rw.Write([]byte("image data"))
	}))

	defer server.Close()
	defer func() {
		_ = os.Remove("test.jpg")
	}()

	err := DownloadPhoto(server.URL, "test.jpg")
	if err != nil {
		t.Error("DownloadPhoto() failed")
	}

}
