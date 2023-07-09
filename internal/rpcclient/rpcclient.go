package rpcclient

import (
	"net/rpc"

	"github.com/cenkalti/backoff/v4"
)

type client struct {
	rpcClient *rpc.Client
}

type Client interface {
	Call(method string, args any, reply any) error
}

func NewRPCClient(address string) (Client, error) {
	var rpcClient *rpc.Client
	operation := func() error {
		r, err := rpc.DialHTTP("tcp", address)
		if err != nil {
			return err
		}

		rpcClient = r
		return nil
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	return client{rpcClient: rpcClient}, nil
}

func (c client) Call(method string, args any, reply any) error {
	operation := func() error {
		return c.rpcClient.Call(method, args, reply)
	}
	return backoff.Retry(operation, backoff.NewExponentialBackOff())
}
