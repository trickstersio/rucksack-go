package secrets

func NewNoop() *Noop {
	return &Noop{}
}

type Noop struct{}

func (s *Noop) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (s *Noop) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}
