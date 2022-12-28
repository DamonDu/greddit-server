package service

var (
	User UserService
	Post PostService
)

func Init() {
	User = NewUserService()
	Post = NewPostService()
}
