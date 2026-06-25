package script

func performCollectBodyEmu(args args, apply applyFn) {
	apply("collect_type", "BODY")
	apply("collect_value", 10)
}
