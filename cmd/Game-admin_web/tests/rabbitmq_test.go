package tests

import (
	"github.com/wagslane/go-rabbitmq"
	"log"
	"testing"
	"time"
)

func publisher() {
	conn, err := rabbitmq.NewConn("amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	newPublisher, err := rabbitmq.NewPublisher(conn, rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	defer newPublisher.Close()

	for {
		time.Sleep(time.Second * 3)
		newPublisher.Publish([]byte("hello rabbitmq"),
			[]string{"mykey"},
			rabbitmq.WithPublishOptionsContentType("app"),
			rabbitmq.WithPublishOptionsExchange("evebts"))
	}
}

func consumer() {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.NackRequeue
		},
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
}

func TestRabbitmq(t *testing.T) {
	go publisher()
	go func() {
		now := time.Now()
		for true {
			if time.Since(now) == time.Second*10 {
				break
			}
			//consumer()
			time.Sleep(time.Second * 2)
		}
	}()
	time.Sleep(time.Hour)
}
