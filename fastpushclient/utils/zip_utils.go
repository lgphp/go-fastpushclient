package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Gzip(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func UnGzip(data []byte) []byte {
	rdata := bytes.NewReader(data)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	return s
}
