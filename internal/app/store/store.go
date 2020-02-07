package store

// IfaceStore interface
type IfaceStore interface {
	User() UserRepository
}
