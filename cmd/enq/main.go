package main

import (
	"flag"
	"fmt"

	"github.com/knry0329/rabbitmq_enq_deq/service"
)

func main() {
	fmt.Println("enqueue start.")
	flag.Parse()
	// new enqueuer
	enq := service.NewEnqueuer(
		"amqp://guest:guest@localhost:5672/",
		"knry_q",
	)
	if err := enq.Enqueue([]byte(flag.Arg(0))); err != nil {
		fmt.Println("enqueue failure.")
		return
	}
	fmt.Println("enqueue success.")
}
