// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package post

func (s *WithUsers) MapInterface(fc func(something *WithUser) interface{}) []interface{} {
	results := make([]interface{}, len(*s))
	for i, something := range *s {
		results[i] = fc(&something)
	}
	return results
}

func (s *WithUsers) GroupByInterface(fc func(something *WithUser) interface{}) map[interface{}]WithUser {
	results := make(map[interface{}]WithUser, len(*s))
	for _, something := range *s {
		results[fc(&something)] = something
	}
	return results
}