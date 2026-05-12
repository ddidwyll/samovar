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

func (r Request) LastFrom() string {
	length := len(r.From)

	if length == 0 {
		return ""
	} else {
		return r.From[length-1]
	}
}
