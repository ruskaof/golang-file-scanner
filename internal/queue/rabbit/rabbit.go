package rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbit(
	username,
	password,
	host,
	port,
	queueName string,
) (*amqp.Connection, *amqp.Channel, amqp.Queue, error) {
	conn, err := setupConnection(username, password, host, port)

	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	channel, err := conn.Channel()

	if err != nil {
		_ = conn.Close()
		return nil, nil, amqp.Queue{}, err
	}

	queue, err := declareQueue(channel, queueName)

	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	return conn, channel, queue, nil
}

func setupConnection(username string, password string, host string, port string) (*amqp.Connection, error) {
	url := fmt.Sprintf(
		"amqp://%v:%v@%v:%v/",
		username,
		password,
		host,
		port,
	)

	conn, err := amqp.Dial(url)
	return conn, err
}

func declareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	return queue, err
}
