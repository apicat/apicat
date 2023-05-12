package models

type nameToIdMap map[string]uint
type IdToNameMap map[uint]string

type RefContentNameToId struct {
	DefinitionSchemas    nameToIdMap
	DefinitionResponses  nameToIdMap
	DefinitionParameters nameToIdMap
}
