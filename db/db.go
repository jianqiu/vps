package db

//go:generate counterfeiter . DB

type DB interface {
	VirtualGuestDB
}

