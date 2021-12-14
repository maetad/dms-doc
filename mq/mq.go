package mq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"gitlab.com/pakkaparn/dms-doc/user"
	"gorm.io/gorm"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Received(db *gorm.DB) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"users_topic", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // queue name
		"user.*",      // routing key
		"users_topic", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(d.RoutingKey)
			switch d.RoutingKey {
			case "user.create":
				var body user.User
				json.Unmarshal([]byte(d.Body), &body)

				db.Create(&body)
				break
			case "user.update":
				var body user.User
				json.Unmarshal([]byte(d.Body), &body)

				db.Model(&user.User{}).Where("id = ?", body.ID).Updates(user.User{FirstName: body.FirstName, LastName: body.LastName})
				break
			case "user.delete":
				db.Delete(&user.User{}, d.Body)
				break
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
