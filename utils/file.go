package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
)

// CreateFile 创建文件
func CreateFile(filename string) (*os.File, error) {
	dir := path.Dir(filename)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0644)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	if err = LockFile(file); err != nil {
		return nil, err
	}
	err = file.Truncate(0)
	return file, err
}

// Write 写文件
func Write(filename string, content string) (*os.File, error) {
	file, err := CreateFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return file, err
}

// LockFile 文件加锁
func LockFile(file *os.File) error {
	return Flock(int(file.Fd()), LOCK_EX|LOCK_NB)
}

func FileMd5(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()
	b := md5.New()
	io.Copy(b, f)
	return hex.EncodeToString(b.Sum(nil))
}
