package gen

//go:generate genny -in=$GOFILE -out=../../user/slice.gen.go -pkg=user gen "Source=User"
//go:generate genny -in=$GOFILE -out=../../post/slice.gen.go -pkg=post gen "Source=Post"
//go:generate genny -in=$GOFILE -out=../../post/with_user_slice.gen.go -pkg=post gen "Source=WithUser"

type Sources []Source
