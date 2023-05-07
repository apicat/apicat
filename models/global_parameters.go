package models

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/apicat/apicat/common/spec"
)

type GlobalParameters struct {
	ID        uint   `gorm:"type:integer primary key autoincrement"`
	ProjectID uint   `gorm:"type:integer;index;not null;comment:项目id"`
	In        string `gorm:"type:varchar(255);not null;comment:位置:header,cookie,query,path"`
	Name      string `gorm:"type:varchar(255);not null;comment:参数名称"`
	Required  int    `gorm:"type:tinyint(1);not null;comment:是否必传:0-否,1-是"`
	Schema    string `gorm:"type:mediumtext;comment:参数内容"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGlobalParameters(ids ...uint) (*GlobalParameters, error) {
	globalParameters := &GlobalParameters{}
	if len(ids) > 0 {
		if err := Conn.Take(globalParameters, ids[0]).Error; err != nil {
			return globalParameters, err
		}
		return globalParameters, nil
	}
	return globalParameters, nil
}

func (gp *GlobalParameters) List() ([]*GlobalParameters, error) {
	globalParametersQuery := Conn.Where("project_id = ?", gp.ProjectID)

	var globalParameters []*GlobalParameters
	return globalParameters, globalParametersQuery.Find(&globalParameters).Error
}

func (gp *GlobalParameters) GetCountByName() (int64, error) {
	var count int64
	return count, Conn.Model(&GlobalParameters{}).Where("project_id = ? and name = ? and \"in\" = ?", gp.ProjectID, gp.Name, gp.In).Count(&count).Error
}

func (gp *GlobalParameters) GetCountExcludeTheID() (int64, error) {
	var count int64
	return count, Conn.Model(&GlobalParameters{}).Where("project_id = ? and name = ? and \"in\" = ? and id != ?", gp.ProjectID, gp.Name, gp.In, gp.ID).Count(&count).Error
}

func (gp *GlobalParameters) Create() error {
	return Conn.Create(gp).Error
}

func (gp *GlobalParameters) Update() error {
	return Conn.Save(gp).Error
}

func (gp *GlobalParameters) Delete() error {
	return Conn.Delete(gp).Error
}

func GlobalParametersImport(projectID uint, parameters *spec.HTTPParameters) map[string]nameToIdMap {
	var parametersMap = map[string]nameToIdMap{
		"header": make(nameToIdMap),
		"cookie": make(nameToIdMap),
		"query":  make(nameToIdMap),
		"path":   make(nameToIdMap),
	}

	if parameters.Header == nil && parameters.Cookie == nil && parameters.Query == nil && parameters.Path == nil {
		return parametersMap
	}

	var params []*spec.Schema
	parameterList := []string{"header", "cookie", "query", "path"}
	for _, key := range parameterList {
		switch key {
		case "header":
			params = parameters.Header
		case "cookie":
			params = parameters.Cookie
		case "query":
			params = parameters.Query
		case "path":
			params = parameters.Path
		}

		for _, parameter := range params {
			if parameterStr, err := json.Marshal(parameter.Schema); err == nil {
				required := 0
				if parameter.Required {
					required = 1
				}

				record := &GlobalParameters{
					ProjectID: projectID,
					In:        key,
					Name:      parameter.Name,
					Required:  required,
					Schema:    string(parameterStr),
				}

				if Conn.Create(record).Error == nil {
					parametersMap[key][parameter.Name] = record.ID
				}
			}
		}
	}

	return parametersMap
}

func GlobalParametersExport(projectID uint) spec.HTTPParameters {
	var globalParameters []*GlobalParameters
	var specHTTPParameters spec.HTTPParameters

	if err := Conn.Where("project_id = ?", projectID).Find(&globalParameters).Error; err != nil {
		return specHTTPParameters
	}

	for _, globalParameter := range globalParameters {
		var schema spec.Schema
		schema.Name = globalParameter.Name
		required := false
		if globalParameter.Required == 1 {
			required = true
		}
		schema.Required = required
		json.Unmarshal([]byte(globalParameter.Schema), &schema.Schema)

		switch globalParameter.In {
		case "header":
			specHTTPParameters.Header = append(specHTTPParameters.Header, &schema)
		case "cookie":
			specHTTPParameters.Cookie = append(specHTTPParameters.Cookie, &schema)
		case "query":
			specHTTPParameters.Query = append(specHTTPParameters.Query, &schema)
		case "path":
			specHTTPParameters.Path = append(specHTTPParameters.Path, &schema)
		}
	}

	return specHTTPParameters
}

type GlobalParametersIDToNameMap struct {
	Header IdToNameMap
	Cookie IdToNameMap
	Query  IdToNameMap
	Path   IdToNameMap
}

func (gpMap *GlobalParametersIDToNameMap) GlobalParametersIDToNameMapInit(projectID uint) {
	var globalParameters []*GlobalParameters
	if err := Conn.Where("project_id = ?", projectID).Find(&globalParameters).Error; err == nil {
		for _, globalParameter := range globalParameters {
			switch globalParameter.In {
			case "header":
				gpMap.Header[globalParameter.ID] = globalParameter.Name
			case "cookie":
				gpMap.Cookie[globalParameter.ID] = globalParameter.Name
			case "query":
				gpMap.Query[globalParameter.ID] = globalParameter.Name
			case "path":
				gpMap.Path[globalParameter.ID] = globalParameter.Name
			}
		}
	}
}

type GlobalExceptsID struct {
	Header []uint `json:"header"`
	Path   []uint `json:"path"`
	Cookie []uint `json:"cookie"`
	Query  []uint `json:"query"`
}

type GlobalExceptsName struct {
	Header []string `json:"header"`
	Path   []string `json:"path"`
	Cookie []string `json:"cookie"`
	Query  []string `json:"query"`
}

func GlobalParametersExceptsIDToName(content string, gpMap GlobalParametersIDToNameMap) string {
	re := regexp.MustCompile(`"globalExcepts":{[^}]*}`)
	match := re.FindString(content)

	if match == "" {
		return content
	}

	oldGlobalExcepts := match[16:]

	var globalExceptsID GlobalExceptsID
	globalExceptsName := GlobalExceptsName{
		Header: []string{},
		Path:   []string{},
		Cookie: []string{},
		Query:  []string{},
	}

	if err := json.Unmarshal([]byte(oldGlobalExcepts), &globalExceptsID); err != nil {
		return content
	}

	for _, id := range globalExceptsID.Header {
		globalExceptsName.Header = append(globalExceptsName.Header, gpMap.Header[id])
	}
	for _, id := range globalExceptsID.Path {
		globalExceptsName.Path = append(globalExceptsName.Path, gpMap.Path[id])
	}
	for _, id := range globalExceptsID.Cookie {
		globalExceptsName.Cookie = append(globalExceptsName.Cookie, gpMap.Cookie[id])
	}
	for _, id := range globalExceptsID.Query {
		globalExceptsName.Query = append(globalExceptsName.Query, gpMap.Query[id])
	}

	if newGlobalExcepts, err := json.Marshal(globalExceptsName); err == nil {
		content = strings.Replace(content, oldGlobalExcepts, string(newGlobalExcepts), -1)
	}

	return content
}
