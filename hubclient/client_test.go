package hubclient_test

import (
	"testing"

	"github.com/ariefdarmawan/rpchub/hubclient"
	"github.com/ariefdarmawan/rpchub/hubserver"

	"github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	convey.Convey("server", t, func() {
		s := hubserver.NewServer().Register(newObj("Hub"))
		err := s.Start(":9910")
		convey.So(err, convey.ShouldBeNil)
		defer s.Stop()

		convey.Convey("ping", func() {
			client, err := hubclient.NewClient(":9910")
			convey.So(err, convey.ShouldBeNil)
			defer client.Close()

			res := client.Call("obj.ping")
			convey.So(res.Err(), convey.ShouldBeNil)
			convey.So(string(res.Data), convey.ShouldEqual, "RPC Object OK")
		})
	})
}
