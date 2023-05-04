package util

import (
	"regexp"
	"strconv"
	"strings"
)

func ReplaceIDToName(content string, idToNameMap map[uint]string, prefix string) string {
	re := regexp.MustCompile(prefix + `\d+`)
	reID := regexp.MustCompile(`\d+`)

	for {
		match := re.FindString(content)
		if match == "" {
			break
		}

		schemasIDStr := reID.FindString(match)
		if schemasIDInt, err := strconv.Atoi(schemasIDStr); err == nil {
			schemasID := uint(schemasIDInt)
			content = strings.Replace(content, match, prefix+idToNameMap[schemasID], -1)
		}
	}

	return content
}
