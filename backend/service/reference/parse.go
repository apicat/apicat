package reference

import (
	"github.com/apicat/apicat/v2/backend/model/collection"
	"github.com/apicat/apicat/v2/backend/model/definition"
)

func ParseRefSchemasFromCollection(c *collection.Collection) ([]uint, error) {
	specC, err := c.ContentToSpec()
	if err != nil {
		return nil, err
	}

	refSchemaIDs := specC.GetRefModelIDs()

	var list []uint
	for _, v := range refSchemaIDs {
		list = append(list, uint(v))
	}

	return list, nil
}

func ParseRefSchemasFromResponse(r *definition.DefinitionResponse) ([]uint, error) {
	specR, err := r.ToSpec()
	if err != nil {
		return nil, err
	}

	refSchemaIDs := specR.RefIDs()

	var list []uint
	for _, v := range refSchemaIDs {
		list = append(list, uint(v))
	}

	return list, nil
}

func ParseRefSchemasFromSchema(s *definition.DefinitionSchema) ([]uint, error) {
	specS, err := s.ToSpec()
	if err != nil {
		return nil, err
	}

	refSchemaIDs := specS.RefIDs()

	var list []uint
	for _, v := range refSchemaIDs {
		list = append(list, uint(v))
	}

	return list, nil
}

func ParseRefResponsesFromCollection(c *collection.Collection) ([]uint, error) {
	specC, err := c.ContentToSpec()
	if err != nil {
		return nil, err
	}

	refResponseIDs := specC.GetRefResponseIDs()

	var list []uint
	for _, v := range refResponseIDs {
		list = append(list, uint(v))
	}

	return list, nil
}
