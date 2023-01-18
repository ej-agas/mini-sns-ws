package domain

type Transport interface {
	Send(to, message string) error
}
