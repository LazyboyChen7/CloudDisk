package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type Sha256Stream struct {
	_sha256 hash.Hash
}

func (obj *Sha256Stream) Update(data []byte) {
	if obj._sha256 == nil {
		obj._sha256 = sha256.New()
	}
	obj._sha256.Write(data)
}

func (obj *Sha256Stream) Sum() string {
	return hex.EncodeToString(obj._sha256.Sum([]byte("")))
}

func Sha256(data []byte) string {
	_sha256 := sha256.New()
	_sha256.Write(data)
	return hex.EncodeToString(_sha256.Sum([]byte("")))
}

func FileSha256(file *os.File) string {
	_sha256 := sha256.New()
	io.Copy(_sha256, file)
	return hex.EncodeToString(_sha256.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}