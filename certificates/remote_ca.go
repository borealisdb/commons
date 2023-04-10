package certificates

const RemoteCA = "remote_ca"

type RemoteCAStrategy struct {
}

func (r RemoteCAStrategy) Request(params RequestParams) (RequestResponse, error) {
	// TODO make a http call to the remote CA to request certificates
	return RequestResponse{}, nil
}

func (r RemoteCAStrategy) Deposit(params DepositParams) (DepositResponse, error) {
	return DepositResponse{}, nil
}

func (r RemoteCAStrategy) GetLocations(params GetLocationParams) (Location, error) {
	return Location{}, nil
}
