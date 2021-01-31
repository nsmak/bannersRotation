package rabbit

import (
	"github.com/nsmak/bannersRotation/cmd/config"
	"github.com/nsmak/bannersRotation/internal/app"
	"github.com/streadway/amqp"
)

type mqError struct {
	app.BaseError
}

func (e *mqError) Error() string {
	if e.Err != nil {
		e.Message = "[rmq] " + e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}

func newError(msg string, err error) *mqError {
	return &mqError{BaseError: app.BaseError{Message: msg, Err: err}}
}

var ErrChannelIsNil = newError("channel is nil", nil)

func declareChannel(cfg config.Rabbit, conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, newError("can't get channel", err)
	}
	err = channel.ExchangeDeclare(
		cfg.ExchangeName,
		cfg.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, newError("can't declare exchange", err)
	}

	return channel, nil
}
