package rpchub_test

import (
	"testing"
	"time"

	"github.com/eaciit/toolkit"

	"github.com/ariefdarmawan/rpchub"
	"github.com/smartystreets/goconvey/convey"
)

func TestHub(t *testing.T) {
	convey.Convey("prepare Hub", t, func() {
		hub := rpchub.NewHub()
		obj := newObj("Hub")
		rpchub.RegisterToHub(hub, obj)
		res := new(rpchub.Response)
		err := hub.Call(rpchub.Request{Name: "obj.Ping"}, res)
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("validate", func() {
			convey.So(res.Err(), convey.ShouldBeNil)
			data := string(res.Data)
			convey.So(data, convey.ShouldEqual, "RPC Object OK")
		})

		convey.Convey("check int", func() {
			res := new(rpchub.Response)
			_ = hub.Call(rpchub.Request{Name: "obj.Int"}, res)
			convey.So(res.Err(), convey.ShouldBeNil)

			x := rpchub.BytesToInt(res.Data)
			convey.So(x, convey.ShouldBeBetween, 1, 1000)
			convey.Println("\nX value: ", x)
		})

		convey.Convey("check float", func() {
			res := new(rpchub.Response)
			err := hub.Call(rpchub.Request{Name: "obj.Float", Parm: []interface{}{800.50}}, res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)
			convey.Printf("\nData: %v\n", res.Data)

			float := rpchub.BytesToFloat(res.Data)
			convey.So(float, convey.ShouldEqual, float64(2)*800.50)
		})

		convey.Convey("check internal variable", func() {
			res := new(rpchub.Response)
			err := hub.Call(rpchub.Request{Name: "obj.SetObjData", Parm: []interface{}{"Arief"}}, res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)

			convey.So(string(res.Data), convey.ShouldEqual, "Hub Arief")
		})

		convey.Convey("check struct", func() {
			res := new(rpchub.Response)
			t0 := time.Now()
			err := hub.Call(rpchub.NewRequest("obj.Struct", t0), res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)

			dummy := new(Obj1)
			dummy.Title = "Obj 1"
			dummy.Time = t0
			dummyJsonTxt := toolkit.JsonString(dummy)
			convey.Println("struct data:", res.Data)

			resObj := new(Obj1)
			e := toolkit.FromBytes(res.Data, "", resObj)
			convey.So(e, convey.ShouldBeNil)

			resObjJsonTxt := toolkit.JsonString(resObj)
			convey.So(resObjJsonTxt, convey.ShouldEqual, dummyJsonTxt)
		})

		convey.Convey("check bytes", func() {
			err := hub.Call(rpchub.NewRequest("obj.Bytes", "Arief Darmawan"), res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)

			convey.So(string(res.Data), convey.ShouldEqual, "Arief Darmawan")
		})

		convey.Convey("check time", func() {
			t0 := time.Now()
			err := hub.Call(rpchub.NewRequest("obj.Time", t0), res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)

			var t1 time.Time
			err = toolkit.FromBytes(res.Data, "", &t1)
			convey.So(err, convey.ShouldBeNil)
			convey.So(t1.Sub(t0), convey.ShouldEqual, time.Duration(10)*time.Minute)
		})

		convey.Convey("check slice", func() {
			err := hub.Call(rpchub.NewRequest("obj.Slice", "hello"), res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)

			var hellos []string
			err = toolkit.FromBytes(res.Data, "", &hellos)
			convey.So(err, convey.ShouldBeNil)
			convey.So(len(hellos), convey.ShouldEqual, 10)
		})

	})
}
