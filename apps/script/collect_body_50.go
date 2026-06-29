package script

func performCollectBodyEmuX(args args, apply applyFn) {
	apply("collect_type", "BODY")
	apply("collect_value", 50)
}
