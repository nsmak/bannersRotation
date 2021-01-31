package rabbit

import (
	"fmt"

	"github.com/nsmak/bannersRotation/cmd/config"
	"github.com/streadway/amqp"
)

type Producer struct {
	cfg     config.Rabbit
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(cfg config.Rabbit) (*Producer, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, newError("can't connect to rmq", err)
	}
	return &Producer{conn: conn, cfg: cfg}, nil
}

func (p *Producer) CloseConn() error {
	return p.conn.Close()
}

func (p *Producer) Publish(body []byte) error {
	if p.channel == nil {
		return ErrChannelIsNil
	}

	err := p.channel.Publish(
		p.cfg.ExchangeName,
		p.cfg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			Body:            body,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	if err != nil {
		return newError("can't publish", err)
	}

	return nil
}

func (p *Producer) OpenChannel() error {
	var err error
	p.channel, err = declareChannel(p.cfg, p.conn)
	if err != nil {
		return newError("can't create channel", err)
	}
	return nil
}

func (p *Producer) CloseChannel() error {
	return p.channel.Close()
}
