package script

import (
	clc "samovar/lib/calc"
)

type args = clc.Args
type applyFn = clc.ApplyFn

type scriptFn func(args, applyFn)
type scriptMap map[string]scriptFn
