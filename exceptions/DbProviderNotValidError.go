package exceptions

type DbProviderNotValidError struct {
	Key   string
	Value interface{}
}

func (e *DbProviderNotValidError) Error() string {
	return "Db provider Type cast error for key: " + e.Key
}
