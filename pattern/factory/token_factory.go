package factory

type TokenGenerator interface {
	GenerateToken() (token string)
}

type AdminTokenGenerator struct {
}

type CustomerTokenGenerator struct {
}

func (at AdminTokenGenerator) GenerateToken() (token string) {
	return
}

func (at CustomerTokenGenerator) GenerateToken() (token string) {
	return
}
