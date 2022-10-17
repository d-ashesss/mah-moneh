package accounts

// Accounts is a service responsible for managing accounts.
type Accounts struct {
	db AccountStore
}

// NewService initializes a new accounts service.
func NewService(db AccountStore) Accounts {
	return Accounts{db: db}
}
