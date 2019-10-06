package rpchub_test

import (
	"fmt"

	"github.com/eaciit/toolkit"
)

type Obj struct {
	pre string
}

func newObj(pre string) *Obj {
	o := new(Obj)
	o.pre = pre
	return o
}

func (obj *Obj) Ping() string {
	return "RPC Object OK"
}

func (obj *Obj) Int() int {
	return toolkit.RandInt(1000)
}

func (obj *Obj) Float(f float64) float64 {
	return f * float64(2)
}

func (obj *Obj) SetObjData(data string) string {
	return fmt.Sprintf("%s %s", obj.pre, data)
}
