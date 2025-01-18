package models

type BaseResponse struct {
	Code    int    `msgpack:"code"`
	Message string `msgpack:"message"`
}

type User struct {
	Id    int64  `msgpack:"id"`
	Name  string `msgpack:"name"`
	Email string `msgpack:"email"`
}

type Url struct {
	Id    int64  `msgpack:"id"`
	User  int64  `msgpack:"user_id"`
	Long  string `msgpack:"long"`
	Short string `msgpack:"short"`
}

type UserResponse struct {
	BaseResponse
	Id int64 `msgpack:"id"`
}

type UrlResponse struct {
	BaseResponse
	Url Url `msgpack:"url"`
}
