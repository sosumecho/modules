package i18n

import (
	"bytes"
	"embed"
	"github.com/sosumecho/modules/exceptions"
	"github.com/sosumecho/modules/i18n/types"
	"github.com/sosumecho/modules/logger"
	"github.com/sosumecho/modules/utils"
	"go.uber.org/zap"
	"io/fs"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/text/language"
)

// Translator 翻译器
var (
	Translator  *I18N
	DefaultLang = language.Chinese
	once        sync.Once
)

//func init() {
//	absDir := utils.GetAbsDir()
//	Translator = New().Load(fmt.Sprintf("%s/public/i18n", absDir), yaml.Unmarshal)
//}

// I18N 国际化
type I18N struct {
	Lang     language.Tag
	Messages map[language.Tag]map[string]string
	matcher  language.Matcher
	Tags     []language.Tag
	logger   *logger.Logger
	loader   types.Loader
}

// SetLang 设置国际化语言
func (i *I18N) SetLang(lang language.Tag) *I18N {
	i.Lang = lang
	return i
}

// Matcher  匹配器
func (i *I18N) Matcher() language.Matcher {
	return i.matcher
}

// Load 加载配置文件
func (i *I18N) Load() error {
	if i.loader != nil {
		messages, err := i.loader.Load()
		if err != nil {
			i.logger.Error("load i18n", zap.Error(err))
			return err
		}
		i.Messages = messages
		var tags = make([]language.Tag, 0, len(messages))
		for tag := range messages {
			tags = append(tags, tag)
		}
		if len(tags) > 0 {
			i.matcher = language.NewMatcher(tags)
		}
		return nil
	}
	return exceptions.InvalidI18NLoader
}

// LoadByFs 使用 fs加载
func (i *I18N) LoadByFs(dir embed.FS, unmarshal types.Unmarshal, allowExts []string) *I18N {
	_ = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fileInfo := strings.Split(d.Name(), ".")
		filename := fileInfo[0]
		ext := fileInfo[1]
		if !utils.IsStrInArr(allowExts, ext) {
			return nil
		}
		content, err := dir.ReadFile(path)
		if err != nil {
			return err
		}
		result := make(map[string]string)
		err = unmarshal(content, &result)
		if err != nil {
			return err
		}
		// 解析语言为tag
		t, err := language.Parse(filename)
		if err != nil {
			return err
		}
		i.Messages[t] = result
		i.Tags = append(i.Tags, t)
		return nil
	})
	i.matcher = language.NewMatcher(i.Tags)
	return i
}

// Tr 翻译
func (i *I18N) Tr(tags []language.Tag, key string, params map[string]interface{}) string {
	var rs bytes.Buffer
	for _, tag := range tags {
		if tmpl, exist := i.Messages[tag][key]; exist {
			t, err := template.New(key).Parse(tmpl)
			if err != nil {
				i.logger.Debug("解析多语言模板失败, ", zap.Error(err))
				//return key
				continue
			}
			err = t.Execute(&rs, params)
			if err != nil {
				i.logger.Debug("执行解析多语言失败, ", zap.Error(err))
				//return key
				continue
			}
			return rs.String()
		}
	}
	return key
}

func (i *I18N) SetLoader(loader types.Loader) *I18N {
	i.loader = loader
	return i
}

func (i *I18N) GetRootLanguages(lang string, defaultLang language.Tag) []language.Tag {
	var (
		rs      = make([]language.Tag, 0)
		tagDict = make(map[string]struct{})
	)

	acceptLanguages, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		return []language.Tag{defaultLang}
	}
	for _, item := range acceptLanguages {
		rootLanguage := i.GetRootLanguage(item.String(), defaultLang)
		if _, ok := tagDict[rootLanguage.String()]; ok {
			continue
		}
		rs = append(rs, rootLanguage)
		tagDict[rootLanguage.String()] = struct{}{}
	}
	return rs
}

func (i *I18N) GetRootLanguage(lang string, defaultLang language.Tag) language.Tag {
	languageTag, err := language.Parse(lang)
	if err != nil {
		return defaultLang
	}
	for {
		tag := languageTag.Parent()
		if tag.IsRoot() {
			break
		}
		languageTag = tag
	}
	return languageTag
}

// New 创建新的翻译器
func New(logf *logger.Logger) *I18N {
	if Translator == nil {
		once.Do(func() {
			Translator = &I18N{
				Lang:     DefaultLang,
				Messages: make(map[language.Tag]map[string]string),
				Tags:     make([]language.Tag, 0),
				logger:   logf,
			}
		})
	}
	return Translator
}
