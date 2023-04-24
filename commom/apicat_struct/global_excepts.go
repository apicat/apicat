package apicat_struct

type GlobalExceptsObject struct {
	Path   []int `json:"path"`
	Query  []int `json:"query"`
	Header []int `json:"header"`
	Cookie []int `json:"cookie"`
}

func (g *GlobalExceptsObject) CheckPathRef(id int) bool {
	for i, n := range g.Path {
		if n == id {
			g.Path = append(g.Path[:i], g.Path[i+1:]...)
			return true
		}
	}
	return false
}

func (g *GlobalExceptsObject) CheckQueryRef(id int) bool {
	for i, n := range g.Query {
		if n == id {
			g.Query = append(g.Query[:i], g.Query[i+1:]...)
			return true
		}
	}
	return false
}

func (g *GlobalExceptsObject) CheckHeaderRef(id int) bool {
	for i, n := range g.Header {
		if n == id {
			g.Header = append(g.Header[:i], g.Header[i+1:]...)
			return true
		}
	}
	return false
}

func (g *GlobalExceptsObject) CheckCookieRef(id int) bool {
	for i, n := range g.Cookie {
		if n == id {
			g.Cookie = append(g.Cookie[:i], g.Cookie[i+1:]...)
			return true
		}
	}
	return false
}
