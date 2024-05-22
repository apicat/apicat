package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

type CollectionNode struct {
	Node
}

type Node interface {
	NodeType() string
}

type CollectionNodes []*CollectionNode

var nodeTypes = make(map[string]reflect.Type)

func RegisterNode(n Node) {
	name := n.NodeType()
	t := reflect.TypeOf(n)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if _, ok := nodeTypes[name]; ok {
		panic(fmt.Errorf("type is not registered:%s", name))
	}
	nodeTypes[name] = t
}

func NewCollectionNodesFromJson(c string) (CollectionNodes, error) {
	if c == "" {
		return nil, errors.New("empty json content")
	}
	var collectionNodes CollectionNodes
	if err := json.Unmarshal([]byte(c), &collectionNodes); err != nil {
		return nil, err
	}
	return collectionNodes, nil
}

func (n *CollectionNode) ToHttpUrl() *CollectionHttpUrl {
	return n.Node.(*CollectionHttpUrl)
}

func (n *CollectionNode) ToHttpRequest() *CollectionHttpRequest {
	return n.Node.(*CollectionHttpRequest)
}

func (n *CollectionNode) ToHttpResponse() *CollectionHttpResponse {
	return n.Node.(*CollectionHttpResponse)
}

func (n CollectionNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Node)
}

func (n *CollectionNode) UnmarshalJSON(b []byte) error {
	var _node struct{ Type string }
	if err := json.Unmarshal(b, &_node); err != nil {
		return err
	}
	t, ok := nodeTypes[_node.Type]
	if !ok {
		return errors.New("unknown node type")
	}
	v := reflect.New(t).Interface()
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	n.Node = v.(Node)
	return nil
}

func (ns *CollectionNodes) DerefModel(ref *DefinitionModel) error {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().DerefModel(ref)
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().DerefModel(ref)
		}
	}
	return nil
}

func (ns *CollectionNodes) DeepDerefAll(params *GlobalParameters, definitions *Definitions) error {
	helper := jsonschema.NewDerefHelper(definitions.Schemas.ToJsonSchemaMap())

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DerefGlobalParameters(params)
			if err := node.ToHttpRequest().DeepDerefModelByHelper(helper); err != nil {
				return err
			}
		case NODE_HTTP_RESPONSE:
			res := node.ToHttpResponse()
			if err := res.DerefAllResponses(definitions.Responses); err != nil {
				return err
			}
			if err := res.DeepDerefModelByHelper(helper); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ns *CollectionNodes) DerefGlobalParameter(in string, param *Parameter) {
	if param == nil {
		return
	}

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Attrs.Parameters.Add(in, param)
		}
	}
}

func (ns *CollectionNodes) DerefGlobalParameters(params *GlobalParameters) {
	if params == nil {
		return
	}

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DerefGlobalParameters(params)
		}
	}
}

func (ns *CollectionNodes) DerefResponse(ref *DefinitionResponse) error {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().DerefResponse(ref)
		}
	}
	return nil
}

func (ns *CollectionNodes) DelRefModel(ref *DefinitionModel) {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelRefModel(ref)
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefModel(ref)
		}
	}
}

func (ns *CollectionNodes) DelRefResponse(ref *DefinitionResponse) {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefResponse(ref)
		}
	}
}

func (ns *CollectionNodes) DelGlobalExcept(in string, id int64) {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelGlobalExcept(in, id)
		}
	}
}

func (ns *CollectionNodes) GetUrl() *CollectionHttpUrl {
	for _, node := range *ns {
		if node.NodeType() == NODE_HTTP_URL {
			return node.ToHttpUrl()
		}
	}
	return nil
}

func (ns *CollectionNodes) GetRequest() *CollectionHttpRequest {
	for _, node := range *ns {
		if node.NodeType() == NODE_HTTP_REQUEST {
			return node.ToHttpRequest()
		}
	}
	return nil
}

func (ns *CollectionNodes) GetResponse() *CollectionHttpResponse {
	for _, node := range *ns {
		if node.NodeType() == NODE_HTTP_RESPONSE {
			return node.ToHttpResponse()
		}
	}
	return nil
}

func (ns *CollectionNodes) GetGlobalExcepts() *HttpRequestGlobalExcepts {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExcepts()
		}
	}
	return nil
}

func (ns *CollectionNodes) GetGlobalExceptToMap() map[string][]int64 {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExceptToMap()
		}
	}
	return nil
}

func (ns *CollectionNodes) GetRefModelIDs() []int64 {
	ids := make([]int64, 0)
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			ids = append(ids, node.ToHttpRequest().GetRefModelIDs()...)
		case NODE_HTTP_RESPONSE:
			ids = append(ids, node.ToHttpResponse().GetRefModelIDs()...)
		}
	}
	return ids
}

func (ns *CollectionNodes) GetRefResponseIDs() []int64 {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().GetRefResponseIDs()
		}
	}
	return nil
}

func (ns *CollectionNodes) AddReqParameter(in string, p *Parameter) {
	if p == nil {
		return
	}

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Attrs.Parameters.Add(in, p)
		}
	}
}

func (ns *CollectionNodes) SortResponses() {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().Sort()
		}
	}
}

func (ns *CollectionNodes) ToJson() (string, error) {
	res, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
