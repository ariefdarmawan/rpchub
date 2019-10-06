package hubserver_test

import (
	"fmt"
	"time"

	"github.com/eaciit/toolkit"
)

type Obj struct {
	pre string
}

type Obj1 struct {
	Title string
	Time  time.Time
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

func (obj *Obj) Struct(d time.Time) *Obj1 {
	o := new(Obj1)
	o.Title = "Obj 1"
	o.Time = d
	return o
}

func (obj *Obj) Time(d time.Time) time.Time {
	return d.Add(10 * time.Minute)
}

func (obj *Obj) Bytes(v string) []byte {
	return []byte(v)
}

func (obj *Obj) Slice(v string) []string {
	x := make([]string, 10)
	for i := 0; i < 10; i++ {
		x[i] = v
	}
	return x
}
