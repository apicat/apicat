package except

import "github.com/apicat/apicat/v2/backend/model/collection"

func ParseExceptParamsFromCollection(c *collection.Collection) ([]uint, error) {
	specC, err := c.ContentToSpec()
	if err != nil {
		return nil, err
	}

	ExceptParamIDs := specC.GetGlobalExceptAll()

	var list []uint
	for _, gp := range ExceptParamIDs {
		for _, v := range gp {
			list = append(list, uint(v))
		}
	}

	return list, nil
}
