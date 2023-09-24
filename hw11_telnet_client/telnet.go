package main

import (
	"io"
	"net"
	"time"
)

const bufferSize = 4096

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
	buffer := make([]byte, bufferSize)
	i, err := c.in.Read(buffer)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(buffer[:i])

	return err
}

func (c *TnClient) Receive() error {
	buffer := make([]byte, bufferSize)
	i, err := c.conn.Read(buffer)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buffer[:i])

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
