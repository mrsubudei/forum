package hasher

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type BcryptHasher struct {
	raw string
}

func NewBcryptHasher(raw string) *BcryptHasher {
	return &BcryptHasher{raw: raw}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	//
	var err error
	return "", err
}
