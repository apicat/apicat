package spec

import "encoding/json"

type CollectionNode struct {
	Node
}

type Node interface {
	NodeType() string
}

type CollectionNodes []*CollectionNode

func NewCollectionNodesFromJson(c string) (CollectionNodes, error) {
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

func (ns *CollectionNodes) GetUrlInfo() (method string, path string) {
	for _, node := range *ns {
		if node.NodeType() == NODE_HTTP_URL {
			url := node.ToHttpUrl()
			return url.Attrs.Method, url.Attrs.Path
		}
	}
	return method, path
}
