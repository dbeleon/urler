package models

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
