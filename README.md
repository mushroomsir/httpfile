# httpfile
[![Build Status](https://img.shields.io/travis/mushroomsir/httpfile.svg?style=flat-square)](https://travis-ci.org/mushroomsir/httpfile)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/httpfile.svg?style=flat-square)](https://coveralls.io/github/mushroomsir/httpfile?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/mushroomsir/httpfile/blob/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/mushroomsir/httpfile)

## Features

- Upload multi file by formdata
- Upload file by stream
- Custom file `Content-Type` for upload
- Support add extra field to formdata
- Download file 

## Installation

```sh
go get github.com/mushroomsir/httpfile
```

## Usage
```go
package main

import "github.com/mushroomsir/httpfile"
import "fmt"

func main() {
	res, err := httpfile.Upload(httpfile.UploadOptions{
		FileItems: httpfile.NewFileItems("test.gif"),
		TargetURL: "TargetURL",
	})

	fmt.Println(err)
	fmt.Println(res)
}

```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/httpfile/blob/master/LICENSE).
