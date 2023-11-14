package i18n

import (
	"embed"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales
var languageFile embed.FS

// TT 翻译模板
type TT struct {
	ID string      // 语言ID
	TD interface{} // 模板数据TemplateData
	PC int         // 复数PluralCount
}

var (
	languageMap map[string]*i18n.Localizer
	defaultLang = "en"
)

func Init() {
	// 创建新的I18n包和本地化管理器
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	ls, err := languageFile.ReadDir("locales")
	if err != nil {
		panic(err)
	}
	languageMap = make(map[string]*i18n.Localizer)
	for _, v := range ls {
		if _, err := bundle.LoadMessageFileFS(languageFile, "locales/"+v.Name()); err != nil {
			panic(err)
		}
		lang := strings.TrimSuffix(v.Name(), ".toml")
		languageMap[lang] = i18n.NewLocalizer(bundle, lang)
	}
}

// Trasnlate 翻译
// ctx: gin上下文
// content: 翻译内容
func Trasnlate(ctx *gin.Context, content *TT) string {
	lang := defaultLang
	langs := strings.Split(ctx.GetHeader("Accept-Language"), ";")[0]
	for _, l := range strings.Split(langs, ",") {
		if _, ok := languageMap[l]; ok {
			lang = l
			break
		}
	}

	if content.PC > 0 {
		return languageMap[lang].MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    content.ID,
				One:   content.ID,
				Other: content.ID,
			},
			PluralCount:  content.PC,
			TemplateData: content.TD,
		})
	}

	return languageMap[lang].MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    content.ID,
			Other: content.ID,
		},
		TemplateData: content.TD,
	})
}

// DefaultTrasnlate 默认翻译
// content: 翻译内容
func DefaultTrasnlate(content *TT) string {
	if content.PC > 0 {
		return languageMap[defaultLang].MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    content.ID,
				One:   content.ID,
				Other: content.ID,
			},
			PluralCount:  content.PC,
			TemplateData: content.TD,
		})
	}

	return languageMap[defaultLang].MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    content.ID,
			Other: content.ID,
		},
		TemplateData: content.TD,
	})
}
