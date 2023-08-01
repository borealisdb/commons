package credentials

type Factory struct {
	Providers map[string]Credentials
}

func (f Factory) Get(providerName string) Credentials {
	return f.Providers[providerName]
}
