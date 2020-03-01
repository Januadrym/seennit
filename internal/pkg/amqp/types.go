package amqp

import "github.com/streadway/amqp"

type (
	Exchange struct {
		name                                 string
		kind                                 string
		duable, autoDelete, internal, noWait bool
		args                                 amqp.Table
	}

	Queue struct {
		name                                  string
		duable, autoDelete, exclusive, noWait bool
		arge                                  amqp.Table
	}

	connect = *amqp.Connection
	channel = *amqp.Channel

	AMQPObject struct {
		conn connect
		ch   channel
	}
)

const (
	MyExchange     = "ex-change"
	MyQueue        = "q-ue-ue"
	MyExchangeType = "direct"
)
