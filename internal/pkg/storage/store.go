package storage

type Store interface {
	Ping() error
}
