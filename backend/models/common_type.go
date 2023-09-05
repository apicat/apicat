package models

type IdToNameMap map[uint]string
type virtualIDToIDMap map[int64]uint

type RefContentVirtualIDToId struct {
	DefinitionSchemas    virtualIDToIDMap
	DefinitionResponses  virtualIDToIDMap
	DefinitionParameters virtualIDToIDMap
	GolbalParameters     virtualIDToIDMap
}
