package model

import (
	"strconv"
	"strings"
)

type IdToNameMap map[uint]string
type VirtualIDToIDMap map[int64]uint

type RefContentVirtualIDToId struct {
	DefinitionSchemas    VirtualIDToIDMap
	DefinitionResponses  VirtualIDToIDMap
	DefinitionParameters VirtualIDToIDMap
	GlobalParameters     VirtualIDToIDMap
}

func ReplaceVirtualIDToID(content string, nameIDMap VirtualIDToIDMap, prefix string) string {
	for virtualID, id := range nameIDMap {
		oldStr := prefix + strconv.Itoa(int(virtualID))
		newStr := prefix + strconv.Itoa(int(id))

		content = strings.Replace(content, oldStr, newStr, -1)
	}
	return content
}
