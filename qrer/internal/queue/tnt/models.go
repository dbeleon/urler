package tnt

type BaseRequest struct{}

type BaseResponse struct {
	Code    int    `msgpack:"code"`
	Message string `msgpack:"message"`
}

type PublishRequest struct {
	Url      string `msgpack:"url"`
	Priority uint   `msgpack:"pri"`
	TTL      uint   `msgpack:"ttl"`
	Delay    uint   `msgpack:"delay"`
	TTR      uint   `msgpack:"ttr"`
}

type PublishResponse struct {
	BaseResponse
	Id int64 `msgpack:"id"`
}

type AckRequest struct {
	Id int64 `msgpack:"id"`
}

type AckResponse struct {
	BaseResponse
}

type ConsumeRequest struct {
	Timeout int `msgpack:"timeout"`
}

type ConsumeResponse struct {
	BaseResponse
	Id  int64  `msgpack:"id"`
	Url string `msgpack:"url"`
}
