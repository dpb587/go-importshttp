package stringsutil

func Coalesce(inputs ...string) string {
	for _, input := range inputs {
		if len(input) > 0 {
			return input
		}
	}

	return ""
}
