package exceptions

type DbNameMissingError struct {
	Key string
}

func (e *DbNameMissingError) Error() string {
	return "key " + e.Key + " is missing"
}
