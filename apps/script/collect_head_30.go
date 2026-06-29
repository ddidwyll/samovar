package script

func performCollectHead30(args args, apply applyFn) {
	apply("collect_type", "HEAD")
	apply("collect_value", 30)
}
