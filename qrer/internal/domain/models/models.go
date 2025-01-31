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

type QRTask struct {
	Id    int64
	Host  string
	Short string
	QR    []byte
}

type NotifTask struct {
	Short   string
	UserIDs []int64
	QR      []byte
}
