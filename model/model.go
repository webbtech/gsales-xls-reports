package model

// DBHandler interface
type DBHandler interface {
	Close()
	// FetchQuote(string) (*Quote, error)
}
