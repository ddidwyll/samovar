package script

func performIdle(args args, apply applyFn) {
	apply("collect_type", "OFF")
	apply("collect_value", 0)
}
