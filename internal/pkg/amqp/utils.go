package amqp

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func InitConnect(addr string) (*AMQPObject, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		logrus.Errorf("failed to connect to RabbitMQ, err: %v", err)
		return nil, err
	}

	ch, er := conn.Channel()
	if er != nil {
		logrus.Errorf("failed to create message queue channel , err: %v", err)
		return nil, er
	}
	return &AMQPObject{
		conn,
		ch,
	}, nil
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}

func convertChannel(msgs <-chan amqp.Delivery) <-chan []byte {
	channel := make(chan []byte)
	go func() {
		for d := range msgs {
			channel <- d.Body
		}
	}()
	return channel
}
