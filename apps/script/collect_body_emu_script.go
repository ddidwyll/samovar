package script

func performCollectBodyEmu(args args, apply applyFn) {
	apply("collect_type", "HEAD")
	apply("collect_value", 10)
}
