package util

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserLanguages(ctx *gin.Context) []string {
	langs := strings.Split(ctx.GetHeader("Accept-Language"), ";")
	var result []string

	if len(langs) > 2 {
		result = strings.Split(langs[0], ",")

		for i := 1; i < len(langs); i++ {
			tmp := strings.Split(langs[i], ",")
			result = append(result, tmp[len(tmp)-1])
		}
	} else if len(langs) > 1 {
		result = strings.Split(langs[0], ",")
	}

	return result
}

func GetUserFullLanguage(ctx *gin.Context) string {
	lang := "zh-CN"
	langs := GetUserLanguages(ctx)
	if len(langs) > 0 {
		return langs[0]
	}
	return lang
}

func GetUserLanguage(ctx *gin.Context) string {
	lang := GetUserFullLanguage(ctx)
	return strings.Split(lang, "-")[0]
}
