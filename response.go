package httpfile

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Response ...
type Response struct {
	err       error
	resp      *http.Response
	filePath  string
	targetURL string
}

func (a *Response) Error() error {
	if a.err != nil {
		return a.err
	}
	if a.resp != nil && a.resp.StatusCode >= 400 {
		res, err := a.BodyString()
		if err != nil {
			return err
		}
		return errors.New(res)
	}
	return a.err
}

// Bytes ...
func (a *Response) Bytes() ([]byte, error) {
	if a.err != nil {
		return nil, a.err
	}
	defer a.resp.Body.Close()
	return ioutil.ReadAll(a.resp.Body)
}

// BodyString ...
func (a *Response) BodyString() (string, error) {
	res, err := a.Bytes()
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// StatusCode ...
func (a *Response) StatusCode() int {
	if a.resp != nil {
		return a.resp.StatusCode
	}
	return 0
}

// Unmarshal ...
func (a *Response) Unmarshal(result interface{}) error {
	if a.err != nil {
		return a.err
	}
	defer a.resp.Body.Close()
	return json.NewDecoder(a.resp.Body).Decode(result)
}

// GetHeader ...
func (a *Response) GetHeader(k string) string {
	if a.err != nil {
		return ""
	}
	return a.resp.Header.Get(k)
}

// Close ...
func (a *Response) Close() {
	if a.resp != nil {
		a.resp.Body.Close()
	}
}

// ClearError ...
func (a *Response) ClearError() {
	a.err = nil
}

// FileSize ...
func (a *Response) FileSize() (int64, error) {
	if a.err != nil {
		return 0, a.err
	}
	stat, err := os.Stat(a.filePath)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

// TargetURL ...
func (a *Response) TargetURL() string {
	return a.targetURL
}

// ContentLength ...
func (a *Response) ContentLength() *Int64Result {
	result := &Int64Result{
		Err: a.err,
	}
	if result.Err != nil {
		return result
	}
	result.Length, result.Err = strconv.ParseInt(a.resp.Header.Get("Content-Length"), 10, 64)
	return result
}

// FileName get FileName from filePath or Content-Disposition of response
func (a *Response) FileName() string {
	_, file := filepath.Split(a.filePath)
	return file
}

// Resp ...
func (a *Response) Resp() *http.Response {
	return a.resp
}

// Body ...
func (a *Response) Body() io.ReadCloser {
	return a.resp.Body
}

// Int64Result ...
type Int64Result struct {
	Length int64
	Err    error
}
