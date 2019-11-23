package cony

import "github.com/streadway/amqp"

// Declaration is a callback type to declare AMQP queue/exchange/binding
type Declaration func(Declarer) error

// Declarer is implemented by *amqp.Channel
type Declarer interface {
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Confirm(noWait bool) error
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(c chan amqp.Return) chan amqp.Return
}

// DeclareQueue is a way to declare AMQP queue
func DeclareQueue(q *Queue) Declaration {
	name := q.Name
	return func(c Declarer) error {
		q.Name = name
		realQ, err := c.QueueDeclare(q.Name,
			q.Durable,
			q.AutoDelete,
			q.Exclusive,
			false,
			q.Args,
		)
		q.l.Lock()
		q.Name = realQ.Name
		q.l.Unlock()
		return err
	}
}

// DeclareExchange is a way to declare AMQP exchange
func DeclareExchange(e Exchange) Declaration {
	return func(c Declarer) error {
		return c.ExchangeDeclare(e.Name,
			e.Kind,
			e.Durable,
			e.AutoDelete,
			false,
			false,
			e.Args,
		)
	}
}

// DeclareBinding is a way to declare AMQP binding between AMQP queue and exchange
func DeclareBinding(b Binding) Declaration {
	return func(c Declarer) error {
		return c.QueueBind(b.Queue.Name,
			b.Key,
			b.Exchange.Name,
			false,
			b.Args,
		)
	}
}

func DeclareConfirm(noWait bool) Declaration {
	return func(c Declarer) error {
		return c.Confirm(noWait)
	}
}

func DeclareNotifyPublish(confirm chan amqp.Confirmation) Declaration {
	return func(c Declarer) error {
		c.NotifyPublish(confirm)
		return nil
	}
}

func DeclareNotifyReturn(rc chan amqp.Return) Declaration {
	return func(c Declarer) error {
		c.NotifyReturn(rc)
		return nil
	}
}
