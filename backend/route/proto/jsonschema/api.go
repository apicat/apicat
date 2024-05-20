package jsonschema

import (
	"github.com/apicat/apicat/v2/backend/route/proto/jsonschema/base"
	"github.com/gin-gonic/gin"
)

type JsonSchemaApi interface {
	Parse(*gin.Context, *base.JsonSchemaOption) (*base.JsonSchemaOption, error)
}
