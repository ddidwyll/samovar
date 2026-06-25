package script

func performCollectBodyEmuY(args args, apply applyFn) {
	apply("collect_type", "HEAD")
	apply("collect_value", 30)
}
