package nats

import "log"

func (stan *HandNats) StopProcess() {
	err := stan.Sub.sub.Unsubscribe()
	if err != nil {
		log.Fatalf("UNSUBSCRIDE ERROR: %s", err)
	}
	err = (*stan.Conn).Close()
	if err != nil {
		log.Fatalf("CLOSE CONNECTION ERROR: %s", err)
	}
	return
}
