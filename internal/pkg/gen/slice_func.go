package gen

//go:generate genny -in=$GOFILE -out=../../user/slice_func.gen.go -pkg=user gen "Source=User Target=int64"
//go:generate genny -in=$GOFILE -out=../../post/slice_func.gen.go -pkg=post gen "Source=Post Target=int64,WithUser"
//go:generate genny -in=$GOFILE -out=../../post/with_user_slice_func.gen.go -pkg=post gen "Source=WithUser Target=interface{}"
func (s *Sources) MapTarget(fc func(something *Source) Target) []Target {
	results := make([]Target, len(*s))
	for i, something := range *s {
		results[i] = fc(&something)
	}
	return results
}

func (s *Sources) GroupByTarget(fc func(something *Source) Target) map[Target]Source {
	results := make(map[Target]Source, len(*s))
	for _, something := range *s {
		results[fc(&something)] = something
	}
	return results
}
