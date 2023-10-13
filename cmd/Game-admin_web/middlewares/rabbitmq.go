package middlewares

import (
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

// 每次都建立新连接，不搞复用了

func InitRabbitmq() {

}

func Publisher(exChangeName string, kind string, message []byte, routerKey []string) error {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		zap.S().Info(err)
		return err

	}
	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exChangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind(kind),
	)
	if err != nil {
		zap.S().Info(err)
		return err
	}
	defer publisher.Close()

	err = publisher.Publish(
		message,
		routerKey,
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exChangeName),
	)
	if err != nil {
		zap.S().Info(err)
		return err
	}
	return nil
}

func Consumer(exChangeName string, queue string, routerKey string, work rabbitmq.Handler) error {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		zap.S().Fatal(err)
	}
	defer conn.Close()

	if work == nil {
		work = func(d rabbitmq.Delivery) rabbitmq.Action {
			zap.S().Info("consumed: %v", string(d.Body))

			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.NackRequeue
		}
	}
	consumer, err := rabbitmq.NewConsumer(
		conn,
		work,
		queue,
		rabbitmq.WithConsumerOptionsRoutingKey(routerKey),
		rabbitmq.WithConsumerOptionsExchangeName(exChangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		zap.S().Fatal(err)
		return err
	}
	defer consumer.Close()
	return nil
}
