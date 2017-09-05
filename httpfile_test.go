package httpfile

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	SetHTTPClient(&http.Client{})
}

func TestHead(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	_, err := Head("")
	require.NotNil(err)

	res, err := Head("https://http2.golang.org/reqinfo", map[string]string{"filename": "test.gif"})
	require.Nil(err)
	assert.Equal("text/plain", res.Header.Get("Content-Type"))
}
func TestError(t *testing.T) {
	assert := assert.New(t)
	res, err := UploadFile("", "")
	assert.Nil(res)
	assert.NotNil(err)
}
func TestSimpleUpload(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	_, err := Upload(UploadOptions{
		FileItems: NewFileItems(uploadDir("test.gif"), ""),
		TargetURL: "",
	})
	assert.NotNil(err)

	res, err := Upload(UploadOptions{
		FileItems: NewFileItems(uploadDir("test.gif"), ""),
		TargetURL: fileURL(),
	})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)

	resp, err := Download(fileURL(), downloadDir("test.gif"), map[string]string{"filename": "test.gif"})
	require.Nil(err)
	assert.Equal(200, resp.StatusCode)
	assert.Equal("bytes", resp.Header.Get("Accept-Ranges"))

	resp, err = Download(fileURL(), "", map[string]string{"filename": "test.gif"})
	require.Nil(err)
	assert.Equal(200, resp.StatusCode)
	assert.Equal("bytes", resp.Header.Get("Accept-Ranges"))
}
func TestUploadReader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	res, err := UploadReader(strings.NewReader(`{"Username": "12124", "Password": "testinasg", "Channel": "M"}`), "")
	assert.Nil(res)
	assert.NotNil(err)
	file, err := os.Open(uploadDir("test.gif"))
	assert.Nil(err)
	defer file.Close()

	res, err = UploadReader(file, fileURL(), map[string]string{"testheader": "123"})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)
	assert.Equal("123", res.Header.Get("testheader"))
}
func TestDownload(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	filename := "test.gif"
	res, err := Upload(UploadOptions{
		FileItems: NewFileItems(uploadDir(filename), ""),
		TargetURL: fileURL(),
	})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)

	resp, err := Download(fileURL(), "", map[string]string{"filename": filename})
	require.Nil(err)
	assert.Equal(200, resp.StatusCode)
	assert.Equal("bytes", resp.Header.Get("Accept-Ranges"))
	_, err = os.Stat(filename)
	assert.False(os.IsNotExist(err))
	os.Remove(filename)
}

func TestHeader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	res, err := Upload(UploadOptions{
		FileItems: []FileItem{NewFileItem(uploadDir("test.gif"), "")},
		TargetURL: fileURL(),
		Header:    map[string]string{"testheader": "123"},
	})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)
	assert.Equal("123", res.Header.Get("testheader"))
}
func TestMultiFileAndExtraField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	res, err := Upload(UploadOptions{
		FileItems: []FileItem{
			NewFileItem(uploadDir("test.gif")),
			NewFileItem(uploadDir("test.bmp")),
		},
		TargetURL:  fileURL(),
		ExtraField: map[string]string{"ExtraField": "123"},
	})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)
	assert.Equal("123", res.Header.Get("ExtraField"))
}
func TestUploadStream(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	res, err := UploadFile(uploadDir("test.gif"), fileURL(), map[string]string{"testheader": "123"})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)
	assert.Equal("123", res.Header.Get("testheader"))
}

func downloadDir(fileName string) string {
	return "testdata/download/" + fileName
}

func uploadDir(fileName string) string {
	return "testdata/" + fileName
}

func fileURL() string {
	return testServer.URL + "/file"
}
