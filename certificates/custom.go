package certificates

const Custom = "custom"

type CustomStrategy struct {
	CertificatePath    string
	CertificateKeyPath string
}

func (c *CustomStrategy) Request(params RequestParams) (RequestResponse, error) {
	return RequestResponse{}, nil
}

func (c *CustomStrategy) Deposit(params DepositParams) (DepositResponse, error) {
	return DepositResponse{}, nil
}

func (c *CustomStrategy) GetLocations(params GetLocationParams) (Location, error) {
	return Location{
		CertPath: c.CertificatePath,
		KeyPath:  c.CertificateKeyPath,
	}, nil
}
