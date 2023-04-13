package translator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func UseValidatori18n() gin.HandlerFunc {
	uni := ut.New(en.New(), zh.New())
	validater := binding.Validator.Engine().(*validator.Validate)
	zhtran, _ := uni.GetTranslator("zh")
	entran, _ := uni.GetTranslator("en")
	zh_translations.RegisterDefaultTranslations(validater, zhtran)
	en_translations.RegisterDefaultTranslations(validater, entran)
	return func(ctx *gin.Context) {
		ctx.Set("ValidatorTrans", uni)
		// 处理非ShouldBind直接终止的情况
		ctx.Next()
		if err := ctx.Errors.ByType(gin.ErrorTypeBind); err != nil {
			ctx.Writer.Write([]byte(ValiadteTransErr(ctx, err.Last().Err).Error()))
		}
	}
}

func ValiadteTransErr(ctx *gin.Context, err error) error {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		if tran, ok := ctx.Get("ValidatorTrans"); ok {
			langs := strings.Split(ctx.GetHeader("Accept-Language"), ";")[0]
			t, _ := tran.(*ut.UniversalTranslator).FindTranslator(strings.Split(langs, ",")...)
			ls := make([]string, 0)
			for _, v := range errs.Translate(t) {
				ls = append(ls, v)
			}
			return fmt.Errorf(strings.Join(ls, ";"))
		}
	}
	return err
}
