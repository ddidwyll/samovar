package scripts

func PerformCollectWaste(args args, apply applyFn) {
	apply("collect_type", "WASTE")
	apply("collect_value", 50)
}
