package spec

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

const TYPE_RESPONSE = "response"

type DefinitionResponse struct {
	BasicResponse
	ParentId int64               `json:"parentid,omitempty" yaml:"parentid,omitempty"`
	Type     string              `json:"type,omitempty" yaml:"type,omitempty"`
	Items    DefinitionResponses `json:"items,omitempty" yaml:"items,omitempty"`
}

type DefinitionResponses []*DefinitionResponse

func NewDefinitionResponseFromJson(str string) (*DefinitionResponse, error) {
	if str == "" {
		return nil, errors.New("empty json content")
	}
	s := &DefinitionResponse{}
	if err := json.Unmarshal([]byte(str), s); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *DefinitionResponse) RefIDs() (ids []int64) {
	for _, v := range r.Content {
		if v.Schema != nil {
			ids = append(ids, v.Schema.DeepGetRefID()...)
		}
	}
	return
}

func (r *DefinitionResponse) Deref(ref *DefinitionModel) error {
	if ref == nil {
		return errors.New("model is nil")
	}
	ref.Schema.ID = ref.ID

	for _, v := range r.Content {
		if v.Schema != nil {
			refSchemas := v.Schema.DeepFindRefById(strconv.FormatInt(ref.ID, 10))
			if len(refSchemas) > 0 {
				for _, schema := range refSchemas {
					if err := schema.ReplaceRef(ref.Schema); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *DefinitionResponse) DeepDeref(refs DefinitionModels) error {
	if len(refs) == 0 {
		return nil
	}

	helper := jsonschema.NewDerefHelper(refs.ToJsonSchemaMap())

	for _, v := range r.Content {
		if v.Schema != nil {
			s, err := helper.DeepDeref(v.Schema)
			if err != nil {
				return err
			}
			v.Schema = &s
		}
	}
	return nil
}

func (r *DefinitionResponse) DelRef(ref *DefinitionModel) {
	if ref == nil {
		return
	}
	ref.Schema.ID = ref.ID

	for _, v := range r.Content {
		if v.Schema != nil {
			v.Schema.DelRef(ref.Schema)
		}
	}
}

func (r *DefinitionResponse) ItemsTreeToList() DefinitionResponses {
	list := make(DefinitionResponses, 0)
	if r.Type == TYPE_MODEL {
		list = append(list, r)
	}
	if r.Items != nil && len(r.Items) > 0 {
		for _, v := range r.Items {
			list = append(list, v.ItemsTreeToList()...)
		}
	}
	return list
}

func (r *DefinitionResponses) FindByName(name string) *DefinitionResponse {
	for _, v := range *r {
		if v.Type == TYPE_CATEGORY {
			return v.Items.FindByName(name)
		}
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (r *DefinitionResponses) FindByID(id int64) *DefinitionResponse {
	for _, v := range *r {
		if v.Type == TYPE_CATEGORY {
			return v.Items.FindByID(id)
		}
		if id == v.ID {
			return v
		}
	}
	return nil
}

func (r *DefinitionResponses) DelByID(id int64) {
	for i, v := range *r {
		if v.Type == TYPE_CATEGORY {
			v.Items.DelByID(id)
		} else {
			if id == v.ID {
				*r = append((*r)[:i], (*r)[i+1:]...)
			}
		}
	}
}

func (r *DefinitionResponses) ToMap() map[int64]*DefinitionResponse {
	m := make(map[int64]*DefinitionResponse)
	for _, v := range *r {
		if v.Type == TYPE_RESPONSE {
			m[v.ID] = v
		}
	}
	return m
}
