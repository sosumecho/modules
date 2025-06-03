package utils

import (
	"github.com/spf13/cast"
	"log"
	"os"
	path2 "path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetAbsDir 得到程序运行的绝对路径
func GetAbsDir() string {
	workingDir, _ := os.Getwd()
	binPath, err := filepath.Abs(workingDir)
	if err != nil {
		log.Fatalln(err)
	}
	return binPath
}

// RuntimeDir 运行时目录
func RuntimeDir(path ...string) string {
	dirArr := make([]string, 0)
	dirArr = append(dirArr, GetAbsDir())
	dirArr = append(dirArr, "runtime")
	if len(path) > 0 {
		dirArr = append(dirArr, strings.Join(path, string(os.PathSeparator)))
	}
	dir := strings.Join(dirArr, string(os.PathSeparator))
	if !IsDir(dir) {
		log.Println(dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	return dir
}

// IsDir 是否是文件夹
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Mkdir(path []string) (err error) {
	dirArr := make([]string, 0)
	dirArr = append(dirArr, GetAbsDir())
	dirArr = append(dirArr, path...)
	dir := path2.Join(dirArr...)
	if !IsDir(dir) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	return err
}

func UploadDir(path ...string) string {
	cwd, _ := os.Getwd()
	var uploadDir = cwd + "/public/upload"
	if len(path) > 0 {
		uploadDir += "/" + strings.Join(path, cast.ToString(os.PathSeparator))
	}

	if !IsDir(uploadDir) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return ""
		}
	}
	return uploadDir
}

func CreateDirIfNotExist(dir string) error {
	if IsDir(dir) {
		return nil
	}
	return Mkdir([]string{dir})
}
