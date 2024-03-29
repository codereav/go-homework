package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TnClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type TnClient struct {
	in      io.ReadCloser
	out     io.Writer
	address string
	timeout time.Duration
	conn    net.Conn
}

func (c *TnClient) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	return nil
}

func (c *TnClient) Send() error {
	_, err := io.Copy(c.conn, c.in)

	return err
}

func (c *TnClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)

	return err
}

func (c *TnClient) Close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
	}
	if c.in != nil {
		err = c.in.Close()
	}
	return err
}
