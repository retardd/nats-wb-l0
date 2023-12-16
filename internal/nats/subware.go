package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"l0/internal/datastorage/cache"
	"l0/internal/datastorage/structure"
	"log"
	"os"
)

type Sub struct {
	subj  string
	cache *cache.Cache
	sub   stan.Subscription
}

type Pub struct {
	scn  *stan.Conn
	subj string
}

func InitSub(cch *cache.Cache) *Sub {
	sub := Sub{}
	sub.subj = os.Getenv("NTS_PUB_SUBJ")
	sub.cache = cch
	return &sub
}

func InitPub(scn *stan.Conn) *Pub {
	pub := Pub{}
	pub.subj = os.Getenv("NTS_PUB_SUBJ")
	pub.scn = scn
	return &pub
}

func (sub *Sub) GetSubscribe(scn *stan.Conn) {
	var err error
	sub.sub, err = (*scn).Subscribe(sub.subj, sub.AddModel, stan.DurableName(os.Getenv("NTS_DNAME")))
	if err != nil {
		log.Fatalf("SUBSCRIDE ERROR: %s", err)
	}
}

func (sub *Sub) AddModel(mes *stan.Msg) {
	var tempModel structure.Model
	err := json.Unmarshal(mes.Data, &tempModel)

	if err != nil {
		log.Fatalf("IN HANDLING VALIDATION ERROR")
	}
	//КОНТЕКСТ БД
	err = sub.cache.AddData(context.TODO(), &tempModel)

	if err != nil {
		log.Fatalf("ADDING DATA ERROR")
	}
}

func (pub *Pub) GetPublish(model *structure.Model) error {
	data, err := json.Marshal(model)

	if err != nil {
		log.Fatalf("IN PUBLISHING VALIDATION ERROR")
	}
	var ackedNuid string

	ackedNuid, err = (*pub.scn).PublishAsync(pub.subj, data,
		func(ackedNuid string, err error) {
			if err == nil {
				log.Printf("GOOD ACK MESSAGE (%s)\n", ackedNuid)
			}
		})
	if err != nil {
		log.Fatalf("PUBLICATE ERROR (%s)", ackedNuid)
	}

	return err
}
