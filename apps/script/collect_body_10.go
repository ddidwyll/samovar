package script

func performCollectBody10(args args, apply applyFn) {
	apply("collect_type", "BODY")
	apply("collect_value", 10)
}
