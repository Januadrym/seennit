package amqp

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (q *AMQPObject) CreateExchange(ex Exchange) error {
	return q.ch.ExchangeDeclare(
		ex.name,
		ex.kind,
		ex.duable,
		ex.autoDelete,
		ex.internal,
		ex.noWait,
		ex.args,
	)
}

func (q *AMQPObject) CreateQueue(que Queue) (*amqp.Queue, error) {
	queue, err := q.ch.QueueDeclare(
		que.name,
		que.duable,
		que.autoDelete,
		que.exclusive,
		que.noWait,
		que.arge,
	)
	if err != nil {
		logrus.Errorf("failed to create message queue, err: %v", err)
		return nil, err
	}
	return &queue, nil
}

func (q *AMQPObject) DeleteQueue(queueName string) error {
	_, err := q.ch.QueueDelete(queueName, false, false, true)
	if err != nil {
		logrus.Errorf("faield to delete message queue, err %v", err)
		return err
	}
	return nil
}

func (q *AMQPObject) Subcribe(exchangeName, queueName string, routeKeys []string) {
	if exchangeName == "" {
		exchangeName = MyExchange
	}
	for _, routeKey := range routeKeys {
		if er := q.ch.QueueBind(queueName, routeKey, exchangeName, false, nil); er != nil {
			HandleError(er, "failed to bind queue with key: "+routeKey)
			logrus.Errorf("faield to bind queue %v with key %v, err: %v", queueName, routeKey, er)
		}
	}
}

func (q *AMQPObject) Publish(exchangeName string, routeKeys []string, data []byte) error {
	if exchangeName == "" {
		exchangeName = MyExchange
	}
	for _, routeKey := range routeKeys {
		err := q.ch.Publish(exchangeName, routeKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
		if err != nil {
			logrus.Errorf("failed to publish data to exchange %s with routekey %s, err: %v", exchangeName, routeKey, err)
		}
	}
	return nil
}

func (q *AMQPObject) Listen(queueName string, routeKeys []string) (<-chan []byte, error) {
	ex := Exchange{
		name:       MyExchange,
		kind:       MyExchangeType,
		duable:     true,
		autoDelete: false,
		internal:   false,
		noWait:     false,
		args:       nil,
	}
	err := q.CreateExchange(ex)
	if err != nil {
		logrus.Errorf("failed to create an Exchange, err: %v", err)
		HandleError(err, "cannot create Exchange"+ex.name)
	}
	que := Queue{
		name:       queueName,
		duable:     true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		arge:       nil,
	}
	queue, err := q.CreateQueue(que)
	if err != nil {
		logrus.Errorf("failed to create queue, err: %v", err)
		HandleError(err, "cannot create Queue with name: "+queue.Name)
	}

	q.Subcribe(MyExchange, queue.Name, routeKeys)

	msgs, er := q.ch.Consume(queue.Name, "", true, false, false, false, nil)
	if er != nil {
		logrus.Errorf("failed to register a consumer, err: %v", er)
		HandleError(er, "cannot register a consumer")
		return nil, er
	}
	return convertChannel(msgs), nil
}

func (q *AMQPObject) CloseMQ() {
	q.conn.Close()
	q.ch.Close()
}
