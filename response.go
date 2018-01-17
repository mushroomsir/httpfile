package httpfile

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

// Body ...
func (a *Response) Body() ([]byte, error) {
	if a.err != nil {
		return nil, a.err
	}
	defer a.resp.Body.Close()
	return ioutil.ReadAll(a.resp.Body)
}

// BodyString ...
func (a *Response) BodyString() (string, error) {
	res, err := a.Body()
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

// Discard ...
func (a *Response) Discard() {
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

// FileName get FileName from filePath or Content-Disposition of response
func (a *Response) FileName() string {
	_, file := filepath.Split(a.filePath)
	return file
}

// Resp ...
func (a *Response) Resp() *http.Response {
	if a.resp == nil {
		return &http.Response{}
	}
	return a.resp
}
