package loader

import (
	"github.com/sosumecho/modules/i18n/types"
	"github.com/sosumecho/modules/utils"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

type FileLoader struct {
	path      string
	unmarshal types.Unmarshal
	allowExt  []string
}

func NewFileLoader(path string, unmarshal types.Unmarshal, allowExt []string) *FileLoader {
	return &FileLoader{path: path, unmarshal: unmarshal, allowExt: allowExt}
}

func (f *FileLoader) Load() (map[language.Tag]map[string]string, error) {
	messages := make(map[language.Tag]map[string]string)
	tags := make([]language.Tag, 0)
	err := filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		fileInfo := strings.Split(info.Name(), ".")
		filename := fileInfo[0]
		ext := fileInfo[1]
		if !utils.IsStrInArr(f.allowExt, ext) {
			return nil
		}
		content, _ := os.ReadFile(path)
		result := make(map[string]string)
		err = f.unmarshal(content, &result)

		if err != nil {
			return err
		}

		// 解析语言为tag
		t, err := language.Parse(filename)
		if err != nil {
			return err
		}
		messages[t] = result
		tags = append(tags, t)
		return err
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
