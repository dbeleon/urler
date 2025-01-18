package qr

import (
	"bytes"
	"errors"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/compressed"
)

var (
	ErrGenerate = errors.New("could not generate QRCode")
	ErrSave     = errors.New("could not save image")
)

type Code struct {
	opt compressed.Option
}

func New() *Code {
	return &Code{
		opt: compressed.Option{
			Padding:   4, // padding pixels around the qr code.
			BlockSize: 1, // block pixels which represents a bit data.
		},
	}
}

func (c *Code) Encode(text string) ([]byte, error) {
	qrc, err := qrcode.New(text)
	if err != nil {
		return nil, ErrGenerate
	}

	data := make([]byte, 0, 32*1024)
	buf := bytes.NewBuffer(data)
	wc := &WrCl{Buf: buf}
	w := compressed.NewWithWriter(wc, &c.opt)
	if err = qrc.Save(w); err != nil {
		return nil, ErrSave
	}

	return wc.Buf.Bytes(), nil
}

type WrCl struct {
	Buf *bytes.Buffer
}

func (w *WrCl) Write(data []byte) (int, error) {
	return w.Buf.Write(data)
}

func (w *WrCl) Close() error {
	return nil
}
