package jsonschema

import "strconv"

func ReplaceAll(replace *Schema, subjects []*Schema) {
	for _, s := range subjects {
		if s == nil {
			continue
		}

		refs := s.DeepFindRefById(strconv.FormatInt(replace.ID, 10))
		if len(refs) > 0 {
			for _, ref := range refs {
				ref.ReplaceRef(replace)
			}
		}
	}
}
