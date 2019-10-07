package hubclient

import (
	"fmt"
	"net/rpc"

	"github.com/ariefdarmawan/rpchub"
)

type Client struct {
	dialer *rpc.Client
}

func NewClient(addr string) (*Client, error) {
	c := new(Client)

	d, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c.dialer = d

	return c, nil
}

func (c *Client) Call(name string, parms ...interface{}) *rpchub.Response {
	res := new(rpchub.Response)
	if c == nil {
		return rpchub.NewResponseWithErr("client is not yet initialized")
	}
	if c.dialer == nil {
		return rpchub.NewResponseWithErr("client is not connected")
	}

	if err := c.dialer.Call("hub.Call", rpchub.NewRequest(name, parms...), res); err != nil {
		return rpchub.NewResponseWithErr(fmt.Sprintf("error invoking. %s", err.Error()))
	}

	return res
}

func (c *Client) Close() {
	if c.dialer != nil {
		c.dialer.Close()
	}
}
