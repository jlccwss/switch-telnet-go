package telnet

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	TIME_DELAY_AFTER_WRITE = 10
)

type Client struct {
	Address string
	Conn    net.Conn
	buf     [4096]byte
}

func (c *Client) Write(conn net.Conn, bufs []byte) (n int, err error) {
	n, err = conn.Write(bufs)
	if err != nil {
		return n, err
	}
	time.Sleep(time.Millisecond * TIME_DELAY_AFTER_WRITE)
	return n, err
}

func (c *Client) Connect(address string) (err error) {
	c.Conn, err = net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return err
	}
	c.Conn.SetDeadline(time.Now().Add(5 * time.Second))

	for {
		n, err := c.Conn.Read(c.buf[0:])
		if err != nil {
			break
		}
		if strings.Contains(string(c.buf[0:n]), "Username:") {
			break
		}
	}
	c.Conn.SetDeadline(time.Now().Add(15 * time.Second))

	return err
}

func (c *Client) Login(username string, password string) error {
	n, err := c.Write(c.Conn, []byte(username+"\n"))
	if err != nil {
		return err
	}

	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}

	n, err = c.Write(c.Conn, []byte(password+"\n"))
	if err != nil {
		return err
	}
	n, err = c.Conn.Read(c.buf[0:])
	if err != nil {
		return err
	}
	fmt.Printf(" login end %s\n", string(c.buf[0:n]))
	/*
		n, err = c.Write(c.Conn, []byte("enable\n"))
		if err != nil {
			return err
		}

		n, err = c.Conn.Read(c.buf[0:])
		if err != nil {
			return err
		}
		//fmt.Println(string(buf[0:n]))

		n, err = c.Write(c.Conn, []byte(enable+"\n"))
		if err != nil {
			return err
		}

		n, err = c.Conn.Read(c.buf[0:])
		if err != nil {
			return err
		}
		//fmt.Println(string(buf[0:n]))

		n, err = c.Write(c.Conn, []byte("terminal length 0\n"))
		if err != nil {
			return err
		}

		n, err = c.Conn.Read(c.buf[0:])
		if err != nil {
			return err
		}
		//fmt.Println(string(buf[0:n]))
	*/
	return err
}

func (c *Client) Cmd(shell string) (context string, err error) {
	_, err = c.Write(c.Conn, []byte(shell+"\n"))
	if err != nil {
		return "", err
	}
	for {
		n, err := c.Conn.Read(c.buf[0:])
		if err != nil {
			break
		}
		context += string(c.buf[0:n])
		if strings.HasSuffix(string(c.buf[0:n]), ">") || strings.HasSuffix(string(c.buf[0:n]), "]") || strings.HasSuffix(string(c.buf[0:n]), "#") || strings.HasSuffix(string(c.buf[0:n]), "Password: ") {
			break
		}
	}
	return context, err
}
