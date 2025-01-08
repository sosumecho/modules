package upload

import (
	"github.com/dustin/go-humanize"
	"github.com/sosumecho/modules/exceptions"
	"github.com/sosumecho/modules/utils"
	"mime/multipart"
	"strings"
)

type FileChecker struct {
	AllowTypes []string
	MaxSize    uint64
}

func NewFileChecker(conf *Conf) *FileChecker {
	size, err := humanize.ParseBytes(conf.MaxSize)
	if err != nil {
		size = 1024 * 1024 * 20
	}
	return &FileChecker{AllowTypes: conf.AllowTypes, MaxSize: size}
}

func (t *FileChecker) GetFileType(f *multipart.FileHeader) (string, error) {
	fileType := strings.Split(f.Filename, ".")
	return fileType[len(fileType)-1], nil
	//switch s {
	//case "image/png":
	//	return "png", nil
	//case "image/jpg":
	//	return "jpg", nil
	//case "image/jpeg":
	//	return "jpeg", nil
	//case "image/gif":
	//	return "gif", nil
	//case "video/mp4":
	//	return "mp4", nil
	//case "image/webp":
	//	return "webp", nil
	//case "video/quicktime":
	//	return "mov", nil
	//case "application/zip":
	//	return "zip", nil
	//default:
	//	return "", exceptions.InvalidUploadFileType
	//}
}

func (t *FileChecker) Check(f *multipart.FileHeader, size uint64) error {
	fileType, err := t.GetFileType(f)
	if err != nil {
		return err
	}
	if !utils.InArray(t.AllowTypes, fileType) {
		return exceptions.InvalidUploadFileType
	}

	if size > t.MaxSize {
		return exceptions.ExceedMaxFileSize
	}

	return nil
}
