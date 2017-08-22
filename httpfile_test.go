package httpfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleUpload(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
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
