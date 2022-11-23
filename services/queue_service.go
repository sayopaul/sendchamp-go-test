package services

import (
	"encoding/json"
	"log"

	"github.com/sayopaul/sendchamp-go-test/config"
	"github.com/streadway/amqp"
)

type Queue struct {
	Conn         *amqp.Connection
	Channel      *amqp.Channel
	Name         string
	ExchangeName string
}
type QueueService struct {
	configEnv config.Config
}

func NewQueueService(configEnv config.Config) QueueService {
	return QueueService{
		configEnv: configEnv,
	}
}
func (qs QueueService) NewQueue(configEnv config.Config) *Queue {
	conn, err := amqp.Dial(configEnv.AMQPUrl)
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		configEnv.QueueName, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Panicf("%s: %s", "Failed to declare a queue", err)
	}

	erro := ch.ExchangeDeclare(
		configEnv.ExchangeName, // name
		"direct",               // exchangeType
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if erro != nil {
		log.Panicf("%s: %s", "Failed to declare an exchange", err)
	}

	errCh := ch.QueueBind(configEnv.QueueName, "", configEnv.ExchangeName, false, nil)
	if errCh != nil {
		log.Panicf("%s: %s", "Failed to bind the exchange and queue", err)
	}

	return &Queue{
		Conn:         conn,
		Channel:      ch,
		Name:         q.Name,
		ExchangeName: configEnv.ExchangeName,
	}
}

func (qs *QueueService) PublishMessage(message map[string]interface{}, configEnv config.Config) error {
	queue := qs.NewQueue(configEnv)
	toJson, _ := json.Marshal(message)

	err := queue.Channel.Publish(
		queue.ExchangeName, // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(toJson),
		})
	defer queue.Conn.Close()
	defer queue.Channel.Close()
	return err

}
