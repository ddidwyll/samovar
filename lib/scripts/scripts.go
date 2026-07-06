package scripts

import (
	clc "samovar/lib/calc"

	"fmt"
	"math"
	"runtime"
)

type args = clc.Args
type applyFn = clc.ApplyFn

func debug(some any, title string) {
	fmt.Printf(">>>>>>>>>>>>>>> %s <<<<<<<<<<<<<<<<<\n", title)
	fmt.Printf("%#v\n", some)

	if pc, _, _, ok := runtime.Caller(1); !ok {
		fmt.Println("<<<<<<<<<<<<<>>>>>>>>>>>>>>")
	} else {
		parentName := runtime.FuncForPC(pc).Name()
		fmt.Printf("< %s >\n", parentName)
	}
}

func calcCollectGHByRefluxRatio(power, refluxRatio float64) int64 {
	return int64(power * 3600 / ((refluxRatio + 1) * 911))
}

func calcCollectPercByCollectGH(collect, maxCollect int64) int64 {
	return int64(math.Floor(float64(collect) / float64(maxCollect) * 100.0))
}
