// Code generated by go-bindata. DO NOT EDIT.
// sources:
// data/gjisho.glade
// data/gjisho.glade~
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dataGjishoGlade = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5d\x6f\x77\xe2\x36\x97\x7f\xcf\xa7\xd0\xfa\xc5\x73\xda\xb3\x90\xc9\x64\xfa\x74\xa7\x6d\x42\x97\x24\x24\x71\x87\x40\x16\x9c\x66\xe7\x15\x47\xd8\x17\xd0\xc4\x48\x5e\x49\x40\xe8\xd9\x2f\xb4\x5f\x63\x3f\xd9\x73\x2c\x43\x82\xc1\x7f\x64\x1b\x26\x24\x75\xdf\x74\x62\xa4\x9f\xa4\x2b\xdd\x3f\xba\xba\xd2\x3d\xfd\xfd\x69\xe2\xa2\x19\x70\x41\x18\x3d\x33\x3e\x1e\x1d\x1b\x08\xa8\xcd\x1c\x42\x47\x67\xc6\xbd\x75\x55\xfb\x6c\xfc\x5e\xaf\x9c\xfe\x5b\xad\x86\xae\x81\x02\xc7\x12\x1c\x34\x27\x72\x8c\x46\x2e\x76\x00\x7d\x3a\xfa\xf4\xf3\xd1\x31\xaa\x54\xac\x31\xa0\x5b\xd3\x42\x2d\x62\x03\x15\x80\x7e\xb8\x35\xad\x1f\x2b\x95\x0b\xe6\x2d\x38\x19\x8d\x25\xfa\xc1\xfe\x11\x9d\x1c\x9f\x1c\x23\x13\x53\xf4\x07\x1b\x53\xc1\x68\xa5\x72\x07\x7c\x42\x84\xdf\x3a\x22\x02\x8d\x81\xc3\x60\x81\x46\x1c\x53\x09\x4e\x15\x0d\x39\x00\x62\x43\x64\x8f\x31\x1f\x41\x15\x49\x86\x30\x5d\x20\x0f\xb8\x60\x14\xb1\x81\xc4\x84\x12\x3a\x42\x18\xd9\xcc\x5b\x54\xd8\x10\xc9\x31\x11\x48\xb0\xa1\x9c\x63\x0e\x08\x53\x07\x61\x21\x98\x4d\x54\xbf\x1d\x66\x4f\x27\x40\x25\x96\x7e\x7b\x43\xe2\x82\x40\x3f\xc8\x31\x20\xa3\xb7\xac\x61\xfc\xa8\x1a\x71\x00\xbb\x15\x42\x91\xff\xdb\xea\x27\x35\x6a\x36\x95\x88\x83\x90\x9c\xd8\x3e\x46\x15\x11\x6a\xbb\x53\x9f\x5a\xcf\x3f\xbb\x64\x42\x96\x2d\xf8\xd5\xd5\xe0\x45\x45\x32\x34\x15\x50\x55\xfd\xac\xa2\x09\x73\xc8\xd0\xff\x3f\xa8\x61\x79\xd3\x81\x4b\xc4\xb8\x8a\x1c\xe2\x43\x0f\xa6\x12\xaa\x48\xf8\x1f\x15\x2d\xab\xfe\x38\x3e\x30\x8e\x04\xb8\x6e\xc5\x66\x1e\x01\x81\xd4\x58\x5f\x7a\xa7\xca\xf8\x5d\xf7\x7c\x82\xca\x25\x89\x84\xff\x65\x3e\x66\x93\xf0\x48\x88\xa8\x0c\xa7\x9c\x12\x31\x06\x55\xc7\x61\x48\x30\xd5\xe2\x37\xb0\xa5\xff\xc5\x2f\x3e\x64\xae\xcb\xe6\xfe\xd0\x6c\x46\x1d\xe2\x8f\x48\xfc\x1a\x4c\x34\x1e\xb0\x19\xa8\xb1\x04\x73\x4b\x99\x24\x76\x40\x6e\x35\x01\xde\xcb\xac\x2e\x7f\x12\x63\xec\xba\x68\x00\x4b\x82\x81\x83\x08\xad\xf8\x9f\x56\xc3\xe1\x7e\xf3\x42\x62\x2a\x09\x76\x91\xc7\xb8\x6a\x6f\x73\x98\x47\x95\x8a\x75\xd3\x44\xbd\xce\x95\xf5\xd0\xe8\x36\x91\xd9\x43\x77\xdd\xce\x9f\xe6\x65\xf3\x12\x19\x8d\x1e\x32\x7b\x46\x15\x3d\x98\xd6\x4d\xe7\xde\x42\x0f\x8d\x6e\xb7\xd1\xb6\xbe\xa2\xce\x15\x6a\xb4\xbf\xa2\x2f\x66\xfb\xb2\x8a\x9a\xff\x7d\xd7\x6d\xf6\x7a\xa8\xd3\xad\x98\xb7\x77\x2d\xb3\x79\x59\x45\x66\xfb\xa2\x75\x7f\x69\xb6\xaf\xd1\xf9\xbd\x85\xda\x1d\x0b\xb5\xcc\x5b\xd3\x6a\x5e\x22\xab\x83\xfc\x06\x97\x50\x66\xb3\xe7\x83\xdd\x36\xbb\x17\x37\x8d\xb6\xd5\x38\x37\x5b\xa6\xf5\xb5\x5a\xb9\x32\xad\xb6\x8f\x79\xd5\xe9\xa2\x06\xba\x6b\x74\x2d\xf3\xe2\xbe\xd5\xe8\xa2\xbb\xfb\xee\x5d\xa7\xd7\x44\x8d\xf6\x25\x6a\x77\xda\x66\xfb\xaa\x6b\xb6\xaf\x9b\xb7\xcd\xb6\x75\x84\xcc\x36\x6a\x77\x50\xf3\xcf\x66\xdb\x42\xbd\x9b\x46\xab\xe5\x37\x55\x69\xdc\x5b\x37\x9d\xae\xdf\x3f\x74\xd1\xb9\xfb\xda\x35\xaf\x6f\x2c\x74\xd3\x69\x5d\x36\xbb\x3d\x74\xde\x44\x2d\xb3\x71\xde\x6a\x06\x4d\xb5\xbf\xa2\x8b\x56\xc3\xbc\xad\xa2\xcb\xc6\x6d\xe3\xba\xa9\x6a\x75\xac\x9b\x66\xb7\xe2\x17\x0b\x7a\x87\x1e\x6e\x9a\xfe\x27\xbf\xbd\x46\x1b\x35\x2e\x2c\xb3\xd3\xf6\x87\x71\xd1\x69\x5b\xdd\xc6\x85\x55\x45\x56\xa7\x6b\x3d\x57\x7d\x30\x7b\xcd\x2a\x6a\x74\xcd\x9e\x4f\x90\xab\x6e\xe7\xb6\x5a\xf1\xc9\xd9\xb9\xf2\x8b\x98\x6d\xbf\x5e\xbb\x19\xa0\xf8\xa4\x46\xa1\x19\xe9\x74\xd5\xdf\xf7\xbd\xe6\x33\x20\xba\x6c\x36\x5a\x66\xfb\xba\x87\xcc\x76\x68\xfa\x8e\x2a\x95\xc6\x54\x8e\x19\xff\x75\x5d\x28\xa0\x53\x82\xa9\xc7\xc9\x04\x8e\xff\x79\xfc\xcb\x7f\x8e\x26\x98\xb8\x47\x36\x9b\xd4\x2b\x95\x5a\xad\x5e\x39\x25\x54\x02\x1f\x62\x1b\xea\x15\x84\x4e\x39\xfc\xcf\x94\x70\x10\xc8\x25\x83\x33\x63\x24\x1f\xff\xdd\x78\x91\x67\x9f\x8e\x4e\x4e\x8c\x0f\xaa\x9c\x2f\xc1\x9e\x6b\xd6\x96\x9c\x55\x93\x0b\x0f\x90\xcf\x2f\x3e\xf2\x56\x29\x8a\x27\x80\xae\xff\x20\x62\xcc\xa2\x0b\x38\x20\x6c\x4e\x3c\xc5\xec\x0d\xf4\x07\xf6\x30\x05\x01\xc8\x09\xa4\x03\xe6\x0b\x34\x64\x1c\x5d\xb7\x3b\xb7\xcd\xa3\x68\x84\x17\x1e\xda\x94\x8d\xd1\xe5\xb1\xa2\x97\xd0\xa2\xd7\x0a\x81\x05\x6c\x6d\xbb\x58\x88\x33\xe3\x5a\x3e\xde\x31\x8f\xcd\x80\xdf\x02\x9d\x1a\x88\x38\x67\xc6\x04\x13\xaa\xfe\xf2\xcb\x23\x74\xea\x71\xe6\x01\x97\x0b\xe4\x53\xe0\xcc\xb0\x31\xed\x0f\x99\x3d\x15\x46\xfd\x0a\xbb\x02\x4e\x3f\xac\x0a\x2c\xcb\xdb\x63\xe2\x3a\xc1\xbf\xa3\xda\x3b\x67\x4f\xc6\xea\xd7\x6d\xf4\x19\x11\x64\xe0\x82\x51\xb7\xf8\x74\x0b\x3a\x4f\x77\xa2\xea\x30\x4e\x56\x82\xdf\xa8\xcf\x80\x4b\x62\x63\x37\xb2\x62\x68\x2c\xd1\xe3\xb9\x65\x0e\xb8\xe7\x53\x29\x7d\xb0\xb5\x92\x39\xc6\x96\x32\xbe\x2c\xd5\x38\xd8\x40\x66\x20\xfa\x0e\x0c\xf1\xd4\x95\xd9\x6a\x63\xb5\x64\xfb\xfe\x1f\x46\x1d\x7b\xde\x11\x1e\xb0\xa9\xd4\xad\x2d\xe1\x49\x1a\x48\x72\x4c\x85\x8b\x25\x1e\xb8\x70\x66\x2c\x40\x18\xf5\x86\x8f\xb2\xe4\xa1\x68\xb0\xd3\x0f\x01\x7d\x43\xdf\x3c\x6c\x3f\x12\x3a\x4a\x6e\x14\x9e\x3c\x4c\x9d\x84\x45\x10\x55\x69\x48\x5c\x37\x1b\x65\x3c\x26\x48\xb0\x6e\x4e\xe2\x46\xb0\xd5\xdd\xd3\x0f\x61\x96\xd8\x18\xe3\xf6\xf8\x36\x1b\x15\xd3\xc1\x44\x31\xa4\xcf\x9a\x3a\x0b\xfc\xa5\x97\x1f\xb7\x8b\x6f\xf4\x70\xad\x77\xeb\x3d\xdb\x5e\xe8\x0d\xcf\x73\x89\xad\xb8\xe6\x81\x50\x87\xcd\x03\x71\x81\x3d\x6f\xf9\x67\x4e\x79\xb1\x39\x91\x33\xa0\x52\x18\xf5\xeb\xcb\x2f\xfd\xf3\x7b\xcb\xea\xb4\xfb\x4a\x47\xf7\x6f\x1b\xbd\x2f\xe8\x7f\x91\xff\xfd\x4b\xf3\xeb\xf6\xc7\x9e\xd5\xbd\xbf\xb0\xee\xbb\x4d\xf5\x31\xad\x99\x25\x5f\xf4\xe7\xc4\x91\x63\xa3\xfe\xf9\xf8\x58\xb7\xc6\x18\x7c\x11\x6d\xd4\x7f\x8e\xa8\x22\xc8\x88\x62\x77\x59\x61\xa0\x84\x42\xcd\xe3\x20\x44\x4d\x0d\xcb\x40\x63\x4c\x1d\x17\xf8\x99\x31\x57\x34\x0b\xe4\xc6\x9d\x5f\xc2\x40\x62\x8e\x3d\x0f\x9c\x33\x83\xb2\x40\x4f\x6d\x02\x3e\xc2\x22\x11\xed\x0b\x2c\x12\xa0\x0e\x50\x2e\x6b\x88\xd7\x2e\xcc\x00\xbb\xc0\x83\xc5\x26\x00\x73\x7b\xfc\xfc\x6d\xbf\xf2\x36\x93\x28\x51\xe2\x4e\x31\x5d\xdf\xb7\x23\x8c\xba\x70\x89\x03\x35\xa5\xce\x63\x31\xb6\xc6\xaf\x33\x35\x71\x3d\x50\x2b\xb9\xef\x9b\x40\x20\xa4\x51\xff\xb4\xbd\x3c\x0b\x13\xab\x08\xc1\xa2\xea\x6a\x2b\xe2\x64\x8a\x65\xa1\x5a\xa1\xb1\x17\x1d\xff\x4e\x68\x90\x4c\x87\x68\x5a\xf4\x14\xdf\x34\xa9\xe4\x8b\x75\x46\x0a\x3e\xc4\xc1\x14\x24\x54\x0a\xb1\xf2\x42\x4c\x30\x1f\x11\xda\x17\x12\x73\x69\xd4\xff\x59\x00\x02\x7c\x7b\xa1\x08\x80\x64\x5e\x31\x80\x01\x93\x92\x4d\x72\x62\xf8\x56\x3e\xe6\x8b\x3e\xb1\x9f\x4d\x35\x70\x88\xac\x0d\x09\x75\x6a\x62\x31\x19\x30\x97\xd8\x85\x81\x7d\x53\x70\x16\x58\x71\x9a\x0b\x3c\x15\x52\x80\x92\x94\xb3\x4c\x80\x9a\x8a\x70\x6d\x5d\x27\x6a\xc3\xf4\x36\x02\xa4\x9a\x3d\xc6\x74\x04\xce\x56\x13\x17\xab\xef\x5a\xe0\x51\x76\x6d\x98\x5e\x51\x36\x6e\x12\x45\xb5\xec\xdd\x24\x80\x74\xdb\x37\xa9\xf6\x8b\x85\x99\xa8\x65\x56\xa3\x4f\x1c\xde\x86\x7d\x1c\xfe\x29\x9e\x6e\xc9\x34\x2b\x44\xaf\xfc\xb4\xca\x4d\xa7\x04\x1a\xc5\xd2\x27\x8b\x36\xec\xd9\x9c\xb9\x2e\x38\x21\x5b\x3d\xad\xf7\xbb\x52\x8c\x79\xaa\x8f\x85\xea\xf1\x00\xf3\xbe\xe7\x4b\xb2\x85\x51\xa7\x30\x03\x9e\x8a\x13\xe2\x62\x70\x46\x50\xe3\x80\xed\x71\x04\x0f\x77\x41\x4c\x5d\x29\x9a\xce\x08\xba\xab\x22\x1a\xec\x9c\x55\xf9\xfe\x49\x60\xee\x31\x5f\x5b\xbd\x8a\xa2\xcd\x2d\x21\xc4\x18\x3b\x6c\xbe\x34\x67\x29\xa3\x7a\x20\x89\xc4\x89\x26\x50\x8b\x08\xe9\x5b\x6b\x21\x13\x5f\x4d\x4d\x02\xc9\x76\x42\xb6\x5d\x91\x2e\x92\x7c\xe0\x42\xe0\x47\x99\x30\x07\x8c\xfa\x80\xb3\x79\x06\xb4\xd0\x32\xe6\x6c\x5e\x0b\xf0\xe2\x97\x71\x97\xcd\x7b\xcf\x45\x34\x55\x5e\xba\x66\x4a\x14\xcf\xe9\x00\xdf\x55\xb6\xe7\x11\x33\xbb\x10\xed\x11\x4e\x96\xf0\x40\xb3\x8a\xf6\x68\xd2\x44\x14\x7e\x43\xfe\xb2\x18\xfd\x97\xee\x2f\xd3\xf3\x13\xf4\xc0\xc3\x1c\x4b\xf6\x8a\x4e\x81\x37\x34\x19\x31\x2b\x76\x57\x93\xb1\xbd\xf5\x7e\x55\xdf\x4c\xa6\x2d\xb6\xae\x53\x26\xd5\xa0\xda\x87\x87\x25\x6b\xd5\x3c\x46\x54\x16\xab\x32\xc5\xb4\x79\x5d\x47\x4b\x46\x33\x2d\xde\x5f\x54\x78\x28\xbb\x18\x4e\x14\x46\x2e\xdf\x51\x3a\x6d\xa2\xe9\xd3\xc2\x03\x70\x03\x1b\x6d\xb9\x9d\xff\x82\x29\x0e\xbe\xbe\x25\x33\x0d\xcb\x65\x58\x85\x48\x2e\xb8\x5e\x74\x65\xd2\xd9\xd8\x05\x03\xcd\xb0\x3b\x85\x33\xe3\xe3\xd1\x49\xb2\x6d\xe5\x4b\x4f\xdd\xd6\x4e\x85\x5c\xb8\x90\xda\x23\x35\x1d\x61\xa7\x4a\xed\x11\x53\x9c\xde\x91\x54\xfc\x74\x4b\x50\xc7\x4f\x81\x76\xe1\xab\x88\x02\xc9\xe6\xaf\x88\x42\xc8\xe4\xb3\x40\xe9\x7e\x0b\x94\x6e\x1c\xef\x90\xd3\xbe\x91\xbf\x35\xab\xfd\xc7\xc9\xe7\x43\x61\xb6\x6f\xa4\xe4\xb6\x74\x04\xed\xed\xd1\x0b\x45\x5e\x9d\xdb\x80\x4a\xbe\xb8\x04\x89\x89\x2b\xde\x1e\xb7\x15\x3f\x28\x49\x80\xd1\x3f\x2c\x49\x00\x51\x07\x26\x1f\x35\xa5\xef\x6e\x4e\x4d\xa2\x70\xe6\x1c\x7b\xc5\xe6\xe8\x09\xbb\x64\xa4\xaf\x48\x36\xdd\x39\xcb\x43\x15\xa8\xb9\x84\x3e\xae\xf9\x73\x56\xdf\x5b\xea\xb3\xb6\x0f\x27\x97\x40\x53\x6b\xbd\xe6\x04\x8b\xbd\x14\x68\xe9\x08\xa9\xa1\x3f\xdb\x14\xf9\x2e\x02\xcd\x73\xb1\x0d\x63\xe6\x3a\xc0\x93\x0e\xb7\xbe\xbb\x17\x2f\x93\x6f\x2b\x69\x6d\xe4\x75\xf7\xe5\x5b\x02\xb9\xac\xc6\x98\xa9\x8e\x24\x81\xae\x83\x21\x1c\xf9\x32\x61\x1c\x4c\x3a\x64\x31\xb1\x2f\x51\x3d\x7f\x95\x90\x8e\x98\x38\x98\xa9\xf7\xfd\xc3\x3a\x82\x40\x2d\xed\x88\x98\x42\x24\x2c\x4a\xc6\xa8\xfa\xdf\x29\x32\x24\xda\x6f\xba\x33\xaa\xec\x82\x32\xe5\x01\x7a\xde\x03\xf4\x1c\x0b\x42\x62\xfb\xb1\x37\x27\xd2\x1e\x47\x4a\x99\xb8\x01\x1d\x84\x43\x6c\xbc\xb4\xcb\x6c\xa0\x32\xfd\xa8\x38\x0a\xa1\x0c\x2f\x5a\x6d\xc0\xfd\x75\x60\xd4\x57\x7a\x47\x2d\x8b\xbf\x35\x9f\x6a\x6c\x63\xbf\x3b\x9f\x86\x6d\x83\xe0\xd3\xa1\xb3\x6c\x8e\x2d\xba\x56\xf0\xca\xce\x46\x99\x32\xd2\x22\x30\x79\x83\x5a\xf4\x89\x17\x4d\x40\x8d\x10\x94\xb8\x2e\xe7\x22\x5f\x14\x50\x5e\xbf\x46\x14\x56\xae\xd0\x94\x67\x30\x2d\x22\xa2\x1c\x87\x44\x49\x5d\xce\x4d\xc7\x28\xb0\x22\xb4\x8c\xc2\xcb\x7d\x90\x14\x02\xd5\xa6\x2b\x4a\x70\xc4\xe9\xd6\xdf\x31\x85\xa3\x00\x8b\x52\x39\x0a\x33\xa7\x57\x2e\x09\xd2\x0d\x1c\x98\x11\x17\xae\xd4\xb1\x01\x9a\x73\x22\x09\x1d\x89\x5d\xb4\x95\xd9\xf1\x15\x02\xd3\x3f\x0f\x88\xae\xb6\xf2\xdf\x05\x37\x70\x56\x87\x03\x03\xe6\x3a\x69\x0e\xac\x14\xc4\xcd\xd3\x86\x9f\x3e\xfd\x12\xfa\x2f\x0b\xbc\xfe\x31\xc4\x7a\x9d\x74\xe7\x59\xa8\xbc\x96\x23\x2d\x5c\x63\x07\x4e\xb5\x24\xc0\xec\x0e\xb6\x24\xb4\xcc\x67\x75\x21\xb0\x74\xc7\x5b\xb8\x78\x06\xc1\xb5\x1b\x21\xa7\xcc\x36\x75\x94\xf4\xb0\x64\xce\x52\xf8\xe5\x83\xcc\xee\xcb\x4f\x42\x3b\x08\xf1\x96\x39\xca\x20\x84\x56\x4a\x9f\x5d\x4a\x1f\xcd\xb3\xcb\x10\xd8\xc1\x4b\x9f\x52\xca\xe4\x84\xd4\xf7\xc8\xa4\x01\x26\xda\x6c\xb8\x34\xd9\xd2\x11\x4b\x93\xed\x60\x85\xa6\xe6\xf9\x68\x08\xec\xe0\x85\xe6\xca\x64\xc3\xa5\xc5\x56\x5a\x6c\xa5\xc5\x76\xb0\xc2\xe7\xd3\x21\x09\x1f\xfd\x79\xd2\x04\xd5\x03\xd4\x00\xdb\x5f\xb8\x0e\x55\x57\xd1\xb3\x5a\x30\x5b\x71\x04\x44\xfa\x2c\x15\x61\x20\x3d\x64\x42\x7e\xad\xa8\xc2\xf2\xc8\x42\x9b\x78\xd1\x04\x2c\x8f\x2c\xd6\xc1\x0a\x1c\x59\x84\x6e\xd7\x2a\x9f\x93\xff\xe5\x3d\x9e\x62\x6c\xde\xb8\xcd\x4a\x65\x85\xb9\x75\xef\x76\x15\x94\xb9\x7e\xf1\xf6\x99\x8c\x5d\x36\x6f\xbc\xfc\x9e\x21\x62\xf3\xb9\x3d\x35\xb3\xc8\x5f\x16\x67\xc6\x5a\x28\x9f\xb6\x61\xf9\x37\xdd\xe4\xc7\x6f\xa0\xdb\x0c\xa9\xe9\x29\x0d\xc2\x84\xd2\xa5\x05\x13\x2a\x15\x65\xc1\x64\x5a\x44\xfa\xe6\xcb\x97\x22\xb0\x6f\xe6\x52\x45\x69\xfe\xac\xf0\x74\x9e\x23\x81\x27\x3c\xf1\x5c\xc8\xfc\x12\xc9\x73\x23\xa5\x8d\x75\x40\x36\xd6\x6a\x36\x4b\x33\x2b\x01\xb3\x34\x7b\x72\x61\x26\x9a\x3d\xab\x85\x57\x5a\x3e\x09\xa5\x4b\xcb\x27\x54\x2a\xca\xf2\xc9\xba\x8e\xf4\x8d\x9f\x66\x41\xe4\x57\xb8\x83\xb5\xbf\x88\xe6\xbc\x51\xb7\xbb\x0a\x68\xd6\x20\xe1\x3e\x5e\xee\xdb\xfb\xbd\xb0\xcc\x77\x96\x8a\x5e\x0c\x4b\xb0\xc5\xf7\x71\x31\xac\xa1\x94\xef\x39\x3e\x9c\x3b\x60\x59\xae\x6f\x05\x2f\xe5\x2b\x43\x69\x80\xed\xc7\xc8\x97\xf3\x0b\x8f\x26\xaa\x72\xd6\x67\x47\x0b\x6e\x28\x36\xab\x67\x7b\x91\xff\x19\x25\xb4\x7d\xb0\x5d\x62\x3f\x86\x76\x0e\x14\xcf\xc8\x08\x4b\x38\x57\x97\x01\xf6\xf0\x78\xa1\x39\xc1\x23\x38\xf8\x3b\x05\x1b\x18\x6b\x4f\xe2\x8e\x58\xcd\xe3\x30\x23\x6c\x2a\x32\x3c\x8a\x7b\x00\x8f\xda\x1d\xcc\xcb\xa1\xeb\x1c\x3b\x64\x7c\x8e\xb9\x53\x32\xed\x4e\x98\xf6\x2a\xa0\x66\xc9\xb7\xd1\x7c\x4b\xe1\x49\xbe\x55\x9e\xdd\xf9\x93\x90\xd9\x78\xd6\x62\xa3\x91\x0b\xeb\x9c\xbb\xba\x23\x16\xfa\xe5\x3b\x30\xf0\x01\x30\xa0\x54\x43\x5e\x67\xc0\x30\x31\x4a\xfe\x8b\xe0\xbf\x19\x81\x79\xcd\xa7\xd3\x1b\xe3\x40\x6c\x3f\x2e\xfd\x7f\x40\x9d\xac\xab\x6d\x8f\x3a\xf7\x0d\x6f\x6f\x12\xf6\x89\xfa\xdb\x9b\xe2\x8f\xa3\x66\x79\x29\xf4\x00\x12\x3b\x85\x7e\x0c\xb9\x3b\x95\x7b\x64\xb0\xb6\x7d\xdb\x96\x1e\x37\x80\x1d\xe0\xe1\x2d\x5e\x66\x89\x91\x47\x42\xe8\xbb\x72\xe2\xd3\x78\x6d\x7b\xe5\xd9\xbc\x6f\xbb\x4c\x40\x7f\xb0\xd4\x3b\xb1\x5d\xd6\x78\x5c\x76\x5b\xb7\x05\xcf\x60\x27\x68\xb6\x5c\xc2\x36\xa7\x26\x2b\xa6\xc1\xd2\x34\xd7\xfa\x50\x13\xf5\x96\xae\x33\x21\x4e\x4f\xbd\x8e\x23\x21\x5e\x1f\x65\x49\x6d\x52\xe0\xf9\xec\x5c\xcf\x1d\xdf\x02\x9d\x1e\xdc\xba\x53\x55\xfa\x8c\xf6\xd5\xfe\x23\xe3\x8b\xc9\xc5\xd2\x07\x7a\x41\x4a\xc9\x20\x5f\x9d\x4f\x9c\xd8\x8a\xef\x6c\x91\x32\x0f\x68\x6d\x02\x74\xba\xdf\x45\xaa\xab\x39\xf5\x8c\xa1\xcc\x9b\x98\x02\x9a\x30\x25\xc3\xe0\x80\x4d\xe5\x25\xc1\x2e\x1b\x2d\x73\x0b\xae\x7d\x58\x62\x65\x9c\xb8\xed\x85\x2d\xc8\x5f\x09\xc9\x8c\xb6\xe2\xdf\x99\x83\x63\x2c\x89\x2d\x55\xb9\xf0\xa0\x3f\x26\x54\x1a\x75\x47\x75\x39\xb5\x82\x7a\x2a\x0b\xa8\xec\x0f\x19\x57\x79\x36\x83\x68\x89\xb4\x7a\x1e\x67\x23\x8e\x27\xcb\x35\x17\xad\x87\xb7\x38\x23\xc8\x82\x6b\xd4\x8f\x8f\x3e\x1e\xa5\xa6\x3a\x7c\x4e\x44\x1b\xa9\xf9\x5f\xd2\x78\xff\xff\xff\x6d\x65\xaa\x4d\x87\x9e\x4c\x54\x82\xc7\xc8\x04\xa1\x29\x69\x73\xd3\xc0\x97\x09\x71\x8d\xfa\x7a\x46\xdc\x7f\xb8\xf2\xb7\xe8\xa4\xb8\xff\x18\xc9\xdf\xd2\x20\x5d\x36\x62\xeb\x29\xb6\x88\x2f\x8a\x6a\x2a\xa1\x35\x4d\x9d\xe1\x65\x66\xe1\x25\x0f\x4e\x88\xac\x3d\x7d\xdc\x62\xaa\xb0\xca\x77\xc0\x05\x09\x5b\x79\xad\xc6\xc4\x81\x07\xe2\x8c\x40\x6a\xe4\x89\x8c\xc1\x20\x74\x4c\x06\x44\xb6\x55\x7a\xd6\x54\x10\x0e\xc2\x63\x54\x40\x86\x4e\x04\x46\xae\xca\x4f\x4c\xb1\x5b\x53\x7f\x9e\x19\xb3\xc1\xda\x9b\x19\x59\xf3\x4e\xee\x37\xbd\xef\xa6\xa1\xea\x61\x9b\xd0\x51\xcc\x16\x20\x66\x78\xcb\x5c\xb9\x98\x03\x36\x92\xcd\x84\xc0\x44\x48\xcf\x8d\x90\x37\xcd\x81\x8b\x17\x6c\x2a\xfb\xea\xa5\xce\x64\x89\x1f\xad\x76\xe3\x1f\x92\xcc\x72\x5c\x95\x09\xe5\x55\x52\x66\x64\xaa\xb2\xaf\x04\x26\xd1\x64\xda\xd9\x8e\x72\x1b\x5e\x5b\x0b\xaf\x27\xf7\x55\x31\x90\xcb\x07\x91\x77\x9a\xe5\x37\x25\x3e\x12\x2d\xdf\xa5\xdd\xa1\x7a\xde\x48\xf8\x1b\x95\xbd\x37\xa6\x46\x52\xc2\xdf\xbc\x2a\xfd\x60\x64\x7e\x4a\x56\xe0\xd8\x38\xce\x5d\x7a\x23\x74\xab\x64\x89\xc5\xd4\xd9\xb7\xc5\x84\x3f\xee\x7a\xd7\x96\x49\xde\x68\x87\x30\xee\x36\x79\xf0\xfb\x4b\xff\xbb\xf9\xd0\xca\x52\x8a\x5d\x8c\x31\xc7\xb6\x04\x9e\x14\x2e\xb8\xaf\x13\x90\x5c\x47\x98\xdf\xa6\x42\x92\xe1\x42\xf7\x51\x4b\x8d\x40\xbe\x94\xa0\xbd\x93\xa3\x9f\x3e\xc7\x26\x5b\x48\x0b\xd1\x2b\x33\x76\xc6\xd9\x4b\x28\xe3\x02\x7e\x73\x19\x3b\x67\x79\x33\x76\x96\xf9\x34\x13\x40\x72\x5c\xbe\x08\xdf\xf6\x5b\x0a\xbe\x9e\xe4\xec\x11\x3a\x3c\x3d\xc8\x7a\xdf\xd7\x31\x0a\xbd\xe3\x9f\xfd\x85\xdf\x28\x94\x43\xcd\x8d\x51\x14\xa4\x70\x6a\x8c\xe7\x2d\xb0\x3e\x82\x5e\xec\xb3\x5e\x96\x02\x94\x24\x3c\xd7\x8b\xbc\x9f\xa4\xa5\xaf\xa6\xde\x5e\x39\x44\x25\xc6\x3e\xeb\x4d\x07\x6a\x6f\xf8\x3e\xcd\xb3\x77\xb4\x7c\x52\xa2\xc7\xf7\x1e\x95\x18\xa5\xe3\xba\x80\x1d\x42\x47\xb7\x80\x29\xa1\xa3\xb8\xec\xd1\x07\xb5\x7c\x32\x6a\xa2\x42\x1a\xa8\xf0\xfe\x2b\x0a\x44\x5f\x63\xa4\xd8\x79\xe9\x1a\xe2\x5d\xcb\x65\xec\x38\x9a\x64\x8c\x65\xc9\x94\xa7\x6f\xf6\xcd\x92\x4d\x45\xb9\x58\xf3\xf2\x55\xb7\x2a\x19\xb7\x18\x31\xca\xe9\x92\xd8\xb2\x0b\xc3\xd4\xa7\xbf\x0e\x72\xfb\x71\x08\x99\x2d\xc0\x75\x89\x27\xc8\x5f\xea\x34\xcc\x71\xdc\x5c\x03\xc9\xf0\xe4\x57\x6e\x3b\x6f\xc3\xc1\xee\x26\x4e\x78\xf6\xbb\xa8\x07\xb9\x40\xe2\xef\x96\x5e\xbe\x1c\xc6\x72\x18\x02\x07\x6a\xeb\xdc\xe1\x2b\xad\x6c\xa3\xfe\x53\x29\x93\xf7\x2b\x93\xff\x6b\x0a\x7c\x71\xc1\x1c\x28\xa5\x72\x29\x95\x43\x45\xdf\xbb\x54\x56\x2b\x1f\xd9\xfe\xd2\x2f\x65\x71\x7c\xcd\x17\x59\x9c\xc2\x1c\xbb\x0a\xea\xcf\x13\x68\xba\xf7\x93\xfe\xd3\x0f\x2a\x70\x65\x88\x6d\xa8\x57\xfe\x15\x00\x00\xff\xff\x25\x41\xcf\xe8\xea\xa1\x00\x00")

func dataGjishoGladeBytes() ([]byte, error) {
	return bindataRead(
		_dataGjishoGlade,
		"data/gjisho.glade",
	)
}

func dataGjishoGlade() (*asset, error) {
	bytes, err := dataGjishoGladeBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/gjisho.glade", size: 41450, mode: os.FileMode(436), modTime: time.Unix(1593293437, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataGjishoGlade2 = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5d\xdd\x76\xe2\xb8\x96\xbe\xe7\x29\x34\xbe\x38\xab\x7b\x0d\xa4\x52\xa9\x3e\x67\xea\x74\x27\xf4\x90\x84\x24\xee\x22\x90\x01\xa7\x33\x75\xc5\x12\xf6\x06\x54\x31\x92\x47\x12\x10\x7a\xcd\x0b\xcd\x6b\xcc\x93\x9d\x65\xd9\x24\x18\xfc\x23\xdb\x50\x21\x29\xfa\xa6\x2b\x46\xfa\x24\x6d\x69\xff\x68\x6b\x4b\xfb\xf4\xf7\xa7\x89\x8b\x66\xc0\x05\x61\xf4\xcc\xf8\x78\x74\x6c\x20\xa0\x36\x73\x08\x1d\x9d\x19\xf7\xd6\x55\xed\xb3\xf1\x7b\xbd\x72\xfa\x6f\xb5\x1a\xba\x06\x0a\x1c\x4b\x70\xd0\x9c\xc8\x31\x1a\xb9\xd8\x01\xf4\xe9\xe8\xd3\x3f\x8e\x8e\x51\xa5\x62\x8d\x01\xdd\x9a\x16\x6a\x11\x1b\xa8\x00\xf4\xd3\xad\x69\xfd\x5c\xa9\x5c\x30\x6f\xc1\xc9\x68\x2c\xd1\x4f\xf6\xcf\xe8\xe4\xf8\xe4\x18\x99\x98\xa2\x3f\xd8\x98\x0a\x46\x2b\x95\x3b\xe0\x13\x22\xfc\xd6\x11\x11\x68\x0c\x1c\x06\x0b\x34\xe2\x98\x4a\x70\xaa\x68\xc8\x01\x10\x1b\x22\x7b\x8c\xf9\x08\xaa\x48\x32\x84\xe9\x02\x79\xc0\x05\xa3\x88\x0d\x24\x26\x94\xd0\x11\xc2\xc8\x66\xde\xa2\xc2\x86\x48\x8e\x89\x40\x82\x0d\xe5\x1c\x73\x40\x98\x3a\x08\x0b\xc1\x6c\xa2\xfa\xed\x30\x7b\x3a\x01\x2a\xb1\xf4\xdb\x1b\x12\x17\x04\xfa\x49\x8e\x01\x19\xbd\xb0\x86\xf1\xb3\x6a\xc4\x01\xec\x56\x08\x45\xfe\x6f\xcb\x9f\xd4\xa8\xd9\x54\x22\x0e\x42\x72\x62\xfb\x18\x55\x44\xa8\xed\x4e\x7d\x6a\x3d\xff\xec\x92\x09\x09\x5b\xf0\xab\xab\xc1\x8b\x8a\x64\x68\x2a\xa0\xaa\xfa\x59\x45\x13\xe6\x90\xa1\xff\x7f\x50\xc3\xf2\xa6\x03\x97\x88\x71\x15\x39\xc4\x87\x1e\x4c\x25\x54\x91\xf0\x3f\x2a\x5a\x56\xfd\x71\x7c\x60\x1c\x09\x70\xdd\x8a\xcd\x3c\x02\x02\xa9\xb1\xbe\xf4\x4e\x95\xf1\xbb\xee\xf9\x04\x95\x21\x89\x84\xff\x65\x3e\x66\x93\xe8\x48\x88\xa8\x0c\xa7\x9c\x12\x31\x06\x55\xc7\x61\x48\x30\xd5\xe2\x37\xb0\xa5\xff\xc5\x2f\x3e\x64\xae\xcb\xe6\xfe\xd0\x6c\x46\x1d\xe2\x8f\x48\xfc\x1a\x4c\x34\x1e\xb0\x19\xa8\xb1\x04\x73\x4b\x99\x24\x76\x40\x6e\x35\x01\xde\xcb\xac\x86\x3f\x89\x31\x76\x5d\x34\x80\x90\x60\xe0\x20\x42\x2b\xfe\xa7\xe5\x70\xb8\xdf\xbc\x90\x98\x4a\x82\x5d\xe4\x31\xae\xda\x5b\x1f\xe6\x51\xa5\x62\xdd\x34\x51\xaf\x73\x65\x3d\x34\xba\x4d\x64\xf6\xd0\x5d\xb7\xf3\xa7\x79\xd9\xbc\x44\x46\xa3\x87\xcc\x9e\x51\x45\x0f\xa6\x75\xd3\xb9\xb7\xd0\x43\xa3\xdb\x6d\xb4\xad\xaf\xa8\x73\x85\x1a\xed\xaf\xe8\x8b\xd9\xbe\xac\xa2\xe6\x7f\xdf\x75\x9b\xbd\x1e\xea\x74\x2b\xe6\xed\x5d\xcb\x6c\x5e\x56\x91\xd9\xbe\x68\xdd\x5f\x9a\xed\x6b\x74\x7e\x6f\xa1\x76\xc7\x42\x2d\xf3\xd6\xb4\x9a\x97\xc8\xea\x20\xbf\xc1\x10\xca\x6c\xf6\x7c\xb0\xdb\x66\xf7\xe2\xa6\xd1\xb6\x1a\xe7\x66\xcb\xb4\xbe\x56\x2b\x57\xa6\xd5\xf6\x31\xaf\x3a\x5d\xd4\x40\x77\x8d\xae\x65\x5e\xdc\xb7\x1a\x5d\x74\x77\xdf\xbd\xeb\xf4\x9a\xa8\xd1\xbe\x44\xed\x4e\xdb\x6c\x5f\x75\xcd\xf6\x75\xf3\xb6\xd9\xb6\x8e\x90\xd9\x46\xed\x0e\x6a\xfe\xd9\x6c\x5b\xa8\x77\xd3\x68\xb5\xfc\xa6\x2a\x8d\x7b\xeb\xa6\xd3\xf5\xfb\x87\x2e\x3a\x77\x5f\xbb\xe6\xf5\x8d\x85\x6e\x3a\xad\xcb\x66\xb7\x87\xce\x9b\xa8\x65\x36\xce\x5b\xcd\xa0\xa9\xf6\x57\x74\xd1\x6a\x98\xb7\x55\x74\xd9\xb8\x6d\x5c\x37\x55\xad\x8e\x75\xd3\xec\x56\xfc\x62\x41\xef\xd0\xc3\x4d\xd3\xff\xe4\xb7\xd7\x68\xa3\xc6\x85\x65\x76\xda\xfe\x30\x2e\x3a\x6d\xab\xdb\xb8\xb0\xaa\xc8\xea\x74\xad\xe7\xaa\x0f\x66\xaf\x59\x45\x8d\xae\xd9\xf3\x09\x72\xd5\xed\xdc\x56\x2b\x3e\x39\x3b\x57\x7e\x11\xb3\xed\xd7\x6b\x37\x03\x14\x9f\xd4\x28\x32\x23\x9d\xae\xfa\xfb\xbe\xd7\x7c\x06\x44\x97\xcd\x46\xcb\x6c\x5f\xf7\x90\xd9\x8e\x4c\xdf\x51\xa5\xd2\x98\xca\x31\xe3\xbf\xae\x0a\x05\x74\x4a\x30\xf5\x38\x99\xc0\xf1\xdf\x8f\xff\xf9\x9f\xa3\x09\x26\xee\x91\xcd\x26\xf5\x4a\xa5\x56\xab\x57\x4e\x09\x95\xc0\x87\xd8\x86\x7a\x05\xa1\x53\x0e\xff\x33\x25\x1c\x04\x72\xc9\xe0\xcc\x18\xc9\xc7\x7f\x37\x5e\xe4\xd9\xa7\xa3\x93\x13\xe3\x83\x2a\xe7\x4b\xb0\xe7\x9a\xb5\x90\xb3\x6a\x72\xe1\x01\xf2\xf9\xc5\x47\xde\x28\x45\xf1\x04\xd0\xf5\x1f\x44\x8c\x59\x7c\x01\x07\x84\xcd\x89\xa7\x98\xbd\x81\xfe\xc0\x1e\xa6\x20\x00\x39\x81\x74\xc0\x7c\x81\x86\x8c\xa3\xeb\x76\xe7\xb6\x79\x14\x8f\xf0\xc2\x43\xeb\xb2\x31\xbe\x3c\x56\xf4\x12\x5a\xf4\x5a\x22\xb0\x80\xad\x6d\x17\x0b\x71\x66\x5c\xcb\xc7\x3b\xe6\xb1\x19\xf0\x5b\xa0\x53\x03\x11\xe7\xcc\x98\x60\x42\xd5\x5f\x7e\x79\x84\x4e\x3d\xce\x3c\xe0\x72\x81\x7c\x0a\x9c\x19\x36\xa6\xfd\x21\xb3\xa7\xc2\xa8\x5f\x61\x57\xc0\xe9\x87\x65\x81\xb0\xbc\x3d\x26\xae\x13\xfc\x3b\xae\xbd\x73\xf6\x64\x2c\x7f\xdd\x44\x9f\x11\x41\x06\x2e\x18\x75\x8b\x4f\x37\xa0\x8b\x74\x27\xae\x0e\xe3\x64\x29\xf8\x8d\xfa\x0c\xb8\x24\x36\x76\x63\x2b\x46\xc6\x12\x3f\x9e\x5b\xe6\x80\x7b\x3e\x95\xd2\x07\x5b\x29\x59\x60\x6c\x19\xe3\xcb\x53\x8d\x83\x0d\x64\x06\xa2\xef\xc0\x10\x4f\x5d\x99\xaf\x36\x56\x4b\xb6\xef\xff\x61\xd4\xb1\xe7\x1d\xe1\x01\x9b\x4a\xdd\xda\x12\x9e\xa4\x81\x24\xc7\x54\xb8\x58\xe2\x81\x0b\x67\xc6\x02\x84\x51\x6f\xf8\x28\x21\x0f\xc5\x83\x9d\x7e\x08\xe8\x1b\xf9\xe6\x61\xfb\x91\xd0\x51\x7a\xa3\xf0\xe4\x61\xea\xa4\x2c\x82\xb8\x4a\x43\xe2\xba\xf9\x28\xe3\x31\x41\x82\x75\x73\x92\x34\x82\x8d\xee\x9e\x7e\x88\xb2\xc4\xda\x18\x37\xc7\xb7\xde\xa8\x98\x0e\x26\x8a\x21\x7d\xd6\xd4\x59\xe0\x2f\xbd\xfc\xb8\x59\x7c\xad\x87\x2b\xbd\x5b\xed\xd9\xe6\x42\x6f\x78\x9e\x4b\x6c\xc5\x35\x0f\x84\x3a\x6c\x1e\x88\x0b\xec\x79\xe1\x9f\x05\xe5\xc5\xfa\x44\xce\x80\x4a\x61\xd4\xaf\x2f\xbf\xf4\xcf\xef\x2d\xab\xd3\xee\x2b\x1d\xdd\xbf\x6d\xf4\xbe\xa0\xff\x45\xfe\xf7\x2f\xcd\xaf\x9b\x1f\x7b\x56\xf7\xfe\xc2\xba\xef\x36\xd5\xc7\xac\x66\x42\xbe\xe8\xcf\x89\x23\xc7\x46\xfd\xf3\xf1\xb1\x6e\x8d\x31\xf8\x22\xda\xa8\xff\x23\xa6\x8a\x20\x23\x8a\xdd\xb0\xc2\x40\x09\x85\x9a\xc7\x41\x88\x9a\x1a\x96\x81\xc6\x98\x3a\x2e\xf0\x33\x63\xae\x68\x16\xc8\x8d\x3b\xbf\x84\x81\xc4\x1c\x7b\x1e\x38\x67\x06\x65\x81\x9e\x5a\x07\x7c\x84\x45\x2a\xda\x17\x58\xa4\x40\xed\xa1\x5c\xd6\x10\xaf\x5d\x98\x01\x76\x81\x07\x8b\x4d\x00\xe6\xf6\xf8\xf9\xdb\x6e\xe5\x6d\x2e\x51\xa2\xc4\x9d\x62\xba\xbe\x6f\x47\x18\x75\xe1\x12\x07\x6a\x4a\x9d\x27\x62\x6c\x8c\x5f\x67\x6a\x92\x7a\xa0\x56\x72\xdf\x37\x81\x40\x48\xa3\xfe\x69\x73\x79\x96\x26\x56\x19\x82\xc5\xd5\xd5\x56\xc4\xe9\x14\xcb\x43\xb5\x52\x63\x2f\x3b\xfe\xad\xd0\x20\x9d\x0e\xf1\xb4\xe8\x29\xbe\x69\x52\xc9\x17\xab\x8c\x14\x7c\x48\x82\x29\x49\xa8\x0c\x62\x15\x85\x98\x60\x3e\x22\xb4\x2f\x24\xe6\xd2\xa8\xff\xbd\x04\x04\xf8\xf6\x42\x19\x00\xc9\xbc\x72\x00\x03\x26\x25\x9b\x14\xc4\xf0\xad\x7c\xcc\x17\x7d\x62\x3f\x9b\x6a\xe0\x10\x59\x1b\x12\xea\xd4\xc4\x62\x32\x60\x2e\xb1\x4b\x03\xfb\xa6\xe0\x2c\xb0\xe2\x34\x17\x78\x26\xa4\x00\x25\x29\x67\xb9\x00\x35\x15\xe1\xca\xba\x4e\xd5\x86\xd9\x6d\x04\x48\x35\x7b\x8c\xe9\x08\x9c\x8d\x26\x2e\x96\xdf\xb5\xc0\xe3\xec\xda\x28\xbd\xe2\x6c\xdc\x34\x8a\x6a\xd9\xbb\x69\x00\xd9\xb6\x6f\x5a\xed\x17\x0b\x33\x55\xcb\x2c\x47\x9f\x3a\xbc\x35\xfb\x38\xfa\x53\x32\xdd\xd2\x69\x56\x8a\x5e\xc5\x69\x55\x98\x4e\x29\x34\x4a\xa4\x4f\x1e\x6d\xd8\xb3\x39\x73\x5d\x70\x22\xb6\x7a\x56\xef\xb7\xa5\x18\x8b\x54\x1f\x0b\xd5\xe3\x01\xe6\x7d\xcf\x97\x64\x0b\xa3\x4e\x61\x06\x3c\x13\x27\xc2\xc5\xe0\x8c\xa0\xc6\x01\xdb\xe3\x18\x1e\xee\x82\x98\xba\x52\x34\x9d\x11\x74\x97\x45\x34\xd8\x39\xaf\xf2\xfd\x93\xc0\xdc\x63\xbe\xb6\x7a\x15\x45\x5b\x58\x42\x88\x31\x76\xd8\x3c\x34\x67\x29\xa3\x7a\x20\xa9\xc4\x89\x27\x50\x8b\x08\xe9\x5b\x6b\x11\x13\x5f\x4d\x4d\x0a\xc9\xb6\x42\xb6\x6d\x91\x2e\x96\x7c\xe0\x42\xe0\x47\x99\x30\x07\x8c\xfa\x80\xb3\x79\x0e\xb4\xc8\x32\xe6\x6c\x5e\x0b\xf0\x92\x97\x71\x97\xcd\x7b\xcf\x45\x34\x55\x5e\xb6\x66\x4a\x15\xcf\xd9\x00\xdf\x55\xb6\x17\x11\x33\xdb\x10\xed\x31\x4e\x96\xe8\x40\xf3\x8a\xf6\x78\xd2\xc4\x14\x7e\x43\xfe\xb2\x04\xfd\x97\xed\x2f\xd3\xf3\x13\xf4\xc0\xc3\x1c\x4b\xf6\x8a\x4e\x81\x37\x34\x19\x09\x2b\x76\x5b\x93\xb1\xb9\xf5\x7e\x55\xdf\x4c\xae\x2d\xb6\xae\x53\x26\xd3\xa0\xda\x85\x87\x25\x6f\xd5\x22\x46\x54\x1e\xab\x32\xc3\xb4\x79\x5d\x47\x4b\x4e\x33\x2d\xd9\x5f\x54\x7a\x28\xdb\x18\x4e\x1c\x46\x21\xdf\x51\x36\x6d\xe2\xe9\xd3\xc2\x03\x70\x03\x1b\x2d\xdc\xce\x7f\xc1\x14\x07\x5f\xdf\x92\x99\x86\x65\x18\x56\x21\xd2\x0b\xae\x16\x5d\x9a\x74\x36\x76\xc1\x40\x33\xec\x4e\xe1\xcc\xf8\x78\x74\x92\x6e\x5b\xf9\xd2\x53\xb7\xb5\x53\x21\x17\x2e\x64\xf6\x48\x4d\x47\xd4\xa9\x52\x7b\xc4\x14\x67\x77\x24\x13\x3f\xdb\x12\xd4\xf1\x53\xa0\x6d\xf8\x2a\xe2\x40\xf2\xf9\x2b\xe2\x10\x72\xf9\x2c\x50\xb6\xdf\x02\x65\x1b\xc7\x5b\xe4\xb4\x6f\xe4\x87\x66\xb5\xff\x38\xf9\xbc\x2f\xcc\xf6\x8d\x1c\xb8\x2d\x1b\x41\x7b\x7b\xf4\x42\x91\x57\xe7\x36\xa0\x92\x2f\x2e\x41\x62\xe2\x8a\xb7\xc7\x6d\xe5\x0f\x4a\x52\x60\xf4\x0f\x4b\x52\x40\xd4\x81\xc9\x47\x4d\xe9\xbb\x9d\x53\x93\x38\x9c\x39\xc7\x5e\xb9\x39\x7a\xc2\x2e\x19\xe9\x2b\x92\x75\x77\x4e\x78\xa8\x02\x35\x97\xd0\xc7\x15\x7f\xce\xf2\x7b\x4b\x7d\xd6\xf6\xe1\x14\x12\x68\x6a\xad\xd7\x9c\x60\xb1\x1f\x04\x5a\x36\x42\x66\xe8\xcf\x26\x45\xbe\x8b\x40\xf3\x5c\x6c\xc3\x98\xb9\x0e\xf0\xb4\xc3\xad\xef\xee\xc5\xcb\xe5\xdb\x4a\x5b\x1b\x45\xdd\x7d\xc5\x96\x40\x21\xab\x31\x61\xaa\x63\x49\xa0\xeb\x60\x88\x46\xbe\x4c\x18\x07\x93\x0e\x59\x42\xec\x4b\x5c\xcf\x5f\x25\xa4\x23\x21\x0e\x66\xea\x7d\xff\xb0\x8e\x20\x50\x4b\x3b\x22\xa6\x14\x09\xcb\x92\x31\xae\xfe\x77\x8a\x0c\x89\xf7\x9b\x6e\x8d\x2a\xdb\xa0\xcc\xe1\x00\xbd\xe8\x01\x7a\x81\x05\x21\xb1\xfd\xd8\x9b\x13\x69\x8f\x63\xa5\x4c\xd2\x80\xf6\xc2\x21\x36\x0e\xed\x32\x1b\xa8\xcc\x3e\x2a\x8e\x43\x38\x84\x17\x2d\x37\xe0\xfe\x3a\x30\xea\x4b\xbd\xa3\x96\xc5\x0f\xcd\xa7\x1a\xdb\xd8\xef\xce\xa7\x51\xdb\x20\xf8\xb4\xef\x2c\x5b\x60\x8b\xae\x15\xbc\xb2\xb5\x51\x66\x8c\xb4\x0c\x4c\xd1\xa0\x16\x7d\xe2\xc5\x13\x50\x23\x04\x25\xa9\xcb\x85\xc8\x17\x07\x54\xd4\xaf\x11\x87\x55\x28\x34\xe5\x19\x4c\x8b\x88\xa8\xc0\x21\x51\x5a\x97\x0b\xd3\x31\x0e\xac\x0c\x2d\xe3\xf0\x0a\x1f\x24\x45\x40\xb5\xe9\x8a\x52\x1c\x71\xba\xf5\xb7\x4c\xe1\x38\xc0\xb2\x54\x8e\xc3\x74\x03\x6f\x63\xcc\xed\x28\xe5\xe3\x47\x73\x4e\x24\xa1\x23\xb1\x8d\xb6\x72\x7b\xa9\x22\x60\xfa\xce\xfb\xf8\x6a\x4b\x67\x5b\x70\x5d\x66\xe9\xc9\x1f\x30\xd7\xc9\xf2\x36\x65\x20\xae\x1f\x0d\xfc\xf2\xe9\x9f\x91\xff\xf2\xc0\xeb\x9f\x19\xac\xd6\xc9\xf6\x74\x45\xca\x6b\x79\xbd\xa2\x35\xb6\xe0\x01\x4b\x03\xcc\xef\x0d\x4b\x43\xcb\x7d\xb0\x16\x01\xcb\xf6\x92\x45\x8b\xe7\x90\x32\xdb\x91\x48\xca\xc6\x52\xe7\x3e\x0f\x21\x73\xfe\x10\x92\x2a\xbf\x97\x3c\x0d\x6d\x2f\x64\x51\xee\xf3\xfb\x08\xda\x41\x54\x6c\x53\x54\x68\x9e\x0a\x46\xc0\xf6\x5e\x54\xbc\x73\x91\x90\xdb\x31\x91\x05\x98\x6a\x0d\xe1\x83\x31\x94\x8d\x78\x30\x86\xf6\x56\xc2\x69\x1e\x13\x46\xc0\xf6\x5e\xc2\x2d\x8d\x21\x7c\xb0\x85\xde\xb6\x28\x3a\xd8\x42\xf9\x01\x77\x25\x29\x3e\xed\x93\xa4\xd0\x9f\x27\x4d\x50\x3d\x40\x0d\xb0\xdd\x85\x98\x50\x75\x7d\x3a\xaf\xb9\xb1\x71\xf6\x4d\xa4\xcf\x52\x31\xd6\xcc\x43\x2e\xe4\xd7\x8a\x84\x3b\xb8\xd9\xb5\x89\x17\x4f\xc0\x83\x9b\x7d\x15\xac\x84\x9b\x3d\x72\x23\x54\xb9\x5e\xfc\x2f\xef\xd1\xf3\xbe\x7e\x4b\x34\x2f\x95\x15\xe6\xc6\x5d\xd1\x65\x20\xe1\xea\x65\xd1\x67\x32\x76\xd9\xbc\xf1\xf2\x7b\x8e\x28\xc3\xe7\xf6\xd4\xcc\x22\x7f\x59\x9c\x19\x2b\xe1\x67\xda\x56\xe0\x0f\xba\x7d\x4e\xde\xed\xb6\x19\x52\xd3\x73\x30\x08\x53\x4a\x1f\x2c\x98\x48\xa9\x38\x0b\x26\xd7\x22\xd2\x37\x5f\xbe\x94\x81\x7d\x33\x17\x01\x0e\xe6\xcf\x12\x4f\xe7\x09\x0d\x78\xc2\x13\xcf\x85\xdc\xaf\x67\x3c\x37\x72\xb0\xb1\xf6\xc8\xc6\x5a\xce\xe6\xc1\xcc\x4a\xc1\x3c\x98\x3d\x85\x30\x53\xcd\x9e\xe5\xc2\x3b\x58\x3e\x29\xa5\x0f\x96\x4f\xa4\x54\x9c\xe5\x93\x77\x1d\xe9\x1b\x3f\xcd\x92\xc8\xaf\x70\x6f\x68\x77\x51\xb8\x45\x23\x45\xb7\x15\x84\xab\x41\xc2\x5d\xbc\x36\xb7\xf3\xbb\x4c\xb9\xef\xd9\x94\xbd\xcc\x94\x62\x8b\xef\xe2\x32\x53\x43\x29\xdf\x73\xbc\x3f\xf7\x96\xf2\x5c\x39\x0a\x5e\x77\x57\x86\xd2\x00\xdb\x8f\xb1\xaf\xbd\x97\x1e\x4d\x5c\xe5\xbc\x4f\x65\x96\xdc\x50\xac\x57\xcf\xf7\x8a\xfc\x33\x4a\x64\xfb\x60\xbb\xc4\x7e\x8c\xec\x1c\x28\x9e\x91\x11\x96\x70\xae\x02\xd8\x77\xf0\xe0\x9e\x39\xc1\x23\xd8\xfb\x38\xf8\x35\x8c\x95\x67\x5c\x47\xac\xe6\x71\x98\x11\x36\x15\x39\x1e\x72\xdd\x83\x87\xd8\xf6\xe6\xb5\xcb\x55\x8e\x1d\x32\x3e\xc7\xdc\x39\x30\xed\x56\x98\xf6\x2a\xa0\xe6\x81\x6f\xe3\xf9\x96\xc2\x93\x7c\xab\x3c\xbb\xf5\x67\x0c\xf3\xf1\xac\xc5\x46\x23\x17\x56\x39\x77\x79\xaf\x29\xf2\xcb\x77\x60\xe0\x3d\x60\x40\xa9\x86\xbc\xca\x80\x51\x62\x1c\xf8\x2f\x86\xff\x66\x04\xe6\x35\x9f\x4e\x6f\x8c\x03\xb1\xfd\x18\xfa\xff\x80\x3a\x79\x57\xdb\x0e\x75\xee\x1b\xde\xde\xa4\xec\x13\xf5\xb7\x37\xe5\x1f\xf4\xcc\xf3\xba\xe5\x1e\x24\x23\x8a\xfc\x18\x71\x77\x2a\xf7\xc8\x60\x65\xfb\xb6\x29\x3d\x6e\x00\x3b\xc0\xa3\x5b\xbc\xdc\x12\xa3\x88\x84\xd0\x77\xe5\x24\xa7\x9e\xda\xf4\xca\xb3\x79\xdf\x76\x99\x80\xfe\x20\xd4\x3b\x89\x5d\xd6\x78\x10\x75\x53\xb7\x05\x4f\x37\xa7\x68\xb6\x42\xc2\xb6\xa0\x26\x2b\xa7\xc1\xb2\x34\xd7\xea\x50\x53\xf5\x96\xae\x33\x21\x49\x4f\xbd\x8e\x23\x21\x59\x1f\xe5\x49\xc7\x51\xe2\xc9\xe7\x42\x4f\xf4\xde\x02\x9d\xee\xdd\xba\x53\x55\xfa\x8c\xf6\xd5\xfe\x23\xe7\x2b\xbf\xe5\x52\xde\x79\x41\x1a\xc4\x20\xc7\x9a\x4f\x9c\xc4\x8a\xef\x6c\x91\x32\x0f\x68\x6d\x02\x74\xba\xdb\x45\xaa\xab\x39\xf5\x8c\xa1\xdc\x9b\x98\x12\x9a\x30\x23\x2b\xde\x80\x4d\xe5\x25\xc1\x2e\x1b\x85\xf9\xf0\x56\x3e\x84\x58\x39\x27\x6e\x73\x61\x0b\xf2\x57\x4a\x02\x9e\x8d\x5b\x3a\xcc\xc1\x09\x96\xc4\x86\xaa\x5c\x78\xd0\x1f\x13\x2a\x8d\xba\xa3\xba\x9c\x59\x41\x3d\xef\x04\x54\xf6\x87\x8c\xab\xdc\x90\x41\xb4\x44\x56\x3d\x8f\xb3\x11\xc7\x93\x70\xcd\xc5\xeb\xe1\x0d\xce\x08\x32\xb7\x1a\xf5\xe3\xa3\x8f\x47\x99\xe9\xf9\x9e\x93\xa7\xc6\x6a\xfe\x97\xd4\xd3\xff\xff\x7f\x1b\xd9\x55\xb3\xa1\x27\x13\x95\x94\x30\x36\xa9\x65\x46\xaa\xd7\x2c\xf0\x30\x89\xab\x51\x5f\xcd\xe2\xfa\x37\x57\xfe\x16\x9f\xc8\xf5\x6f\x23\xf9\x5b\x16\xa4\xcb\x46\x6c\x35\x2d\x14\xf1\x45\x51\x4d\x25\x61\xa6\x99\x33\x1c\x66\xc3\x0d\x79\x70\x42\x64\xed\xe9\xe3\x06\x53\x45\x55\xbe\x03\x2e\x48\xd8\xc8\xc5\x34\x26\x0e\x3c\x10\x67\x04\x52\x23\xb7\x61\x02\x06\xa1\x63\x32\x20\xb2\xad\x52\x8a\x66\x82\x70\x10\x1e\xa3\x02\x72\x74\x22\x30\x72\x55\x4e\x5d\x8a\xdd\x9a\xfa\xf3\xcc\x98\x0d\x56\xde\x79\xc8\x9b\x2b\x71\xb7\x29\x69\xd7\x0d\x55\x0f\xdb\x84\x8e\x12\xb6\x00\x09\xc3\x0b\xf3\xbb\x62\x0e\xd8\x48\x37\x13\x02\x13\x21\xfb\x3d\xff\xa2\x4f\xf3\xbb\x78\xc1\xa6\xb2\xaf\x5e\x97\x4c\x97\xf8\xf1\x6a\x37\xf9\xf1\xc3\x3c\xc7\x55\xb9\x50\x5e\x25\xcd\x43\xae\x2a\xbb\x4a\xba\x11\x4f\xa6\xad\xed\x28\x37\xe1\xb5\xb5\xf0\x6a\x42\x5a\x15\x03\x19\x3e\xe2\xbb\xd5\xcc\xb4\x19\xf1\x91\x28\x7c\x4b\x75\x8b\xea\x79\x2d\x49\x6d\x5c\xc6\xd9\x84\x1a\x69\x49\x6a\x8b\xaa\xf4\xbd\x91\xf9\x19\x99\x6c\x13\xe3\x38\xb7\xe9\x8d\xd0\xad\x92\x27\x16\x53\x67\xdf\x96\x10\xfe\xb8\xed\x5d\x5b\x2e\x79\xa3\x1d\xc2\xb8\xdd\x84\xb7\xef\x2f\x65\xed\xfa\x7b\x23\xa1\x14\xbb\x18\x63\x8e\x6d\x09\x3c\x2d\x5c\x70\x57\x27\x20\x85\x8e\x30\xbf\x4d\x85\x24\xc3\x85\xee\x43\x8c\x1a\x81\x7c\x19\x41\x7b\x27\x47\xbf\x7c\x4e\x4c\x10\x90\x15\xa2\x77\xc8\x32\x99\x64\x2f\xa1\x9c\x0b\xf8\xcd\x65\x99\x9c\x15\xcd\x32\x79\xc8\x01\x99\x02\x52\xe0\xf2\x45\xf4\xb6\x5f\x28\xf8\x7a\x92\xb3\x47\xe8\xf0\xec\x20\xeb\x5d\x5f\xc7\x28\xf5\xf6\x7c\xfe\x57\x69\xe3\x50\xf6\x35\x9f\x43\x59\x90\xd2\xe9\x1c\x9e\xb7\xc0\xfa\x08\x7a\xb1\xcf\x7a\x2f\xeb\xa3\x34\xe1\xb9\x5a\xe4\xfd\x24\xda\x7c\x35\xf5\xf6\xca\x21\x2a\x09\xf6\x59\x6f\x3a\x50\x7b\xc3\xf7\x69\x9e\xbd\xa3\xe5\x93\x11\x3d\xbe\xf3\xa8\xc4\x38\x1d\xd7\x05\xec\x10\x3a\xba\x05\x4c\x09\x1d\x25\x65\x3c\xde\xab\xe5\x93\x53\x13\x95\xd2\x40\xa5\xf7\x5f\x71\x20\xfa\x1a\x23\xc3\xce\xcb\xd6\x10\xef\x5a\x2e\x63\xc7\xd1\x24\x63\x22\x4b\x66\x3c\x7d\xb3\x6b\x96\x6c\x2a\xca\x25\x9a\x97\xaf\xba\x55\xc9\xb9\xc5\x48\x50\x4e\x97\xc4\x96\x5d\x18\x66\xbe\xd3\xb5\x97\xdb\x8f\x7d\xc8\xc6\x00\xae\x4b\x3c\x41\xfe\x52\xa7\x61\x8e\xe3\x16\x1a\x48\x8e\x27\xbf\x0a\xdb\x79\x6b\x0e\x76\x37\x75\xc2\xf3\xdf\x45\xdd\xcb\x05\x92\x7c\xb7\xf4\xf2\xe5\x30\x96\xc3\x10\x38\x50\x5b\xe7\x0e\xdf\xc1\xca\x36\xea\xbf\x1c\x64\xf2\x6e\x65\xf2\x7f\x4d\x81\x2f\x2e\x98\x03\x07\xa9\x7c\x90\xca\x91\xa2\xef\x5d\x2a\xab\x95\x8f\x6c\x7f\xe9\x1f\x64\x71\x72\xcd\x17\x59\x9c\xc1\x1c\xdb\x0a\xea\x2f\x12\x68\xba\xf3\x93\xfe\xd3\x0f\x2a\x70\x65\x88\x6d\xa8\x57\xfe\x15\x00\x00\xff\xff\x4a\x3c\xfa\xc1\x9e\xa0\x00\x00")

func dataGjishoGlade2Bytes() ([]byte, error) {
	return bindataRead(
		_dataGjishoGlade2,
		"data/gjisho.glade~",
	)
}

func dataGjishoGlade2() (*asset, error) {
	bytes, err := dataGjishoGlade2Bytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/gjisho.glade~", size: 41118, mode: os.FileMode(420), modTime: time.Unix(1593293437, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"data/gjisho.glade": dataGjishoGlade,
	"data/gjisho.glade~": dataGjishoGlade2,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"data": &bintree{nil, map[string]*bintree{
		"gjisho.glade": &bintree{dataGjishoGlade, map[string]*bintree{}},
		"gjisho.glade~": &bintree{dataGjishoGlade2, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

