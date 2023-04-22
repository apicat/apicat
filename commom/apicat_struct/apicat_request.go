package apicat_struct

type ApiCatHttpRequestNodeObject struct {
	Type  string         `json:"type"`
	Attrs *RequestObject `json:"attrs"`
}
