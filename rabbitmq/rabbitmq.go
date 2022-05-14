package rabbitmq

import (
	"github.com/ChromaMinecraft/email-service/config"
	"github.com/streadway/amqp"
)

type RabbitMQ interface {
	GetChannel() (*amqp.Channel, error)
	Subscribe(ch *amqp.Channel, q amqp.Queue, cb func(msg amqp.Delivery)) error
}

type rabbitMQ struct {
	AppConfig *config.Config
	AMQP      *amqp.Connection
}

func NewRabbitMQ(AppConfig *config.Config, AMQP *amqp.Connection) RabbitMQ {
	return &rabbitMQ{
		AppConfig: AppConfig,
		AMQP:      AMQP,
	}
}

func (r *rabbitMQ) GetChannel() (*amqp.Channel, error) {
	return r.AMQP.Channel()
}

func (r *rabbitMQ) Subscribe(ch *amqp.Channel, q amqp.Queue, cb func(msg amqp.Delivery)) error {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			cb(d)
		}
	}()

	return nil
}
