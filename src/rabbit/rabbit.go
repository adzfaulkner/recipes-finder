package rabbit

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type ClientInterface interface {
	Cleanup()
	Publish()
	Consume()
}

type RabbitClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func (c *RabbitClient) Cleanup() {
	c.Conn.Close()
	c.Channel.Close()
}

func (c *RabbitClient) Publish(msg []byte) {
	err := c.Channel.Publish(
		"",           // exchange
		c.Queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)

	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}

func (c *RabbitClient) Consume() <-chan amqp.Delivery {
	msgs, err := c.Channel.Consume(
		c.Queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	return msgs
}

func GetClient(n string) RabbitClient {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST")))

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	q, err := ch.QueueDeclare(
		n,     // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return RabbitClient{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}
}
