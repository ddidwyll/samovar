package change

import (
)

type Request struct {
	Key       string
	Value     any
	Timestamp int64
}

type Report struct {
  Changed bool
  Key string
  NewValue string
  OldValue string
  Timestamp int64
}
