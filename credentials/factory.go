package credentials

type Factory struct {
	Providers map[string]Credentials
}

func (f Factory) Get(environment string) Credentials {
	return f.Providers[environment]
}
