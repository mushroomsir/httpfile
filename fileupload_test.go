package httpfile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeadRequest(t *testing.T) {
	require := require.New(t)

	res := NewReq("https://http2.golang.org/reqinfo").Head()
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
