package loader

import (
	"embed"
	"github.com/sosumecho/modules/i18n/types"
	"github.com/sosumecho/modules/utils"
	"golang.org/x/text/language"
	"io/fs"

	"strings"
)

type EmbedLoader struct {
	fs        embed.FS
	unmarshal types.Unmarshal
	allowExt  []string
}

func (e EmbedLoader) Load() (map[language.Tag]map[string]string, error) {
	messages := make(map[language.Tag]map[string]string)
	err := fs.WalkDir(e.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fileInfo := strings.Split(d.Name(), ".")
		filename := fileInfo[0]
		ext := fileInfo[1]
		if !utils.IsStrInArr(e.allowExt, ext) {
			return nil
		}
		content, err := e.fs.ReadFile(path)
		if err != nil {
			return err
		}
		result := make(map[string]string)
		err = e.unmarshal(content, &result)
		if err != nil {
			return err
		}
		// 解析语言为tag
		t, err := language.Parse(filename)
		if err != nil {
			return err
		}
		messages[t] = result
		return nil
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func NewEmbedLoader(fs embed.FS, unmarshal types.Unmarshal, allowExt []string) *EmbedLoader {
	return &EmbedLoader{fs: fs, unmarshal: unmarshal, allowExt: allowExt}
}
