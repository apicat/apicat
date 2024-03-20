package i18n

import (
	"fmt"
	"strings"

	"github.com/apicat/apicat/backend/i18n/lang"
	"github.com/apicat/apicat/backend/route/middleware/jwt"

	"github.com/gin-gonic/gin"
)

var langs = map[string]map[string]map[string]string{}

const defaultLang = "en-US"

func init() {
	langs[defaultLang] = lang.En
	langs["zh-CN"] = lang.Zh
}

func NewErr(text string, params ...string) error {
	return &Translation{
		s:      text,
		params: params,
	}
}

func NewTran(text string, params ...string) *Translation {
	return &Translation{
		s:      text,
		params: params,
	}
}

type Translation struct {
	s      string
	params []string
}

func (t *Translation) Error() string {
	return translate(defaultLang, t.s, t.params...)
}

func (t *Translation) Translate(ctx *gin.Context) string {
	if t.s == "" {
		return ""
	}

	userLang := defaultLang
	if jwt.GetUser(ctx) != nil {
		userLang = jwt.GetUser(ctx).Language
	} else {
		langs := strings.Split(ctx.GetHeader("Accept-Language"), ";")[0]
		if len(langs) > 0 {
			userLang = strings.Split(langs, ",")[0]
		}
	}
	return translate(userLang, t.s, t.params...)
}

func translate(lang string, key string, params ...string) string {
	indexs := strings.Split(key, ".")
	if len(indexs) != 2 {
		return key
	}
	if _, ok := langs[lang][indexs[0]]; !ok {
		return key
	}
	if _, ok := langs[lang][indexs[0]][indexs[1]]; !ok {
		return key
	}
	if len(params) > 0 {
		values := make([]interface{}, 0)
		for _, v := range params {
			values = append(values, v)
		}
		return fmt.Sprintf(langs[lang][indexs[0]][indexs[1]], values...)
	}
	return langs[lang][indexs[0]][indexs[1]]
}
