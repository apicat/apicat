package jsonschema

import (
	"net/http"

	"github.com/apicat/apicat/v2/backend/i18n"
	spec "github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
	"github.com/apicat/apicat/v2/backend/route/proto/jsonschema"
	"github.com/apicat/apicat/v2/backend/route/proto/jsonschema/base"
	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type jsonSchemaApiImpl struct{}

func NewJsonSchemaApi() jsonschema.JsonSchemaApi {
	return &jsonSchemaApiImpl{}
}

func (j *jsonSchemaApiImpl) Parse(ctx *gin.Context, opt *base.JsonSchemaOption) (*base.JsonSchemaOption, error) {
	js, err := spec.NewSchemaFromJson(opt.JsonSchema)
	if err != nil {
		return nil, ginrpc.NewError(http.StatusBadRequest, i18n.NewErr("jsonschema.JsonSchemaIncorrect"))
	}
	return &base.JsonSchemaOption{
		JsonSchema: js.ToJson(),
	}, nil
}
