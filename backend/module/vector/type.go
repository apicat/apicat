package vector

type DataType interface {
	TypeString() string
}

type Property struct {
	Name         string
	DataType     DataType
	Description  string
	Tokenization string
}

type Properties []*Property

type ObjectData struct {
	ID         string
	Properties map[string]interface{}
	Vector     []float32
}

type ObjectDataList []*ObjectData

type SearchOption struct {
	Vector           []float32
	Fields           []string
	AdditionalFields []string
	Offset           int
	Limit            int
	Distance         float32
	Certainty        float32
	WhereCondition   []*WhereCondition
}

type WhereCondition struct {
	PropertyName string
	Operator     string
	Value        DataType
}

type T_TEXT string
type T_TEXT_LIST []string
type T_INT int
type T_INT_LIST []int
type T_BOOL bool
type T_BOOL_LIST []bool
type T_NUM float64
type T_NUM_LIST []float64

func (T_TEXT) TypeString() string {
	return "text"
}
func (T_TEXT_LIST) TypeString() string {
	return "text[]"
}
func (T_INT) TypeString() string {
	return "int"
}
func (T_INT_LIST) TypeString() string {
	return "int[]"
}
func (T_BOOL) TypeString() string {
	return "boolean"
}
func (T_BOOL_LIST) TypeString() string {
	return "boolean[]"
}
func (T_NUM) TypeString() string {
	return "number"
}
func (T_NUM_LIST) TypeString() string {
	return "number[]"
}
