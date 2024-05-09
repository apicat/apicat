package spec2

import "strconv"

const NODE_HTTP_RESPONSE = "apicat-http-response"

type CollectionHttpResponse struct {
	Type  string             `json:"type" yaml:"type"`
	Attrs *HttpResponseAttrs `json:"attr" yaml:"attrs"`
}

type HttpResponseAttrs struct {
	List Responses `json:"list" yaml:"list"`
}

func NewCollectionHttpResponse() *CollectionHttpResponse {
	return &CollectionHttpResponse{
		Type: NODE_HTTP_RESPONSE,
		Attrs: &HttpResponseAttrs{
			List: Responses{},
		},
	}
}

func (r *CollectionHttpResponse) NodeType() string {
	return r.Type
}

func (r *CollectionHttpResponse) DerefResponse(ref *DefinitionResponse) {
	if r == nil || r.Attrs == nil || ref == nil {
		return
	}

	for _, v := range r.Attrs.List {
		if v.Ref() {
			v.ReplaceRef(&ref.Response)
		}
	}
}

func (r *CollectionHttpResponse) DerefModel(ref *DefinitionModel) {
	if r == nil || r.Attrs == nil || ref == nil {
		return
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
						schema.ReplaceRef(ref.Schema)
					}
				}
			}
		}
	}
}

func (r *CollectionHttpResponse) DelRefResponse(ref *DefinitionResponse) {
	if r == nil || r.Attrs == nil || ref == nil {
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
	if r == nil || r.Attrs == nil || ref == nil {
		return
	}

	for _, res := range r.Attrs.List {
		if res.Content == nil {
			continue
		}

		for _, v := range res.Content {
			if v.Schema != nil {
				if v.Schema.Ref() {
					v.Schema.DelRootRef(ref.Schema)
				}
				v.Schema.DelChildrenRef(ref.Schema)
			}
		}
	}
}
