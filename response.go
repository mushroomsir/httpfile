package httpfile

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

// Response ...
type Response struct {
	err        error
	resp       *http.Response
	statusCode int
	filePath   string
}

func (a *Response) Error() error {
	if a.resp != nil && a.resp.StatusCode >= 400 {
		return errors.New(a.BodyString())
	}
	return a.err
}

// Body ...
func (a *Response) Body() []byte {
	if a.err != nil {
		return nil
	}
	defer a.resp.Body.Close()
	respBody, err := ioutil.ReadAll(a.resp.Body)
	a.err = err
	return respBody
}

// BodyString ...
func (a *Response) BodyString() string {
	return string(a.Body())
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
func (a *Response) FileSize() int64 {
	if a.err == nil {
		return 0
	}
	stat, err := os.Stat(a.filePath)
	if err != nil {
		a.err = err
		return stat.Size()
	}
	return 0
}

// Resp ...
func (a *Response) Resp() *http.Response {
	if a.resp == nil {
		return &http.Response{}
	}
	return a.resp
}
