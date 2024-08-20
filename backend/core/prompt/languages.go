package prompt

var langMap = map[string]string{
	"zh":    "Chinese",
	"zh-CN": "Chinese",
	"en":    "English",
	"en-US": "English",
}

func language(lang string) string {
	if l, ok := langMap[lang]; ok {
		return l
	}
	return "English"
}
