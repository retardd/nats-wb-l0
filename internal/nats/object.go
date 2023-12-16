package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"l0/internal/datastorage/cache"
	"log"
	"os"
	"time"
)

type HandNats struct {
	Ident string
	Conn  *stan.Conn
	Pub   *Pub
	Sub   *Sub
}

func InitConnection(cch *cache.Cache) (*HandNats, error) {
	hn := HandNats{}
	hn.Ident = "test"
	time.Sleep(time.Second * 10)
	conn, err := stan.Connect(os.Getenv("NTS_CLUSTER"),
		os.Getenv("NTS_ID"),
		stan.NatsURL(fmt.Sprintf("nats://stan-server:%s", os.Getenv("NTS_HOST"))),
		stan.NatsOptions(
			nats.ReconnectWait(time.Second*20),
			nats.Timeout(time.Second*20),
			nats.ErrorHandler(func(conn *nats.Conn, sub *nats.Subscription, err error) {
				fmt.Printf("NATS error: %v\n", err)
			}),
		),
		stan.Pings(5, 3))
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return nil, err
	}

	hn.Conn = &conn

	hn.Sub = InitSub(cch)
	hn.Pub = InitPub(hn.Conn)

	hn.Sub.GetSubscribe(hn.Conn)

	return &hn, nil
}
