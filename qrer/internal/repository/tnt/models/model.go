package models

type Url struct {
	Id    int64  `msgpack:"id"`
	User  int64  `msgpack:"user_id"`
	Long  string `msgpack:"long"`
	Short string `msgpack:"short"`
	QR    string `msgpack:"qr"`
}

type BaseResponse struct {
	Code    int    `msgpack:"code"`
	Message string `msgpack:"message"`
}

type UserIDsResponse struct {
	BaseResponse
	UserIDs []int64 `msgpack:"user_ids"`
}
