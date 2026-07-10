package scripts

func PerformCollectHeadMin(args args, apply applyFn) {
	apply("collect_type", "RECYC")
	apply("collect_value", 5)
}
