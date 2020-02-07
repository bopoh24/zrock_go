package store

// Store interface
type Store interface {
	User() UserRepository
}
