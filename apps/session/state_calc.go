package session

import (
	"fmt"
	clc "samovar/lib/calc"
)

func testRecordCollectionByType(r clc.Report, mustFetch clc.FetchFn, apply clc.ApplyFn) {
	isSynced := mustFetch("device_raw_state.collect_synced")
	fmt.Printf("### collect_type[%s]: %s => %s\n", isSynced, r.OldValue, r.NewValue)
}

func testRecordCollectionByValue(r clc.Report, mustFetch clc.FetchFn, apply clc.ApplyFn) {
	isSynced := mustFetch("device_raw_state.collect_synced")
	fmt.Printf("### collect_value[%s]: %s => %s\n", isSynced, r.OldValue, r.NewValue)
}
