package rpchub_test

import (
	"testing"

	"git.eaciitapp.com/ariefdarmawan/rpchub"
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
			err := hub.Call(rpchub.Request{Name: "obj.Float", Parm: []interface{}{float64(800.50)}}, res)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res.Err(), convey.ShouldBeNil)
			convey.Printf("\nData: %v\n", res.Data)

			float := rpchub.BytesToFloat(res.Data)
			convey.So(float, convey.ShouldEqual, float64(2)*800.50)
		})
	})
}
