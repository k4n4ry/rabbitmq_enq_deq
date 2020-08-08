package service

import "github.com/streadway/amqp"

// EnqueuerIF is IF
type EnqueuerIF interface {
	Enqueue(data []byte) error
}

type enqueuer struct {
	address string
	qName   string
}

// NewEnqueuer is creating enqueuer
func NewEnqueuer(address, qName string) EnqueuerIF {
	return &enqueuer{
		address: address,
		qName:   qName,
	}
}

func (e *enqueuer) Enqueue(data []byte) error {
	conn, err := amqp.Dial(e.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		e.qName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	if err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "test/plain",
			Body:        data,
		},
	); err != nil {
		return err
	}
	return nil
}

// DequeuerIF is IF
type DequeuerIF interface {
	Dequeue() (<-chan amqp.Delivery, error)
	Dial() error
	Close()
}

type dequeuer struct {
	address string
	qName   string
	conn    *amqp.Connection
	ch      *amqp.Channel
}

// NewDequeuer is creating dequeuer
func NewDequeuer(address, qName string) (DequeuerIF, error) {
	// newのときにdialするか、Dial()という別メソッドを設けるか。
	var deq = dequeuer{
		address: address,
		qName:   qName,
	}
	if err := deq.Dial(); err != nil {
		return nil, err
	}
	return &deq, nil
}

func (d *dequeuer) Close() {
	defer d.conn.Close()
	defer d.ch.Close()
}

// reDial用に公開しておく
func (d *dequeuer) Dial() error {
	conn, err := amqp.Dial(d.address)
	if err != nil {
		return err
	}
	d.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	d.ch = ch
	return nil

}

// consumeを実行して
func (d *dequeuer) Dequeue() (<-chan amqp.Delivery, error) {
	q, err := d.ch.QueueDeclare(
		d.qName,
		false,
		false,
		false,
		false,
		nil,
	)

	msgs, err := d.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
