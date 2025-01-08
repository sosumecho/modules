package upload

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/panjf2000/ants/v2"
	"github.com/sosumecho/modules/drivers/pool"
	"github.com/sosumecho/modules/exception"
	"github.com/sosumecho/modules/i18n"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/response"
	"github.com/sosumecho/modules/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// ResponseFormatter 响应格式
type ResponseFormatter func(c *gin.Context, group, path, name string)
type ResponseFormatter1 func(c *gin.Context, group, path, name string) (string, error)
type Formatter interface {
	RespFormat(c *gin.Context, locale *i18n.I18N, group string, domain string, fileInfo *FileInfo)
	Format(group string, domain string, fileInfo *FileInfo) URLFormatter
	DryFormat(group string, domain string, fileInfo *FileInfo) string
	RespFormatMulti(c *gin.Context, locale *i18n.I18N, group string, domain string, filenames map[string]*FileInfo)
	FormatMulti(group string, domain string, filenames map[string]*FileInfo) ([]URLFormatter, error)
}

type URLFormatter interface {
	GetURL() string
}

type SaveHandler func(fileInfo *FileInfo, checkOnly bool) (*FileInfo, bool, error)

type ThumbHandler func(group, domain string, fileInfo *FileInfo, size int, suffix string, formatter Formatter) error

const (
	Small = 180
	Large = 480

	SmallSuffix = "small"
	LargeSuffix = "large"
)

type Conf struct {
	AllowTypes []string
	MaxSize    string
	Domain     string
}

type Uploader struct {
	Formatter    Formatter
	key          string
	Conf         *Conf
	saveHandler  SaveHandler
	thumbHandler ThumbHandler
	locale       *i18n.I18N
	params       *Params
	thumbTypes   []string
	thumbPool    *ants.Pool
	logger       *logger.Logger
}

func NewUploader(conf *Conf, locale *i18n.I18N, log *logger.Logger) *Uploader {
	return &Uploader{
		Conf:   conf,
		key:    "file",
		locale: locale,
		thumbTypes: []string{
			"jpg",
			"png",
			"ico",
			"jpeg",
		},
		thumbPool: pool.NewPool(200),
		logger:    log,
	}
}

func (u *Uploader) SetKey(key string) *Uploader {
	u.key = key
	return u
}

func (u *Uploader) SetFormatter(formatter Formatter) *Uploader {
	u.Formatter = formatter
	return u
}

func (u *Uploader) SetSaveHandler(saveHandler SaveHandler) *Uploader {
	u.saveHandler = saveHandler
	return u
}

func (u *Uploader) SetThumbHandler(thumbHandler ThumbHandler) *Uploader {
	u.thumbHandler = thumbHandler
	return u
}

func (u *Uploader) Upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := response.New(ctx, u.locale, u.logger)
		err := u.validateParams(ctx)
		if err != nil {
			resp.Fail(exception.NewParamsError(err))
			return
		}

		urlMap := make(uploadFileInfos, 0, len(u.params.File.File[u.key]))
		isMulti := false
		if len(u.params.File.File[u.key]) > 1 {
			isMulti = true
		}
		var (
			wg      sync.WaitGroup
			thumbWg sync.WaitGroup
			rsWg    sync.WaitGroup
		)
		ch := make(chan FileInfo, 10)
		rsChan := make(chan FileInfo, 100)

		for _, file := range u.params.File.File[u.key] {
			err = u.Check(file)
			if err != nil {
				resp.Fail(exception.NewParamsError(err))
				return
			}
			wg.Add(1)
			go func(file *multipart.FileHeader) {
				defer wg.Done()
				uploadInfo, err := u.saveFile(ctx, file)
				if err != nil {
					return
				}
				//log.Logger().Info(
				//	"upload",
				//	zap.String("origin_name", uploadInfo.OriginName),
				//	zap.String("url", uploadInfo.URL),
				//	zap.String("name", uploadInfo.Filename),
				//	zap.String("path", uploadInfo.FilePath),
				//)
				ch <- *uploadInfo
			}(file)
		}

		thumbWg.Add(1)
		go func() {
			defer thumbWg.Done()
			for item := range ch {
				rsChan <- item
				if !utils.InArray(u.thumbTypes, item.Ext) {
					continue
				}
				if u.thumbHandler != nil {
					_ = u.thumbPool.Submit(func(file *FileInfo) func() {
						return func() {
							_ = u.thumbHandler(u.params.Group, u.Conf.Domain, file, Small, SmallSuffix, u.Formatter)
						}
					}(&item))
					_ = u.thumbPool.Submit(func(file *FileInfo) func() {
						return func() {
							_ = u.thumbHandler(u.params.Group, u.Conf.Domain, file, Large, LargeSuffix, u.Formatter)
						}
					}(&item))
				}
			}
		}()

		wg.Wait()
		close(ch)
		rsWg.Add(1)
		go func() {
			defer rsWg.Done()
			for item := range rsChan {
				//log.Logger().Info(
				//	"upload rs",
				//	zap.String("origin_name", item.OriginName),
				//	zap.String("url", item.URL),
				//	zap.String("name", item.Filename),
				//	zap.String("path", item.FilePath),
				//)
				fileInfo := item
				urlMap = append(urlMap, &fileInfo)
			}
		}()
		thumbWg.Wait()
		close(rsChan)

		rsWg.Wait()

		if len(urlMap) == 0 {
			resp.Fail(exception.NewParamsError(nil))
			return
		}

		if u.Formatter == nil {
			if isMulti {
				resp.Data(urlMap.FilePaths())
				return
			}
			resp.Data(urlMap[0].FilePath)
		} else if isMulti {
			u.Formatter.RespFormatMulti(ctx, u.locale, u.params.Group, u.Conf.Domain, urlMap.FileMap())
		} else {
			u.Formatter.RespFormat(ctx, u.locale, u.params.Group, u.Conf.Domain, urlMap[0])
		}
		return
	}
}

func (u *Uploader) validateParams(ctx *gin.Context) error {
	group := ctx.PostForm("group")
	if group == "" {
		group = "upload"
	}
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	u.params = &Params{
		Group: group,
		File:  form,
	}
	return nil
}

//func Upload(
//	key string,
//	locale *i18n.I18N,
//	saveHandler SaveHandler,
//	formatter Formatter,
//	domain string,
//	maxSize string,
//	allowTypes ...string,
//) gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		var uploadInfo *UploadFileInfo
//		urlMap := make(uploadFileInfos, 0, len(params.File.File[key]))
//		isMulti := false
//		if len(params.File.File[key]) > 1 {
//			isMulti = true
//		}
//		var (
//			wg      sync.WaitGroup
//			thumbWg sync.WaitGroup
//			rsWg    sync.WaitGroup
//		)
//		ch := make(chan UploadFileInfo, 10)
//		rsChan := make(chan UploadFileInfo, 100)
//		for _, file := range params.File.File[key] {
//			wg.Add(1)
//			go func() {
//				defer wg.Done()
//				uploadInfo, err = saveFile(c, locale, domain, file, params.Group, saveHandler, formatter, maxSize, allowTypes...)
//				if err != nil {
//					return
//				}
//				ch <- *uploadInfo
//			}()
//			//urlMap = append(urlMap, *uploadInfo)
//		}
//
//		thumbWg.Add(1)
//		go func() {
//			defer thumbWg.Done()
//			for item := range ch {
//				thumbWg.Add(1)
//				go func(item UploadFileInfo) {
//					defer thumbWg.Done()
//					//		err = thumbImage(&item)
//					//		if err == nil {
//					rsChan <- item
//					//		}
//				}(item)
//			}
//		}()
//
//		wg.Wait()
//		close(ch)
//		rsWg.Add(1)
//		go func() {
//			defer rsWg.Done()
//			for item := range rsChan {
//				urlMap = append(urlMap, item)
//			}
//		}()
//		thumbWg.Wait()
//		close(rsChan)
//
//		rsWg.Wait()
//
//		if len(urlMap) == 0 {
//			response.New(c, locale).Fail(exception.NewParamsError(nil))
//			return
//		}
//
//		if formatter == nil {
//			if isMulti {
//				response.New(c, locale).Data(urlMap.FilePaths())
//				return
//			}
//			response.New(c, locale).Data(uploadInfo.FilePath)
//		} else if isMulti {
//			formatter.FormatMulti(c, locale, group, domain, urlMap.FileMap())
//		} else {
//			formatter.Format(c, locale, group, domain, uploadInfo)
//		}
//		return
//	}
//}

func (u *Uploader) getOriginFilename(f *multipart.FileHeader) string {
	filename := f.Filename
	suffix := path.Ext(filename)
	return strings.TrimSuffix(filename, suffix)
}

type FileInfo struct {
	Filename   string
	FilePath   string
	Digest     string
	Ext        string
	OriginName string
	URL        string
	Thumb      struct {
		Large string
		Small string
	}
}

type uploadFileInfos []*FileInfo

func (u uploadFileInfos) FilePaths() []string {
	var rs = make([]string, 0, len(u))
	for _, item := range u {
		rs = append(rs, item.FilePath)
	}
	return rs
}

func (u uploadFileInfos) FileMap() map[string]*FileInfo {
	var rs = make(map[string]*FileInfo)
	for _, item := range u {
		rs[item.Filename] = item
	}
	return rs
}

func (u *Uploader) Check(file *multipart.FileHeader) error {
	//contentType := file.Header.Get("Content-Type")
	fileChecker := NewFileChecker(u.Conf)
	if err := fileChecker.Check(file, uint64(file.Size)); err != nil {
		return err
	}
	return nil
}

func (u *Uploader) GetFileType(file *multipart.FileHeader) string {
	//contentType := file.Header.Get("Content-Type")
	fileChecker := NewFileChecker(u.Conf)
	fileType, err := fileChecker.GetFileType(file)
	if err != nil {
		return ""
	}
	return fileType
}

func (u *Uploader) Md5(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	defer f.Close()
	if err != nil {
		return "", err
	}
	byt, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	digest := utils.Md5(byt)
	return digest, nil
}

func (u *Uploader) saveFile(
	ctx *gin.Context,
	file *multipart.FileHeader,
) (*FileInfo, error) {
	info, err := u.GenerateFileInfo(file)
	if err != nil {
		return nil, err
	}
	needSave := true
	if u.saveHandler != nil {
		url := u.Formatter.DryFormat(u.params.Group, u.Conf.Domain, info)
		info.URL = url
		info, needSave, err = u.saveHandler(info, true)
		if err != nil {
			return nil, err
		}
	}
	if needSave {
		err = ctx.SaveUploadedFile(file, info.Filename)
		if err != nil {
			return nil, err
		}
		_, _, _ = u.saveHandler(info, false)
	}
	return info, nil
}

func (u *Uploader) GenerateFileInfo(file *multipart.FileHeader) (*FileInfo, error) {
	digest, err := u.Md5(file)
	if err != nil {
		return nil, err
	}
	originName := u.getOriginFilename(file)
	uploadDir := utils.UploadDir(u.params.Group)

	id := utils.UniqueID()
	fileType := u.GetFileType(file)
	fileName := fmt.Sprintf("%s/%d.%s", uploadDir, id, fileType)
	filePath := fmt.Sprintf("/upload/%s/%d.%s", u.params.Group, id, fileType)
	filePath = strings.ReplaceAll(filePath, "//", "/")

	info := &FileInfo{
		Filename:   fileName,
		Digest:     digest,
		Ext:        fileType,
		FilePath:   filePath,
		OriginName: originName,
	}
	return info, nil
}

func ThumbImage(uploadInfo *FileInfo, size int, suffix string) (string, error) {
	srcImage, err := imaging.Open(uploadInfo.Filename, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}
	img := imaging.Resize(srcImage, size, 0, imaging.Lanczos)
	rs := strings.Split(uploadInfo.Filename, ".")
	thumbFilePath := uploadInfo.FilePath
	if len(rs) > 0 {
		thumbFilePath = strings.Join(append(rs[:len(rs)-1], fmt.Sprintf("_%s.%s", suffix, rs[len(rs)-1])), "")
	}
	thumbFileUrl := fmt.Sprintf("%s/%s", filepath.Dir(uploadInfo.FilePath), filepath.Base(thumbFilePath))
	_, err = os.Stat(thumbFilePath)
	if err == nil {
		return thumbFileUrl, nil
	}
	err = imaging.Save(img, thumbFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	// 取文件大小
	_, err = os.Stat(thumbFilePath)
	if err != nil {
		return "", err
	}

	return thumbFileUrl, nil
}

type Params struct {
	Group string
	File  *multipart.Form
}
