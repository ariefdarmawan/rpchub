package rpchub

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/eaciit/toolkit"
)

type Hub struct {
	fns map[string]reflect.Value
	log *toolkit.LogEngine
}

func NewHub() *Hub {
	r := new(Hub)
	return r
}

func SetHubLog(r *Hub, l *toolkit.LogEngine) *Hub {
	r.log = l
	return r
}

func HubLog(p *Hub) *toolkit.LogEngine {
	if p.log == nil {
		p.log = toolkit.NewLogEngine(true, false, "", "", "")
	}
	return p.log
}

func (r *Hub) Call(request Request, response *Response) error {
	var err error

	name := strings.ToLower(request.Name)
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				err = fmt.Errorf("%s panic error. %v", name, rec)
			}
		}()

		if r.fns == nil {
			err = fmt.Errorf("%s error. invalid initialization", name)
			return
		}

		fn, ok := r.fns[name]
		if !ok {
			err = fmt.Errorf("service %s is not exist", name)
			return
		}

		var ins []reflect.Value
		var outs []reflect.Value
		if len(request.Parm) > 0 {
			ins = make([]reflect.Value, len(request.Parm))
			for idx, parm := range request.Parm {
				ins[idx] = reflect.ValueOf(parm)
			}
		}
		//fmt.Println("call", name, "number of parms", len(request.Parm), "raw:", ins, "json:", toolkit.JsonString(ins))
		outs = fn.Call(ins)
		o := outs[0]
		t := o.Type()
		k := t.Kind()
		//fmt.Println("return type:", t.String())
		if t.String() == "[]uint8" {
			response.Data = o.Interface().([]byte)
		} else if k == reflect.String {
			response.Data = []byte(o.String())
		} else if k == reflect.Int || k == reflect.Int8 ||
			k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
			response.Data = IntToBytes(o.Int())
		} else if k == reflect.Float32 || k == reflect.Float64 {
			response.Data = FloatToBytes(o.Float())
		} else {
			response.Data = toolkit.ToBytes(outs[0].Interface(), "")
		}
	}()

	return err
}

func RegisterToHub(hub *Hub, objs ...interface{}) *Hub {
	if hub == nil {
		hub = NewHub()
	}

	for _, obj := range objs {
		registerObjToHub(hub, obj)
	}

	return hub
}

func registerObjToHub(hub *Hub, obj interface{}) error {
	rv := reflect.ValueOf(obj)
	rt := rv.Type()

	fnCount := rv.NumMethod()
	HubLog(hub).Infof("registering object to Hub: %s. %d method(s) found",
		rt.Elem().Name(), fnCount)
	for i := 0; i < fnCount; i++ {
		fn := rv.Method(i)
		ft := rt.Method(i)
		fnName := strings.ToLower(fmt.Sprintf("%s.%s",
			rt.Elem().Name(), ft.Name))

		if hub.fns == nil {
			hub.fns = map[string]reflect.Value{}
		}

		hub.fns[fnName] = fn
		HubLog(hub).Infof("adding to RPC Hub: %s", fnName)
	}

	return nil
}

func IntToBytes(v int64) []byte {
	i := uint64(v)
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, i)
	return bs
}

func BytesToInt(bs []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bs))
}

func FloatToBytes(v float64) []byte {
	bs := make([]byte, 8)
	fb := math.Float64bits(v)
	binary.LittleEndian.PutUint64(bs[:], fb)
	return bs
}

func BytesToFloat(bs []byte) float64 {
	bits := binary.LittleEndian.Uint64(bs)
	return math.Float64frombits(bits)
}
