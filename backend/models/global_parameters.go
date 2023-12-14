package models

import (
	"encoding/json"
	"time"

	"github.com/apicat/apicat/backend/common/apicat_struct"
	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
)

type GlobalParameters struct {
	ID        uint   `gorm:"type:bigint;primaryKey;autoIncrement"`
	ProjectID uint   `gorm:"type:bigint;index;not null;comment:项目id"`
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

func GlobalParametersImport(projectID uint, parameters *spec.HTTPParameters) map[int64]uint {
	res := virtualIDToIDMap{}

	if parameters.Header == nil && parameters.Cookie == nil && parameters.Query == nil && parameters.Path == nil {
		return res
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
					res[parameter.ID] = record.ID
				}
			}
		}
	}

	return res
}

func GlobalParametersExport(projectID uint) spec.HTTPParameters {
	globalParameters := []*GlobalParameters{}
	specHTTPParameters := spec.HTTPParameters{
		Header: make([]*spec.Schema, 0),
		Cookie: make([]*spec.Schema, 0),
		Query:  make([]*spec.Schema, 0),
		Path:   make([]*spec.Schema, 0),
	}

	if err := Conn.Where("project_id = ?", projectID).Find(&globalParameters).Error; err != nil {
		return specHTTPParameters
	}

	for _, globalParameter := range globalParameters {
		jsonschema := &jsonschema.Schema{}
		if err := json.Unmarshal([]byte(globalParameter.Schema), jsonschema); err == nil {
			required := false
			if globalParameter.Required == 1 {
				required = true
			}

			schema := spec.Schema{
				ID:       int64(globalParameter.ID),
				Name:     globalParameter.Name,
				Required: required,
				Schema:   jsonschema,
			}

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
	}

	return specHTTPParameters
}

func ReplaceGlobalParametersVirtualIDToID(content string, virtualIDToIDMap virtualIDToIDMap) string {
	docContent := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(content), &docContent); err != nil {
		return content
	}

	var (
		request    []byte
		newContent []byte
		err        error
	)

	for _, v := range docContent {
		if v["type"] == "apicat-http-request" {
			request, err = json.Marshal(v["attrs"])
			if err != nil {
				return content
			}
		}
	}

	apicatRequest := apicat_struct.RequestObject{}
	if err := json.Unmarshal(request, &apicatRequest); err != nil {
		return content
	}

	for k, v := range apicatRequest.GlobalExcepts.Cookie {
		if id, ok := virtualIDToIDMap[int64(v)]; ok {
			apicatRequest.GlobalExcepts.Cookie[k] = int(id)
		}
	}
	for k, v := range apicatRequest.GlobalExcepts.Header {
		if id, ok := virtualIDToIDMap[int64(v)]; ok {
			apicatRequest.GlobalExcepts.Header[k] = int(id)
		}
	}
	for k, v := range apicatRequest.GlobalExcepts.Path {
		if id, ok := virtualIDToIDMap[int64(v)]; ok {
			apicatRequest.GlobalExcepts.Path[k] = int(id)
		}
	}
	for k, v := range apicatRequest.GlobalExcepts.Query {
		if id, ok := virtualIDToIDMap[int64(v)]; ok {
			apicatRequest.GlobalExcepts.Query[k] = int(id)
		}
	}

	for i, v := range docContent {
		if v["type"] == "apicat-http-request" {
			docContent[i]["attrs"] = apicatRequest
		}
	}

	if newContent, err = json.Marshal(docContent); err != nil {
		return content
	}

	return string(newContent)
}
