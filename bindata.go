// Code generated by go-bindata. DO NOT EDIT.
// sources:
// data/gjisho.glade
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

var _dataGjishoGlade = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5d\xdd\x76\xe2\xb8\x96\xbe\xe7\x29\x34\xbe\x38\xab\x7b\x0d\xa4\x52\xa9\x3e\x67\xea\x74\x27\xf4\x90\x84\x24\xee\x22\x90\x01\xa7\x33\x75\xc5\x12\xf6\x06\x54\x31\x92\x47\x12\x10\x7a\xcd\x0b\xcd\x6b\xcc\x93\x9d\x65\xd9\x24\x18\xfc\x23\xdb\x50\x21\x29\xfa\xa6\x2b\x46\xfa\x24\x6d\x69\xff\x68\x6b\x4b\xfb\xf4\xf7\xa7\x89\x8b\x66\xc0\x05\x61\xf4\xcc\xf8\x78\x74\x6c\x20\xa0\x36\x73\x08\x1d\x9d\x19\xf7\xd6\x55\xed\xb3\xf1\x7b\xbd\x72\xfa\x6f\xb5\x1a\xba\x06\x0a\x1c\x4b\x70\xd0\x9c\xc8\x31\x1a\xb9\xd8\x01\xf4\xe9\xe8\xd3\x3f\x8e\x8e\x51\xa5\x62\x8d\x01\xdd\x9a\x16\x6a\x11\x1b\xa8\x00\xf4\xd3\xad\x69\xfd\x5c\xa9\x5c\x30\x6f\xc1\xc9\x68\x2c\xd1\x4f\xf6\xcf\xe8\xe4\xf8\xe4\x18\x99\x98\xa2\x3f\xd8\x98\x0a\x46\x2b\x95\x3b\xe0\x13\x22\xfc\xd6\x11\x11\x68\x0c\x1c\x06\x0b\x34\xe2\x98\x4a\x70\xaa\x68\xc8\x01\x10\x1b\x22\x7b\x8c\xf9\x08\xaa\x48\x32\x84\xe9\x02\x79\xc0\x05\xa3\x88\x0d\x24\x26\x94\xd0\x11\xc2\xc8\x66\xde\xa2\xc2\x86\x48\x8e\x89\x40\x82\x0d\xe5\x1c\x73\x40\x98\x3a\x08\x0b\xc1\x6c\xa2\xfa\xed\x30\x7b\x3a\x01\x2a\xb1\xf4\xdb\x1b\x12\x17\x04\xfa\x49\x8e\x01\x19\xbd\xb0\x86\xf1\xb3\x6a\xc4\x01\xec\x56\x08\x45\xfe\x6f\xcb\x9f\xd4\xa8\xd9\x54\x22\x0e\x42\x72\x62\xfb\x18\x55\x44\xa8\xed\x4e\x7d\x6a\x3d\xff\xec\x92\x09\x09\x5b\xf0\xab\xab\xc1\x8b\x8a\x64\x68\x2a\xa0\xaa\xfa\x59\x45\x13\xe6\x90\xa1\xff\x7f\x50\xc3\xf2\xa6\x03\x97\x88\x71\x15\x39\xc4\x87\x1e\x4c\x25\x54\x91\xf0\x3f\x2a\x5a\x56\xfd\x71\x7c\x60\x1c\x09\x70\xdd\x8a\xcd\x3c\x02\x02\xa9\xb1\xbe\xf4\x4e\x95\xf1\xbb\xee\xf9\x04\x95\x21\x89\x84\xff\x65\x3e\x66\x93\xe8\x48\x88\xa8\x0c\xa7\x9c\x12\x31\x06\x55\xc7\x61\x48\x30\xd5\xe2\x37\xb0\xa5\xff\xc5\x2f\x3e\x64\xae\xcb\xe6\xfe\xd0\x6c\x46\x1d\xe2\x8f\x48\xfc\x1a\x4c\x34\x1e\xb0\x19\xa8\xb1\x04\x73\x4b\x99\x24\x76\x40\x6e\x35\x01\xde\xcb\xac\x86\x3f\x89\x31\x76\x5d\x34\x80\x90\x60\xe0\x20\x42\x2b\xfe\xa7\xe5\x70\xb8\xdf\xbc\x90\x98\x4a\x82\x5d\xe4\x31\xae\xda\x5b\x1f\xe6\x51\xa5\x62\xdd\x34\x51\xaf\x73\x65\x3d\x34\xba\x4d\x64\xf6\xd0\x5d\xb7\xf3\xa7\x79\xd9\xbc\x44\x46\xa3\x87\xcc\x9e\x51\x45\x0f\xa6\x75\xd3\xb9\xb7\xd0\x43\xa3\xdb\x6d\xb4\xad\xaf\xa8\x73\x85\x1a\xed\xaf\xe8\x8b\xd9\xbe\xac\xa2\xe6\x7f\xdf\x75\x9b\xbd\x1e\xea\x74\x2b\xe6\xed\x5d\xcb\x6c\x5e\x56\x91\xd9\xbe\x68\xdd\x5f\x9a\xed\x6b\x74\x7e\x6f\xa1\x76\xc7\x42\x2d\xf3\xd6\xb4\x9a\x97\xc8\xea\x20\xbf\xc1\x10\xca\x6c\xf6\x7c\xb0\xdb\x66\xf7\xe2\xa6\xd1\xb6\x1a\xe7\x66\xcb\xb4\xbe\x56\x2b\x57\xa6\xd5\xf6\x31\xaf\x3a\x5d\xd4\x40\x77\x8d\xae\x65\x5e\xdc\xb7\x1a\x5d\x74\x77\xdf\xbd\xeb\xf4\x9a\xa8\xd1\xbe\x44\xed\x4e\xdb\x6c\x5f\x75\xcd\xf6\x75\xf3\xb6\xd9\xb6\x8e\x90\xd9\x46\xed\x0e\x6a\xfe\xd9\x6c\x5b\xa8\x77\xd3\x68\xb5\xfc\xa6\x2a\x8d\x7b\xeb\xa6\xd3\xf5\xfb\x87\x2e\x3a\x77\x5f\xbb\xe6\xf5\x8d\x85\x6e\x3a\xad\xcb\x66\xb7\x87\xce\x9b\xa8\x65\x36\xce\x5b\xcd\xa0\xa9\xf6\x57\x74\xd1\x6a\x98\xb7\x55\x74\xd9\xb8\x6d\x5c\x37\x55\xad\x8e\x75\xd3\xec\x56\xfc\x62\x41\xef\xd0\xc3\x4d\xd3\xff\xe4\xb7\xd7\x68\xa3\xc6\x85\x65\x76\xda\xfe\x30\x2e\x3a\x6d\xab\xdb\xb8\xb0\xaa\xc8\xea\x74\xad\xe7\xaa\x0f\x66\xaf\x59\x45\x8d\xae\xd9\xf3\x09\x72\xd5\xed\xdc\x56\x2b\x3e\x39\x3b\x57\x7e\x11\xb3\xed\xd7\x6b\x37\x03\x14\x9f\xd4\x28\x32\x23\x9d\xae\xfa\xfb\xbe\xd7\x7c\x06\x44\x97\xcd\x46\xcb\x6c\x5f\xf7\x90\xd9\x8e\x4c\xdf\x51\xa5\xd2\x98\xca\x31\xe3\xbf\xae\x0a\x05\x74\x4a\x30\xf5\x38\x99\xc0\xf1\xdf\x8f\xff\xf9\x9f\xa3\x09\x26\xee\x91\xcd\x26\xf5\x4a\xa5\x56\xab\x57\x4e\x09\x95\xc0\x87\xd8\x86\x7a\x05\xa1\x53\x0e\xff\x33\x25\x1c\x04\x72\xc9\xe0\xcc\x18\xc9\xc7\x7f\x37\x5e\xe4\xd9\xa7\xa3\x93\x13\xe3\x83\x2a\xe7\x4b\xb0\xe7\x9a\xb5\x90\xb3\x6a\x72\xe1\x01\xf2\xf9\xc5\x47\xde\x28\x45\xf1\x04\xd0\xf5\x1f\x44\x8c\x59\x7c\x01\x07\x84\xcd\x89\xa7\x98\xbd\x81\xfe\xc0\x1e\xa6\x20\x00\x39\x81\x74\xc0\x7c\x81\x86\x8c\xa3\xeb\x76\xe7\xb6\x79\x14\x8f\xf0\xc2\x43\xeb\xb2\x31\xbe\x3c\x56\xf4\x12\x5a\xf4\x5a\x22\xb0\x80\xad\x6d\x17\x0b\x71\x66\x5c\xcb\xc7\x3b\xe6\xb1\x19\xf0\x5b\xa0\x53\x03\x11\xe7\xcc\x98\x60\x42\xd5\x5f\x7e\x79\x84\x4e\x3d\xce\x3c\xe0\x72\x81\x7c\x0a\x9c\x19\x36\xa6\xfd\x21\xb3\xa7\xc2\xa8\x5f\x61\x57\xc0\xe9\x87\x65\x81\xb0\xbc\x3d\x26\xae\x13\xfc\x3b\xae\xbd\x73\xf6\x64\x2c\x7f\xdd\x44\x9f\x11\x41\x06\x2e\x18\x75\x8b\x4f\x37\xa0\x8b\x74\x27\xae\x0e\xe3\x64\x29\xf8\x8d\xfa\x0c\xb8\x24\x36\x76\x63\x2b\x46\xc6\x12\x3f\x9e\x5b\xe6\x80\x7b\x3e\x95\xd2\x07\x5b\x29\x59\x60\x6c\x19\xe3\xcb\x53\x8d\x83\x0d\x64\x06\xa2\xef\xc0\x10\x4f\x5d\x99\xaf\x36\x56\x4b\xb6\xef\xff\x61\xd4\xb1\xe7\x1d\xe1\x01\x9b\x4a\xdd\xda\x12\x9e\xa4\x81\x24\xc7\x54\xb8\x58\xe2\x81\x0b\x67\xc6\x02\x84\x51\x6f\xf8\x28\x21\x0f\xc5\x83\x9d\x7e\x08\xe8\x1b\xf9\xe6\x61\xfb\x91\xd0\x51\x7a\xa3\xf0\xe4\x61\xea\xa4\x2c\x82\xb8\x4a\x43\xe2\xba\xf9\x28\xe3\x31\x41\x82\x75\x73\x92\x34\x82\x8d\xee\x9e\x7e\x88\xb2\xc4\xda\x18\x37\xc7\xb7\xde\xa8\x98\x0e\x26\x8a\x21\x7d\xd6\xd4\x59\xe0\x2f\xbd\xfc\xb8\x59\x7c\xad\x87\x2b\xbd\x5b\xed\xd9\xe6\x42\x6f\x78\x9e\x4b\x6c\xc5\x35\x0f\x84\x3a\x6c\x1e\x88\x0b\xec\x79\xe1\x9f\x05\xe5\xc5\xfa\x44\xce\x80\x4a\x61\xd4\xaf\x2f\xbf\xf4\xcf\xef\x2d\xab\xd3\xee\x2b\x1d\xdd\xbf\x6d\xf4\xbe\xa0\xff\x45\xfe\xf7\x2f\xcd\xaf\x9b\x1f\x7b\x56\xf7\xfe\xc2\xba\xef\x36\xd5\xc7\xac\x66\x42\xbe\xe8\xcf\x89\x23\xc7\x46\xfd\xf3\xf1\xb1\x6e\x8d\x31\xf8\x22\xda\xa8\xff\x23\xa6\x8a\x20\x23\x8a\xdd\xb0\xc2\x40\x09\x85\x9a\xc7\x41\x88\x9a\x1a\x96\x81\xc6\x98\x3a\x2e\xf0\x33\x63\xae\x68\x16\xc8\x8d\x3b\xbf\x84\x81\xc4\x1c\x7b\x1e\x38\x67\x06\x65\x81\x9e\x5a\x07\x7c\x84\x45\x2a\xda\x17\x58\xa4\x40\xed\xa1\x5c\xd6\x10\xaf\x5d\x98\x01\x76\x81\x07\x8b\x4d\x00\xe6\xf6\xf8\xf9\xdb\x6e\xe5\x6d\x2e\x51\xa2\xc4\x9d\x62\xba\xbe\x6f\x47\x18\x75\xe1\x12\x07\x6a\x4a\x9d\x27\x62\x6c\x8c\x5f\x67\x6a\x92\x7a\xa0\x56\x72\xdf\x37\x81\x40\x48\xa3\xfe\x69\x73\x79\x96\x26\x56\x19\x82\xc5\xd5\xd5\x56\xc4\xe9\x14\xcb\x43\xb5\x52\x63\x2f\x3b\xfe\xad\xd0\x20\x9d\x0e\xf1\xb4\xe8\x29\xbe\x69\x52\xc9\x17\xab\x8c\x14\x7c\x48\x82\x29\x49\xa8\x0c\x62\x15\x85\x98\x60\x3e\x22\xb4\x2f\x24\xe6\xd2\xa8\xff\xbd\x04\x04\xf8\xf6\x42\x19\x00\xc9\xbc\x72\x00\x03\x26\x25\x9b\x14\xc4\xf0\xad\x7c\xcc\x17\x7d\x62\x3f\x9b\x6a\xe0\x10\x59\x1b\x12\xea\xd4\xc4\x62\x32\x60\x2e\xb1\x4b\x03\xfb\xa6\xe0\x2c\xb0\xe2\x34\x17\x78\x26\xa4\x00\x25\x29\x67\xb9\x00\x35\x15\xe1\xca\xba\x4e\xd5\x86\xd9\x6d\x04\x48\x35\x7b\x8c\xe9\x08\x9c\x8d\x26\x2e\x96\xdf\xb5\xc0\xe3\xec\xda\x28\xbd\xe2\x6c\xdc\x34\x8a\x6a\xd9\xbb\x69\x00\xd9\xb6\x6f\x5a\xed\x17\x0b\x33\x55\xcb\x2c\x47\x9f\x3a\xbc\x35\xfb\x38\xfa\x53\x32\xdd\xd2\x69\x56\x8a\x5e\xc5\x69\x55\x98\x4e\x29\x34\x4a\xa4\x4f\x1e\x6d\xd8\xb3\x39\x73\x5d\x70\x22\xb6\x7a\x56\xef\xb7\xa5\x18\x8b\x54\x1f\x0b\xd5\xe3\x01\xe6\x7d\xcf\x97\x64\x0b\xa3\x4e\x61\x06\x3c\x13\x27\xc2\xc5\xe0\x8c\xa0\xc6\x01\xdb\xe3\x18\x1e\xee\x82\x98\xba\x52\x34\x9d\x11\x74\x97\x45\x34\xd8\x39\xaf\xf2\xfd\x93\xc0\xdc\x63\xbe\xb6\x7a\x15\x45\x5b\x58\x42\x88\x31\x76\xd8\x3c\x34\x67\x29\xa3\x7a\x20\xa9\xc4\x89\x27\x50\x8b\x08\xe9\x5b\x6b\x11\x13\x5f\x4d\x4d\x0a\xc9\xb6\x42\xb6\x6d\x91\x2e\x96\x7c\xe0\x42\xe0\x47\x99\x30\x07\x8c\xfa\x80\xb3\x79\x0e\xb4\xc8\x32\xe6\x6c\x5e\x0b\xf0\x92\x97\x71\x97\xcd\x7b\xcf\x45\x34\x55\x5e\xb6\x66\x4a\x15\xcf\xd9\x00\xdf\x55\xb6\x17\x11\x33\xdb\x10\xed\x31\x4e\x96\xe8\x40\xf3\x8a\xf6\x78\xd2\xc4\x14\x7e\x43\xfe\xb2\x04\xfd\x97\xed\x2f\xd3\xf3\x13\xf4\xc0\xc3\x1c\x4b\xf6\x8a\x4e\x81\x37\x34\x19\x09\x2b\x76\x5b\x93\xb1\xb9\xf5\x7e\x55\xdf\x4c\xae\x2d\xb6\xae\x53\x26\xd3\xa0\xda\x85\x87\x25\x6f\xd5\x22\x46\x54\x1e\xab\x32\xc3\xb4\x79\x5d\x47\x4b\x4e\x33\x2d\xd9\x5f\x54\x7a\x28\xdb\x18\x4e\x1c\x46\x21\xdf\x51\x36\x6d\xe2\xe9\xd3\xc2\x03\x70\x03\x1b\x2d\xdc\xce\x7f\xc1\x14\x07\x5f\xdf\x92\x99\x86\x65\x18\x56\x21\xd2\x0b\xae\x16\x5d\x9a\x74\x36\x76\xc1\x40\x33\xec\x4e\xe1\xcc\xf8\x78\x74\x92\x6e\x5b\xf9\xd2\x53\xb7\xb5\x53\x21\x17\x2e\x64\xf6\x48\x4d\x47\xd4\xa9\x52\x7b\xc4\x14\x67\x77\x24\x13\x3f\xdb\x12\xd4\xf1\x53\xa0\x6d\xf8\x2a\xe2\x40\xf2\xf9\x2b\xe2\x10\x72\xf9\x2c\x50\xb6\xdf\x02\x65\x1b\xc7\x5b\xe4\xb4\x6f\xe4\x87\x66\xb5\xff\x38\xf9\xbc\x2f\xcc\xf6\x8d\x1c\xb8\x2d\x1b\x41\x7b\x7b\xf4\x42\x91\x57\xe7\x36\xa0\x92\x2f\x2e\x41\x62\xe2\x8a\xb7\xc7\x6d\xe5\x0f\x4a\x52\x60\xf4\x0f\x4b\x52\x40\xd4\x81\xc9\x47\x4d\xe9\xbb\x9d\x53\x93\x38\x9c\x39\xc7\x5e\xb9\x39\x7a\xc2\x2e\x19\xe9\x2b\x92\x75\x77\x4e\x78\xa8\x02\x35\x97\xd0\xc7\x15\x7f\xce\xf2\x7b\x4b\x7d\xd6\xf6\xe1\x14\x12\x68\x6a\xad\xd7\x9c\x60\xb1\x1f\x04\x5a\x36\x42\x66\xe8\xcf\x26\x45\xbe\x8b\x40\xf3\x5c\x6c\xc3\x98\xb9\x0e\xf0\xb4\xc3\xad\xef\xee\xc5\xcb\xe5\xdb\x4a\x5b\x1b\x45\xdd\x7d\xc5\x96\x40\x21\xab\x31\x61\xaa\x63\x49\xa0\xeb\x60\x88\x46\xbe\x4c\x18\x07\x93\x0e\x59\x42\xec\x4b\x5c\xcf\x5f\x25\xa4\x23\x21\x0e\x66\xea\x7d\xff\xb0\x8e\x20\x50\x4b\x3b\x22\xa6\x14\x09\xcb\x92\x31\xae\xfe\x77\x8a\x0c\x89\xf7\x9b\x6e\x8d\x2a\xdb\xa0\xcc\xe1\x00\xbd\xe8\x01\x7a\x81\x05\x21\xb1\xfd\xd8\x9b\x13\x69\x8f\x63\xa5\x4c\xd2\x80\xf6\xc2\x21\x36\x0e\xed\x32\x1b\xa8\xcc\x3e\x2a\x8e\x43\x38\x84\x17\x2d\x37\xe0\xfe\x3a\x30\xea\x4b\xbd\xa3\x96\xc5\x0f\xcd\xa7\x1a\xdb\xd8\xef\xce\xa7\x51\xdb\x20\xf8\xb4\xef\x2c\x5b\x60\x8b\xae\x15\xbc\xb2\xb5\x51\x66\x8c\xb4\x0c\x4c\xd1\xa0\x16\x7d\xe2\xc5\x13\x50\x23\x04\x25\xa9\xcb\x85\xc8\x17\x07\x54\xd4\xaf\x11\x87\x55\x28\x34\xe5\x19\x4c\x8b\x88\xa8\xc0\x21\x51\x5a\x97\x0b\xd3\x31\x0e\xac\x0c\x2d\xe3\xf0\x0a\x1f\x24\x45\x40\xb5\xe9\x8a\x52\x1c\x71\xba\xf5\xb7\x4c\xe1\x38\xc0\xb2\x54\x8e\xc3\x74\x03\x6f\x63\xcc\xed\x28\xe5\xe3\x47\x73\x4e\x24\xa1\x23\xb1\x8d\xb6\x72\x7b\xa9\x22\x60\xfa\xce\xfb\xf8\x6a\x4b\x67\x5b\x70\x5d\x66\xe9\xc9\x1f\x30\xd7\xc9\xf2\x36\x65\x20\xae\x1f\x0d\xfc\xf2\xe9\x9f\x91\xff\xf2\xc0\xeb\x9f\x19\xac\xd6\xc9\xf6\x74\x45\xca\x6b\x79\xbd\xa2\x35\xb6\xe0\x01\x4b\x03\xcc\xef\x0d\x4b\x43\xcb\x7d\xb0\x16\x01\xcb\xf6\x92\x45\x8b\xe7\x90\x32\xdb\x91\x48\xca\xc6\x52\xe7\x3e\x0f\x21\x73\xfe\x10\x92\x2a\xbf\x97\x3c\x0d\x6d\x2f\x64\x51\xee\xf3\xfb\x08\xda\x41\x54\x6c\x53\x54\x68\x9e\x0a\x46\xc0\xf6\x5e\x54\xbc\x73\x91\x90\xdb\x31\x91\x05\x98\x6a\x0d\xe1\x83\x31\x94\x8d\x78\x30\x86\xf6\x56\xc2\x69\x1e\x13\x46\xc0\xf6\x5e\xc2\x2d\x8d\x21\x7c\xb0\x85\xde\xb6\x28\x3a\xd8\x42\xf9\x01\x77\x25\x29\x3e\xed\x93\xa4\xd0\x9f\x27\x4d\x50\x3d\x40\x0d\xb0\xdd\x85\x98\x50\x75\x7d\x3a\xaf\xb9\xb1\x71\xf6\x4d\xa4\xcf\x52\x31\xd6\xcc\x43\x2e\xe4\xd7\x8a\x84\x3b\xb8\xd9\xb5\x89\x17\x4f\xc0\x83\x9b\x7d\x15\xac\x84\x9b\x3d\x72\x23\x54\xb9\x5e\xfc\x2f\xef\xd1\xf3\xbe\x7e\x4b\x34\x2f\x95\x15\xe6\xc6\x5d\xd1\x65\x20\xe1\xea\x65\xd1\x67\x32\x76\xd9\xbc\xf1\xf2\x7b\x8e\x28\xc3\xe7\xf6\xd4\xcc\x22\x7f\x59\x9c\x19\x2b\xe1\x67\xda\x56\xe0\x0f\xba\x7d\x4e\xde\xed\xb6\x19\x52\xd3\x73\x30\x08\x53\x4a\x1f\x2c\x98\x48\xa9\x38\x0b\x26\xd7\x22\xd2\x37\x5f\xbe\x94\x81\x7d\x33\x17\x01\x0e\xe6\xcf\x12\x4f\xe7\x09\x0d\x78\xc2\x13\xcf\x85\xdc\xaf\x67\x3c\x37\x72\xb0\xb1\xf6\xc8\xc6\x5a\xce\xe6\xc1\xcc\x4a\xc1\x3c\x98\x3d\x85\x30\x53\xcd\x9e\xe5\xc2\x3b\x58\x3e\x29\xa5\x0f\x96\x4f\xa4\x54\x9c\xe5\x93\x77\x1d\xe9\x1b\x3f\xcd\x92\xc8\xaf\x70\x6f\x68\x77\x51\xb8\x45\x23\x45\xb7\x15\x84\xab\x41\xc2\x5d\xbc\x36\xb7\xf3\xbb\x4c\xb9\xef\xd9\x94\xbd\xcc\x94\x62\x8b\xef\xe2\x32\x53\x43\x29\xdf\x73\xbc\x3f\xf7\x96\xf2\x5c\x39\x0a\x5e\x77\x57\x86\xd2\x00\xdb\x8f\xb1\xaf\xbd\x97\x1e\x4d\x5c\xe5\xbc\x4f\x65\x96\xdc\x50\xac\x57\xcf\xf7\x8a\xfc\x33\x4a\x64\xfb\x60\xbb\xc4\x7e\x8c\xec\x1c\x28\x9e\x91\x11\x96\x70\xae\x02\xd8\x77\xf0\xe0\x9e\x39\xc1\x23\xd8\xfb\x38\xf8\x35\x8c\x95\x67\x5c\x47\xac\xe6\x71\x98\x11\x36\x15\x39\x1e\x72\xdd\x83\x87\xd8\xf6\xe6\xb5\xcb\x55\x8e\x1d\x32\x3e\xc7\xdc\x39\x30\xed\x56\x98\xf6\x2a\xa0\xe6\x81\x6f\xe3\xf9\x96\xc2\x93\x7c\xab\x3c\xbb\xf5\x67\x0c\xf3\xf1\xac\xc5\x46\x23\x17\x56\x39\x77\x79\xaf\x29\xf2\xcb\x77\x60\xe0\x3d\x60\x40\xa9\x86\xbc\xca\x80\x51\x62\x1c\xf8\x2f\x86\xff\x66\x04\xe6\x35\x9f\x4e\x6f\x8c\x03\xb1\xfd\x18\xfa\xff\x80\x3a\x79\x57\xdb\x0e\x75\xee\x1b\xde\xde\xa4\xec\x13\xf5\xb7\x37\xe5\x1f\xf4\xcc\xf3\xba\xe5\x1e\x24\x23\x8a\xfc\x18\x71\x77\x2a\xf7\xc8\x60\x65\xfb\xb6\x29\x3d\x6e\x00\x3b\xc0\xa3\x5b\xbc\xdc\x12\xa3\x88\x84\xd0\x77\xe5\x24\xa7\x9e\xda\xf4\xca\xb3\x79\xdf\x76\x99\x80\xfe\x20\xd4\x3b\x89\x5d\xd6\x78\x10\x75\x53\xb7\x05\x4f\x37\xa7\x68\xb6\x42\xc2\xb6\xa0\x26\x2b\xa7\xc1\xb2\x34\xd7\xea\x50\x53\xf5\x96\xae\x33\x21\x49\x4f\xbd\x8e\x23\x21\x59\x1f\xe5\x49\xc7\x51\xe2\xc9\xe7\x42\x4f\xf4\xde\x02\x9d\xee\xdd\xba\x53\x55\xfa\x8c\xf6\xd5\xfe\x23\xe7\x2b\xbf\xe5\x52\xde\x79\x41\x1a\xc4\x20\xc7\x9a\x4f\x9c\xc4\x8a\xef\x6c\x91\x32\x0f\x68\x6d\x02\x74\xba\xdb\x45\xaa\xab\x39\xf5\x8c\xa1\xdc\x9b\x98\x12\x9a\x30\x23\x2b\xde\x80\x4d\xe5\x25\xc1\x2e\x1b\x85\xf9\xf0\x56\x3e\x84\x58\x39\x27\x6e\x73\x61\x0b\xf2\x57\x4a\x02\x9e\x8d\x5b\x3a\xcc\xc1\x09\x96\xc4\x86\xaa\x5c\x78\xd0\x1f\x13\x2a\x8d\xba\xa3\xba\x9c\x59\x41\x3d\xef\x04\x54\xf6\x87\x8c\xab\xdc\x90\x41\xb4\x44\x56\x3d\x8f\xb3\x11\xc7\x93\x70\xcd\xc5\xeb\xe1\x0d\xce\x08\x32\xb7\x1a\xf5\xe3\xa3\x8f\x47\x99\xe9\xf9\x9e\x93\xa7\xc6\x6a\xfe\x97\xd4\xd3\xff\xff\x7f\x1b\xd9\x55\xb3\xa1\x27\x13\x95\x94\x30\x36\xa9\x65\x46\xaa\xd7\x2c\xf0\x30\x89\xab\x51\x5f\xcd\xe2\xfa\x37\x57\xfe\x16\x9f\xc8\xf5\x6f\x23\xf9\x5b\x16\xa4\xcb\x46\x6c\x35\x2d\x14\xf1\x45\x51\x4d\x25\x61\xa6\x99\x33\x1c\x66\xc3\x0d\x79\x70\x42\x64\xed\xe9\xe3\x06\x53\x45\x55\xbe\x03\x2e\x48\xd8\xc8\xc5\x34\x26\x0e\x3c\x10\x67\x04\x52\x23\xb7\x61\x02\x06\xa1\x63\x32\x20\xb2\xad\x52\x8a\x66\x82\x70\x10\x1e\xa3\x02\x72\x74\x22\x30\x72\x55\x4e\x5d\x8a\xdd\x9a\xfa\xf3\xcc\x98\x0d\x56\xde\x79\xc8\x9b\x2b\x71\xb7\x29\x69\xd7\x0d\x55\x0f\xdb\x84\x8e\x12\xb6\x00\x09\xc3\x0b\xf3\xbb\x62\x0e\xd8\x48\x37\x13\x02\x13\x21\xfb\x3d\xff\xa2\x4f\xf3\xbb\x78\xc1\xa6\xb2\xaf\x5e\x97\x4c\x97\xf8\xf1\x6a\x37\xf9\xf1\xc3\x3c\xc7\x55\xb9\x50\x5e\x25\xcd\x43\xae\x2a\xbb\x4a\xba\x11\x4f\xa6\xad\xed\x28\x37\xe1\xb5\xb5\xf0\x6a\x42\x5a\x15\x03\x19\x3e\xe2\xbb\xd5\xcc\xb4\x19\xf1\x91\x28\x7c\x4b\x75\x8b\xea\x79\x2d\x49\x6d\x5c\xc6\xd9\x84\x1a\x69\x49\x6a\x8b\xaa\xf4\xbd\x91\xf9\x19\x99\x6c\x13\xe3\x38\xb7\xe9\x8d\xd0\xad\x92\x27\x16\x53\x67\xdf\x96\x10\xfe\xb8\xed\x5d\x5b\x2e\x79\xa3\x1d\xc2\xb8\xdd\x84\xb7\xef\x2f\x65\xed\xfa\x7b\x23\xa1\x14\xbb\x18\x63\x8e\x6d\x09\x3c\x2d\x5c\x70\x57\x27\x20\x85\x8e\x30\xbf\x4d\x85\x24\xc3\x85\xee\x43\x8c\x1a\x81\x7c\x19\x41\x7b\x27\x47\xbf\x7c\x4e\x4c\x10\x90\x15\xa2\x77\xc8\x32\x99\x64\x2f\xa1\x9c\x0b\xf8\xcd\x65\x99\x9c\x15\xcd\x32\x79\xc8\x01\x99\x02\x52\xe0\xf2\x45\xf4\xb6\x5f\x28\xf8\x7a\x92\xb3\x47\xe8\xf0\xec\x20\xeb\x5d\x5f\xc7\x28\xf5\xf6\x7c\xfe\x57\x69\xe3\x50\xf6\x35\x9f\x43\x59\x90\xd2\xe9\x1c\x9e\xb7\xc0\xfa\x08\x7a\xb1\xcf\x7a\x2f\xeb\xa3\x34\xe1\xb9\x5a\xe4\xfd\x24\xda\x7c\x35\xf5\xf6\xca\x21\x2a\x09\xf6\x59\x6f\x3a\x50\x7b\xc3\xf7\x69\x9e\xbd\xa3\xe5\x93\x11\x3d\xbe\xf3\xa8\xc4\x38\x1d\xd7\x05\xec\x10\x3a\xba\x05\x4c\x09\x1d\x25\x65\x3c\xde\xab\xe5\x93\x53\x13\x95\xd2\x40\xa5\xf7\x5f\x71\x20\xfa\x1a\x23\xc3\xce\xcb\xd6\x10\xef\x5a\x2e\x63\xc7\xd1\x24\x63\x22\x4b\x66\x3c\x7d\xb3\x6b\x96\x6c\x2a\xca\x25\x9a\x97\xaf\xba\x55\xc9\xb9\xc5\x48\x50\x4e\x97\xc4\x96\x5d\x18\x66\xbe\xd3\xb5\x97\xdb\x8f\x7d\xc8\xc6\x00\xae\x4b\x3c\x41\xfe\x52\xa7\x61\x8e\xe3\x16\x1a\x48\x8e\x27\xbf\x0a\xdb\x79\x6b\x0e\x76\x37\x75\xc2\xf3\xdf\x45\xdd\xcb\x05\x92\x7c\xb7\xf4\xf2\xe5\x30\x96\xc3\x10\x38\x50\x5b\xe7\x0e\xdf\xc1\xca\x36\xea\xbf\x1c\x64\xf2\x6e\x65\xf2\x7f\x4d\x81\x2f\x2e\x98\x03\x07\xa9\x7c\x90\xca\x91\xa2\xef\x5d\x2a\xab\x95\x8f\x6c\x7f\xe9\x1f\x64\x71\x72\xcd\x17\x59\x9c\xc1\x1c\xdb\x0a\xea\x2f\x12\x68\xba\xf3\x93\xfe\xd3\x0f\x2a\x70\x65\x88\x6d\xa8\x57\xfe\x15\x00\x00\xff\xff\x4a\x3c\xfa\xc1\x9e\xa0\x00\x00")

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

	info := bindataFileInfo{name: "data/gjisho.glade", size: 41118, mode: os.FileMode(420), modTime: time.Unix(1592765607, 0)}
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
