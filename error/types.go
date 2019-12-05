package error

import "fmt"

// MongoError struct
type MongoError struct {
	err  string
	more string
	// query bson.M
}

func (e *MongoError) Error() string {
	return fmt.Sprintf("%s: %s", e.more, e.err)
}

// StdError struct
type StdError struct {
	err  string
	more string
}

func (e *StdError) Error() string {
	return fmt.Sprintf("%s: %s", e.more, e.err)
}
