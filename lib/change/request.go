package change

func NewRequest(key string, val any, from string, ts int64) Request {
	return Request{key, val, []string{from}, ts}
}

type Request struct {
	Key       string
	Value     any
	From      []string
	Timestamp int64
}
