package rpchub

import "errors"

type Response struct {
	Data []byte
	err  error
}

func (res *Response) Err() error {
	return res.err
}

func NewResponseWithErr(err string) *Response {
	r := new(Response)
	r.err = errors.New(err)
	return r
}
