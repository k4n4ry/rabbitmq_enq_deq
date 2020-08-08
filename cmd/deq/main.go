package main

import (
	"fmt"
	"time"

	"github.com/knry0329/rabbitmq_enq_deq/service"
)

func main() {
	fmt.Println("dequeue start.")

	deq, err := service.NewDequeuer(
		"amqp://guest:guest@localhost:5672/",
		"knry_q",
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer deq.Close()
	forever := make(chan bool)
	msgs, err := deq.Dequeue()
	if err != nil {
		fmt.Println("dequeue failure. exit")
		return
	}
	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if ok {
					fmt.Printf("dequeue data. %s\n", msg.Body)
				} else {
					fmt.Println("channel closed. redialing ...")
					// リダイアル処理
					deq.Close()
					var successFlg bool
					for !successFlg {
						time.Sleep(5 * time.Second)
						if err := deq.Dial(); err != nil {
							fmt.Println("redial failure. retrying ...")
							continue
						}
						fmt.Println("redial success.")
						successFlg = true
					}
					msgs, err = deq.Dequeue()
					if err != nil {
						fmt.Println("dequeue failure. exit")
						return
					}
				}
			default:
				fmt.Println("sleep...")
				time.Sleep(5 * time.Second)
			}
		}
	}()
	<-forever
}
