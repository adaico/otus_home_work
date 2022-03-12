package main

import (
	"context"
	"fmt"
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

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (tc *telnetClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
	defer cancel()

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", tc.address)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	tc.conn = conn

	return nil
}

func (tc *telnetClient) Receive() error {
	if tc.conn == nil {
		return nil
	}

	_, err := io.Copy(tc.out, tc.conn)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	return nil
}

func (tc *telnetClient) Send() error {
	if tc.conn == nil {
		return nil
	}

	if _, err := io.Copy(tc.conn, tc.in); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (tc *telnetClient) Close() error {
	if tc.conn == nil {
		return nil
	}

	if err := tc.conn.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address,
		timeout,
		in,
		out,
		nil,
	}
}
