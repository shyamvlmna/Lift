package handler

type DriverHandler interface {
}

type Impl struct {
}

func NewDriverHandler() DriverHandler {
	return &Impl{}
}