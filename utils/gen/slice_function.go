package gen

import "github.com/cheekybits/genny/generic"

type Target generic.Type

//go:generate genny -in=$GOFILE -out=../../domain/entity/gen_posts_func.go -pkg=entity gen "Something=Post Target=int64,PostUser"
//go:generate genny -in=$GOFILE -out=../../domain/entity/gen_users_func.go -pkg=entity gen "Something=User Target=int64"
//go:generate genny -in=$GOFILE -out=../../domain/entity/gen_post_users_func.go -pkg=entity gen "Something=PostUser Target=interface{}"
func (s *Somethings) MapTarget(fc func(something *Something) Target) []Target {
	results := make([]Target, len(*s))
	for i, something := range *s {
		results[i] = fc(&something)
	}
	return results
}

func (s *Somethings) GroupByTarget(fc func(something *Something) Target) map[Target]Something {
	results := make(map[Target]Something, len(*s))
	for _, something := range *s {
		results[fc(&something)] = something
	}
	return results
}
