package httpfile

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	os.MkdirAll("testdata/fileserver", 0777)
	os.MkdirAll("testdata/download", 0777)
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
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		handlerStream(w, r)
		return
	}
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
func handlerStream(w http.ResponseWriter, r *http.Request) error {
	fileName := RandomMD5()
	dst, err := os.Create("testdata/fileserver/" + fileName)
	defer dst.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, r.Body); err != nil {
		return err
	}
	w.Header().Set("FileName", fileName)
	return nil
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
	dst, err := os.Create("testdata/fileserver/" + fileHeader.Filename)
	defer dst.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}

func RandomMD5() string {
	data := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
