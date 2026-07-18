package scripts

func PerformCollectHeadMin(args args, apply applyFn) {
	apply("collect_type", "HEAD")
	apply("collect_value", 5)
}
