package match

func matchWildcardSimple(pattern, data string) bool {
	if pattern == "" {
		return data == pattern
	}
	if pattern == "*" {
		return true
	}
	return deepMatch([]rune(data), []rune(pattern), true)
}

func matchWildcardAdvanced(pattern, data string) (matched bool) {
	if pattern == "" {
		return data == pattern
	}
	if pattern == "*" {
		return true
	}
	return deepMatch([]rune(data), []rune(pattern), false)
}

func deepMatch(str, pattern []rune, simple bool) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 && !simple {
				return false
			}
		case '*':
			return deepMatch(str, pattern[1:], simple) ||
				(len(str) > 0 && deepMatch(str[1:], pattern, simple))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(str) == 0 && len(pattern) == 0
}
