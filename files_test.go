package httpfile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeadRequest(t *testing.T) {
	require := require.New(t)

	res := NewReq("").Head()
	require.Equal(ErrEmptyTargetURL, res.Error())

	res = NewReq("x").Head()
	require.NotNil(res.Error())

	res = NewReq("https://http2.golang.org/reqinfo").SetHeader("k", "v").Head()
	require.Equal(nil, res.Error())
	require.Equal("text/plain", res.GetHeader("Content-Type"))
}
func TestHTTPFile(t *testing.T) {
	require := require.New(t)

	err := NewReq("").Upload().Error()
	require.Equal(ErrEmptyTargetURL, err)

	err = NewReq("x", "").Upload().Error()
	require.Equal(ErrEmptyFilePath, err)

	res := NewReq(fileURL(), "testdata/test.gif").Upload()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())

	res = NewReq(fileURL(), "testdata/download/test1.gif").SetHeader("filename", "test.gif").Download()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())
	require.Equal("bytes", res.GetHeader("Accept-Ranges"))

	res = NewReq(fileURL()+"?filename=test.gif", "testdata/download/test2.gif").Download()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())
	require.Equal("bytes", res.GetHeader("Accept-Ranges"))
}

func TestSet(t *testing.T) {
	require := require.New(t)
	err := NewReq("xxx", "xxx").SetHTTPClient(defaultHTTPClient).Upload().Error()
	require.NotNil(err)

}
func TestFileUploadStream(t *testing.T) {
	require := require.New(t)

	res := NewReq("").UploadByStream()
	require.Equal(ErrEmptyTargetURL, res.Error())

	res = NewReq("x", "").UploadByStream()
	require.Equal(ErrEmptyFilePath, res.Error())

	res = NewReq("x", "x").UploadByStream()
	require.NotNil(res.Error())

	res = NewReq(fileURL(), "testdata/download/test1.gif").SetHeader("filename", "test.gif").SetHeader("testheader", "123").UploadByStream()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())
	require.Equal("123", res.GetHeader("testheader"))

	res = NewReq("").Download()
	require.Equal(ErrEmptyTargetURL, res.Error())

	res = NewReq("x").Download()
	require.NotNil(res.Error())

	res = NewReq(fileURL()+"?filename=test.gif", "testdata/download/Streamtest2.gif").Download()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())
	require.Equal("bytes", res.GetHeader("Accept-Ranges"))
	size, err := res.FileSize()
	require.Nil(err)
	require.Equal(int64(185210), size)

	res = NewReq(fileURL()+"?filename=test.gif", "").Download()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())
	require.Equal("bytes", res.GetHeader("Accept-Ranges"))
	size, err = res.FileSize()
	require.Nil(err)
	require.Equal(int64(185210), size)
}
