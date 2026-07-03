package scripts

func PerformCollectRecycSlow(args args, apply applyFn) {
	apply("collect_type", "RECYC")
	apply("collect_value", 24)
}
