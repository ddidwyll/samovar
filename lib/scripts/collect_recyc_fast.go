package scripts

func PerformCollectRecycFast(args args, apply applyFn) {
	apply("collect_type", "RECYC")
	apply("collect_value", 44)
}
