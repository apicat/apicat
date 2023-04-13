package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNodeTypeNotRegistered = errors.New("type is not registered")
)

func createErrNodeTypeNotRegistered(typ string) error {
	return fmt.Errorf("%w:%s", ErrNodeTypeNotRegistered, typ)
}

// NodeProxy 为了对泛型对象进行统一的json编解码 而不是每个类型都实现json接口
// 所以使用代理套一层，统一实现json编解码
type NodeProxy struct {
	// Type string `json:"type"`
	Node
}

// Node 所有的文档节点都必须实现Node接口
type Node interface {
	NodeType() string
}

var regisotrNodeType = make(map[string]reflect.Type)

// RegisterNode 注册node节点类型
func RegisterNode(n Node) {
	name := n.NodeType()
	t := reflect.TypeOf(n)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if _, ok := regisotrNodeType[name]; ok {
		panic(createErrNodeTypeNotRegistered(name))
	}
	regisotrNodeType[name] = t
}

// CreateNodeProxy 讲一个node转为proxy对象
// node必须是已注册的
func CreateNodeProxy(n Node) (*NodeProxy, error) {
	name := n.NodeType()
	if _, ok := regisotrNodeType[name]; !ok {
		return nil, createErrNodeTypeNotRegistered(name)
	}
	return &NodeProxy{Node: n}, nil
}

func MuseCreateNodeProxy(n Node) *NodeProxy {
	p, err := CreateNodeProxy(n)
	if err != nil {
		panic(err)
	}
	return p
}

func (n *NodeProxy) UnmarshalJSON(b []byte) error {
	var _node struct{ Type string }
	if err := json.Unmarshal(b, &_node); err != nil {
		return err
	}
	t, ok := regisotrNodeType[_node.Type]
	if !ok {
		t = regisotrNodeType["doc"]
		// return createErrNodeTypeNotRegistered(_node.Type)
	}
	v := reflect.New(t).Interface()
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	n.Node = v.(Node)
	return nil
}

func (n NodeProxy) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Node)
}
