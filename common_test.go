package httpfile

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	testServer = httptest.NewServer(http.HandlerFunc(uploadHandler))

	retCode := m.Run()
	testServer.Close()
	os.Exit(retCode)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/upload" {
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("testheader", r.Header.Get("testheader"))

	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := r.MultipartForm
	for key, val := range m.Value {
		w.Header().Add(key, val[0])
	}
	for _, val := range m.File {
		for _, file := range val {
			err := handlerFile(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func handlerFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	if err != nil {
		return err
	}
	dst, err := os.Create("testdata/download/" + fileHeader.Filename)
	defer dst.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
