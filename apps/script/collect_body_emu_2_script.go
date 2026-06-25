package script

func performCollectBodyEmuX(args args, apply applyFn) {
	apply("collect_type", "HEAD")
	apply("collect_value", 50)
}
