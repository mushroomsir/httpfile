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
		TargetURL: uploadURL(),
	})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)
}

func TestHeader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	res, err := Upload(UploadOptions{
		FileItems: []FileItem{NewFileItem(uploadDir("test.gif"), "")},
		TargetURL: uploadURL(),
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
		TargetURL:  uploadURL(),
		ExtraField: map[string]string{"ExtraField": "123"},
	})
	require.Nil(err)
	assert.Equal(200, res.StatusCode)
	assert.Equal("123", res.Header.Get("ExtraField"))
}
func TestUploadStream(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	res, err := UploadFile(uploadDir("test.gif"), uploadURL(), map[string]string{"testheader": "123"})
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

func uploadURL() string {
	return testServer.URL + "/upload"
}
