package mtdb

type AuthEntry struct {
	ID        *int64 `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	LastLogin int    `json:"last_login"`
}

type AuthRepository interface {
	GetByUsername(username string) (*AuthEntry, error)
	Create(entry *AuthEntry) error
	Update(entry *AuthEntry) error
	Delete(id int64) error
}
