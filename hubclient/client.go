package hubclient

import (
	"fmt"
	"net/rpc"

	"github.com/ariefdarmawan/rpchub"
)

type client struct {
	dialer *rpc.Client
}

func NewClient(addr string) (*client, error) {
	c := new(client)

	d, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c.dialer = d

	return c, nil
}

func (c *client) Call(name string, parms ...interface{}) *rpchub.Response {
	res := new(rpchub.Response)
	if c.dialer == nil {
		return rpchub.NewResponseWithErr("client is not connected")
	}

	if err := c.dialer.Call("hub.Call", rpchub.NewRequest(name, parms...), res); err != nil {
		return rpchub.NewResponseWithErr(fmt.Sprintf("error invoking. %s", err.Error()))
	}

	return res
}

func (c *client) Close() {
	if c.dialer != nil {
		c.dialer.Close()
	}
}
