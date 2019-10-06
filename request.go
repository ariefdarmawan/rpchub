package rpchub

type Request struct {
	Name string
	Parm []interface{}
}

func NewRequest(name string, parms ...interface{}) Request {
	r := Request{}
	r.Name = name
	r.Parm = parms
	return r
}
