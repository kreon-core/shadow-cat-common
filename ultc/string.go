package ultc

import "unicode"

func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}

	return true
}
