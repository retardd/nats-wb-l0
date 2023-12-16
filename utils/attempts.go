package utils

import (
	"github.com/nats-io/stan.go"
	"time"
)

func TryToConn(some func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		err = some()
		if err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}

func TryToConnNats(some func() (stan.Conn, error), attempts int, delay time.Duration) (conn stan.Conn, err error) {
	for attempts > 0 {
		conn, err = some()
		if err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return conn, nil
	}
	return
}
