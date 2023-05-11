package rabbit

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

// FileMessageQueue is a rabbitmq implementation of queue.FileMessageQueue that is used in the application
type FileMessageQueue struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func (q FileMessageQueue) StartConsumer(
	handle func(filePath string) error,
) error {
	messages, err := q.Channel.Consume(
		q.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for message := range messages {
			var strMessage = string(message.Body)
			log.Printf("rabbit received message: %s", strMessage)
			err = handle(strMessage)
			if err != nil {
				log.Printf("could not preprocess message: %v", err)
			}
		}
	}()

	log.Printf("rabbit waitinig for messages")

	return nil
}

func (q FileMessageQueue) Produce(filename string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := q.Channel.PublishWithContext(ctx,
		"",
		q.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(filename),
		})

	if err != nil {
		return err
	}

	return nil

}
