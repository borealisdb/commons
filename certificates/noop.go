package certificates

const Noop = "noop"

type NoopStrategy struct{}

func (n NoopStrategy) Request(params RequestParams) (RequestResponse, error) {
	return RequestResponse{}, nil
}

func (n NoopStrategy) Deposit(params DepositParams) (DepositResponse, error) {
	return DepositResponse{}, nil
}

func (n NoopStrategy) GetLocations(params GetLocationParams) (Location, error) {
	return Location{}, nil
}
