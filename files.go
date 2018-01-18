package httpfile

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/mushroomsir/mimetypes"
)

var (
	ErrEmptyTargetURL = errors.New("Empty Target URL")
	ErrEmptyFilePath  = errors.New("Empty File Path")
)

// Files ...
type Files struct {
	client    *http.Client
	targetURL string
	filePath  string
	header    map[string]string
}

// NewReq ...
func NewReq(targetURL string, filePath ...string) *Files {
	fp := ""
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	hf := &Files{
		client:    defaultHTTPClient,
		targetURL: targetURL,
		filePath:  fp,
		header:    make(map[string]string),
	}
	return hf
}

// SetHTTPClient ...
func (h *Files) SetHTTPClient(c *http.Client) *Files {
	if c != nil {
		h.client = c
	}
	return h
}

// SetHeader ...
func (h *Files) SetHeader(k, v string) *Files {
	h.header[k] = v
	return h
}

// SetAuthorization ...
func (h *Files) SetAuthorization(v string) *Files {
	h.header["Authorization"] = v
	return h
}

func (h *Files) checkUpload() *Response {
	res := &Response{}
	if h.targetURL == "" {
		res.err = ErrEmptyTargetURL
		return res
	}
	if h.filePath == "" {
		res.err = ErrEmptyFilePath
		return res
	}
	res.targetURL = h.targetURL
	res.filePath = h.filePath
	return res
}

// Upload upload file by FormData
func (h *Files) Upload() *Response {
	res := h.checkUpload()
	if res.err != nil {
		return res
	}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	flieNames := strings.Split(h.filePath, "/")
	fileName := flieNames[len(flieNames)-1]
	var fileWriter io.Writer
	contentType := mimetypes.Lookup(fileName)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	fileWriter, err := createFormFile(bodyWriter, fileName, contentType)
	if err != nil {
		res.err = err
		return res
	}
	fh, err := os.Open(h.filePath)
	if err != nil {
		res.err = err
		return res
	}
	_, res.err = io.Copy(fileWriter, fh)
	fh.Close()
	if err != nil {
		return res
	}
	bodyWriter.Close()
	request, err := http.NewRequest(http.MethodPost, h.targetURL, bodyBuf)
	if err != nil {
		res.err = err
		return res
	}
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	for k, v := range h.header {
		request.Header.Set(k, v)
	}
	res.resp, res.err = h.client.Do(request)
	return res
}

// UploadByStream upload by stream
func (h *Files) UploadByStream() *Response {
	res := h.checkUpload()
	if res.err != nil {
		return res
	}
	file, err := os.Open(h.filePath)
	if err != nil {
		res.err = err
		return res
	}
	defer file.Close()
	request, err := http.NewRequest(http.MethodPost, h.targetURL, file)
	if err != nil {
		res.err = err
		return res
	}
	request.Header.Set("Content-Type", "binary/octet-stream")
	for k, v := range h.header {
		request.Header.Set(k, v)
	}
	res.resp, res.err = h.client.Do(request)
	return res
}

func (h *Files) checkDownload() *Response {
	res := &Response{filePath: h.filePath}
	if h.targetURL == "" {
		res.err = ErrEmptyTargetURL
		return res
	}
	res.targetURL = h.targetURL
	return res
}

// Download will get filename from 'Content-Disposition' if savePath is empty.
func (h *Files) Download() *Response {
	res := h.checkDownload()
	if res.err != nil {
		return res
	}
	request, err := http.NewRequest(http.MethodGet, h.targetURL, nil)
	if err != nil {
		res.err = err
		return res
	}
	for k, v := range h.header {
		request.Header.Set(k, v)
	}
	res.resp, res.err = h.client.Do(request)
	if res.err != nil {
		return res
	}
	if h.filePath == "" {
		_, params, err := mime.ParseMediaType(res.resp.Header.Get("Content-Disposition"))
		if err == nil {
			h.filePath = params["filename"]
		} else {
			h.filePath = "unknown"
		}
		res.filePath = h.filePath
	}
	out, err := os.Create(h.filePath)
	if err != nil {
		res.err = err
		return res
	}
	_, res.err = io.Copy(out, res.resp.Body)
	out.Sync()
	out.Close()
	res.resp.Body.Close()
	return res
}

// Head ...
func (h *Files) Head() *Response {
	res := h.checkDownload()
	if res.err != nil {
		return res
	}
	request, err := http.NewRequest(http.MethodHead, h.targetURL, nil)
	if err != nil {
		res.err = err
		return res
	}
	for k, v := range h.header {
		request.Header.Set(k, v)
	}
	res.resp, res.err = h.client.Do(request)
	return res
}

// Get will get response from target URL.
func (h *Files) Get() *Response {
	res := h.checkDownload()
	if res.err != nil {
		return res
	}
	request, err := http.NewRequest(http.MethodGet, h.targetURL, nil)
	if err != nil {
		res.err = err
		return res
	}
	for k, v := range h.header {
		request.Header.Set(k, v)
	}
	res.resp, res.err = h.client.Do(request)
	_, params, err := mime.ParseMediaType(res.resp.Header.Get("Content-Disposition"))
	if err == nil {
		res.filePath = params["filename"]
	}
	return res
}
