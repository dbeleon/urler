package models

type User struct {
	Id    int64
	Name  string
	Email string
}

type Url struct {
	User  int64
	Long  string
	Short string
	Qr    []byte
}
