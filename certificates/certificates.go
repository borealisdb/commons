package certificates

// CertificateAuthority
// The idea is that in the future Borealis could have its own CA to manage, rotate and issue certificates centrally
type CertificateAuthority interface {
	// Request certificates from anywhere or even autogenerate them
	Request(params RequestParams) (RequestResponse, error)
	// Deposit Store them where ever you want
	Deposit(params DepositParams) (DepositResponse, error)
	// GetLocations get their location
	GetLocations(params GetLocationParams) (Location, error)
}

type Location struct {
	CertPath string
	KeyPath  string
}

type RequestParams struct {
}
type RequestResponse struct {
}

type DepositParams struct {
}
type DepositResponse struct {
}

type GetLocationParams struct {
}
