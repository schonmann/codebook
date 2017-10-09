package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

const (
	gzipCompressor = "gzip"
	noneCompressor = "none"
)

//Compressor is the database compressor interface.
type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

//CreateCompressor is the builder method for each compressor type.
func CreateCompressor(t string) Compressor {
	switch t {
	case gzipCompressor:
		return GzipCompressor{}
	case noneCompressor:
		return NoneCompressor{}
	default:
		panic("Wrong compressor providr '" + t + "'")
	}
}

//GzipCompressor is a compressor implementation.
type GzipCompressor struct{}

//Compress ...
func (c GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		slog.Errorf(err.Error())
		return nil, err
	}
	if err := zw.Close(); err != nil {
		slog.Errorf(err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

//Decompress ...
func (c GzipCompressor) Decompress(data []byte) ([]byte, error) {
	var bufin bytes.Buffer
	if _, err := bufin.Write(data); err != nil {
		return nil, err
	}
	var bufout bytes.Buffer
	zipout, err := gzip.NewReader(&bufin)
	if err != nil {
		return nil, err
	}
	bufout.ReadFrom(zipout)
	if err := zipout.Close(); err != nil {
		return nil, err
	}
	return bufout.Bytes(), nil
}

func main() {
  	compressor := CreateCompressor(gzipCompressor)
  	foo := "I will be compressed!"
  	bar, err := compressor.Compress([]byte(foo))
  	if err != nil {
		panic("Something went bad!")
  	}
  	fmt.Printf("Compressed foo: %s\n", string(bar))
}
