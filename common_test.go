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
	"time"
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
	if r.URL.Path != "/file" {
		return
	}
	if r.Method == "POST" {
		w.Header().Set("testheader", r.Header.Get("testheader"))
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
	} else if r.Method == "GET" {
		filename := r.Header.Get("filename")
		if filename == "" {
			filename = r.URL.Query().Get("filename")
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		f, err := os.Open(FileServer(filename))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		http.ServeContent(w, r, filename, time.Now().Add(time.Hour), f)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
func handlerStream(w http.ResponseWriter, r *http.Request) error {
	fileName := RandomMD5()
	dst, err := os.Create(FileServer(fileName))
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
	dst, err := os.Create(FileServer(fileHeader.Filename))
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

func FileServer(filename string) string {
	return "testdata/fileserver/" + filename
}
