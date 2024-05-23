package spec

import (
	"errors"
	"strconv"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

const NODE_HTTP_RESPONSE = "apicat-http-response"

type CollectionHttpResponse struct {
	Type  string             `json:"type" yaml:"type"`
	Attrs *HttpResponseAttrs `json:"attrs" yaml:"attrs"`
}

type HttpResponseAttrs struct {
	List Responses `json:"list" yaml:"list"`
}

func init() {
	RegisterNode(&CollectionHttpResponse{
		Type: NODE_HTTP_RESPONSE,
	})
}

func NewCollectionHttpResponse() *CollectionHttpResponse {
	return &CollectionHttpResponse{
		Type: NODE_HTTP_RESPONSE,
		Attrs: &HttpResponseAttrs{
			List: make(Responses, 0),
		},
	}
}

func (r *CollectionHttpResponse) NodeType() string {
	return r.Type
}

func (r *CollectionHttpResponse) DerefResponse(ref *DefinitionResponse) error {
	if ref == nil {
		return errors.New("response is nil")
	}

	for _, v := range r.Attrs.List {
		if v.Ref() {
			if err := v.ReplaceRef(&ref.BasicResponse); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *CollectionHttpResponse) DerefAllResponses(refs DefinitionResponses) error {
	if len(refs) == 0 {
		return nil
	}

	refsMap := refs.ToMap()

	for _, res := range r.Attrs.List {
		if res.Ref() {
			if ref, ok := refsMap[res.GetRefID()]; ok {
				if err := res.ReplaceRef(&ref.BasicResponse); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *CollectionHttpResponse) DerefModel(ref *DefinitionModel) error {
	if ref == nil {
		return errors.New("model is nil")
	}

	for _, res := range r.Attrs.List {
		if res.Content == nil {
			continue
		}

		for _, v := range res.Content {
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
	}
	return nil
}

func (r *CollectionHttpResponse) DeepDerefModel(refs DefinitionModels) error {
	if len(refs) == 0 {
		return nil
	}

	helper := jsonschema.NewDerefHelper(refs.ToJsonSchemaMap())
	return r.DeepDerefModelByHelper(helper)
}

func (r *CollectionHttpResponse) DeepDerefModelByHelper(helper *jsonschema.DerefHelper) error {
	if helper == nil {
		return errors.New("helper is nil")
	}

	for _, res := range r.Attrs.List {
		if res.Content == nil {
			continue
		}

		for _, v := range res.Content {
			if v.Schema != nil {
				new, err := helper.DeepDeref(v.Schema)
				if err != nil {
					return err
				}
				v.Schema = &new
			}
		}
	}
	return nil
}

func (r *CollectionHttpResponse) DelRefResponse(ref *DefinitionResponse) {
	if ref == nil {
		return
	}

	for i, v := range r.Attrs.List {
		if v.Ref() && v.GetRefID() == ref.ID {
			r.Attrs.List = append(r.Attrs.List[:i], r.Attrs.List[i+1:]...)
			return
		}
	}
}

func (r *CollectionHttpResponse) DelRefModel(ref *DefinitionModel) {
	if ref == nil {
		return
	}

	for _, res := range r.Attrs.List {
		if res.Content == nil {
			continue
		}

		for _, v := range res.Content {
			if v.Schema != nil {
				v.Schema.DelRef(ref.Schema)
			}
		}
	}
}

func (r *CollectionHttpResponse) GetRefModelIDs() []int64 {
	ids := make([]int64, 0)
	for _, res := range r.Attrs.List {
		if res.Content == nil {
			continue
		}

		for _, v := range res.Content {
			if v.Schema != nil {
				ids = append(ids, v.Schema.DeepGetRefID()...)
			}
		}
	}
	return ids
}

func (r *CollectionHttpResponse) GetRefResponseIDs() []int64 {
	ids := make([]int64, 0)
	for _, res := range r.Attrs.List {
		if res.Ref() {
			ids = append(ids, res.GetRefID())
		}
	}
	return ids
}

func (r *CollectionHttpResponse) ToCollectionNode() *CollectionNode {
	return &CollectionNode{
		Node: r,
	}
}

func (r *CollectionHttpResponse) Sort() {
	r.Attrs.List.Sort()
}
