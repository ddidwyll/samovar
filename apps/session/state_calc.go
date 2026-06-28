package session

import (
	clc "samovar/lib/calc"
)

func testRecordCollectionByType(r clc.Report, _ clc.FetchFn, apply clc.ApplyFn) {
	if r.OldValue.IsNil() {
		return
	}
}
