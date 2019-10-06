package rpchub

type Response struct {
	Data []byte
	err  error
}

func (res *Response) Err() error {
	return res.err
}
