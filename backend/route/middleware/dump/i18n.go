package dump

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/apicat/apicat/backend/i18n"
	"github.com/apicat/apicat/backend/route/middleware/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var trans = map[string]*ut.Translator{}

func init() {
	uni := ut.New(en.New(), en.New(), zh.New())
	validater := binding.Validator.Engine().(*validator.Validate)

	if zhtran, found := uni.GetTranslator("zh"); found {
		trans["zh"] = &zhtran
		trans["zh-CN"] = &zhtran
		if err := zh_translations.RegisterDefaultTranslations(validater, zhtran); err != nil {
			slog.Info("zh translator register failed")
		}
	} else {
		slog.Info("zh translator not found")
	}

	if entran, found := uni.GetTranslator("en"); found {
		trans["en"] = &entran
		if err := en_translations.RegisterDefaultTranslations(validater, entran); err != nil {
			slog.Info("en translator register failed")
		}
	} else {
		slog.Info("en translator not found")
	}
}

func transValidErr(ctx *gin.Context, err error) error {
	var userLang string

	if invalid, ok := err.(*validator.InvalidValidationError); ok {
		slog.DebugContext(ctx, "validator.InvalidValidationError", "err", invalid)
		return fmt.Errorf(i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx))
	} else if errs, ok := err.(validator.ValidationErrors); ok {
		if jwt.GetUser(ctx) != nil {
			userLang = jwt.GetUser(ctx).Language
		} else {
			langs := strings.Split(ctx.GetHeader("Accept-Language"), ";")[0]
			if len(langs) > 0 {
				userLang = strings.Split(langs, ",")[0]
			} else {
				userLang = "en"
			}
		}

		if t, exist := trans[userLang]; exist {
			// ls := make([]string, 0)
			// for _, v := range errs.Translate(*t) {
			// 	ls = append(ls, v)
			// }
			// return fmt.Errorf(strings.Join(ls, ";"))
			for _, e := range errs {
				return fmt.Errorf(e.Translate(*t))
			}
		}
		return fmt.Errorf(i18n.NewTran("common.RequestParameterIncorrect").Translate(ctx))
	} else if e, ok := err.(*i18n.Translation); ok {
		return fmt.Errorf(e.Translate(ctx))
	} else {
		return err
	}
}
