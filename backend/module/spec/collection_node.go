package spec

type CollectionNode struct {
	Node
}

type Node interface {
	NodeType() string
}

type CollectionNodes []*CollectionNode

func (n *CollectionNode) ToHttpUrl() *CollectionHttpUrl {
	return n.Node.(*CollectionHttpUrl)
}

func (n *CollectionNode) ToHttpRequest() *CollectionHttpRequest {
	return n.Node.(*CollectionHttpRequest)
}

func (n *CollectionNode) ToHttpResponse() *CollectionHttpResponse {
	return n.Node.(*CollectionHttpResponse)
}
