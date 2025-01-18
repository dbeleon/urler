package models

type User struct {
	Id    int64
	Name  string
	Email string
}

type Url struct {
	Id    int64
	User  int64
	Long  string
	Short string
	QR    []byte
}

type NotifyTask struct {
	Id      int64
	Short   string
	UserIds []int64
	QR      []byte
}
