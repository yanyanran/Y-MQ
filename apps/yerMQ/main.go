package main

import (
	"Y-MQ/yerMQ"
	"github.com/judwhite/go-svc"
	"log"
	"sync"
	"syscall"
)

type program struct {
	once  sync.Once
	yerMq *yerMQ.YERMQ
}

func (p *program) Init(environment svc.Environment) error {
	opts := yerMQ.NewOptions()
}

func (p *program) Start() error {

}

func (p *program) Stop() error {

}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatalf("[yerMQ]%s", err)
	}
}
