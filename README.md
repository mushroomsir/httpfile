# httpfile
[![Build Status](https://img.shields.io/travis/mushroomsir/httpfile.svg?style=flat-square)](https://travis-ci.org/mushroomsir/httpfile)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/httpfile.svg?style=flat-square)](https://coveralls.io/github/mushroomsir/httpfile?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/mushroomsir/httpfile/blob/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/mushroomsir/httpfile)

## Features

- Easy to use
- Upload file by http FormFata
- Upload file by http Stream
- Download file to local

## Installation

```sh
go get -u github.com/mushroomsir/httpfile
```

## Usage
```go
func TestHTTPFile(t *testing.T) {
	require := require.New(t)
	
	res := httpfile.NewReq(fileURL(), "testdata/test.gif").Upload()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())

	res = httpfile.NewReq(fileURL(), "testdata/download/test1.gif").SetHeader("filename", "test.gif").Download()
	require.Nil(res.Error())
	require.Equal(200, res.StatusCode())
	require.Equal("bytes", res.GetHeader("Accept-Ranges"))
}

```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/httpfile/blob/master/LICENSE).
