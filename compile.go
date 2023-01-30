package match

import (
	"errors"
	"sort"
	"strings"
)

type part struct {
	static   bool
	content  string
	patterns []string
}

type prepared struct {
	prefix  string
	pattern string
	suffix  string

	prefixLen int
	suffixLen int

	advancedPattern bool
}

func parseQueryIntoParts(query string) ([]part, error) {
	parts := []part{}
	cp := part{static: query[0] != '['}
	for i, r := range query {
		switch r {
		case '[':
			if i == 0 {
				continue
			}
			cp.content = strings.TrimSpace(cp.content)
			parts = append(parts, cp)
			cp = part{static: false}
		case ']':
			if cp.static {
				return parts, errors.New("invalid query: expected '[' before ']'")
			}
			cp.content = strings.TrimSpace(cp.content)
			parts = append(parts, cp)
			if len(query) > i+1 {
				cp = part{static: query[i+1] != '['}
			} else {
				cp = part{}
			}
		default:
			cp.content += string(r)
		}
	}
	if len(cp.content) > 0 {
		cp.content = strings.TrimSpace(cp.content)
		parts = append(parts, cp)
	}
	return parts, nil
}

func parsePatterns(parts []part) []part {
	for i := range parts {
		if parts[i].static {
			continue
		}
		for _, pattern := range strings.Split(parts[i].content, "|") {
			parts[i].patterns = append(parts[i].patterns, strings.TrimSpace(pattern))
		}
	}
	return parts
}

func extractPrefixFromParts(parts []part) (string, []part) {
	if len(parts) == 0 || !parts[0].static {
		return "", parts
	}
	return parts[0].content, parts[1:]
}

func extractSuffixFromParts(parts []part) (string, []part) {
	if len(parts) == 0 || !parts[len(parts)-1].static {
		return "", parts
	}
	return parts[len(parts)-1].content, parts[0 : len(parts)-1]
}

func extractPreAndSuffixFromParts(parts []part) (string, []part, string) {
	pre, rest := extractPrefixFromParts(parts)
	suf, rest := extractSuffixFromParts(rest)
	return pre, rest, suf
}

func prepareForCartesianProduct(parts []part) [][]string {
	permuatable := [][]string{}
	for _, part := range parts {
		if part.static {
			permuatable = append(permuatable, []string{part.content})
		} else {
			permuatable = append(permuatable, part.patterns)
		}
	}
	return permuatable
}

func generateCartesianProduct(parts []part) [][]string {
	permutable := prepareForCartesianProduct(parts)
	return cartesian(permutable...)
}

func extractPrefixAndSuffixFromProduct(data []string) (string, string, string) {
	prefix := ""
	suffix := ""
	datastr := strings.Join(data, "")
	for i, r := range datastr {
		if r == '*' || r == '?' {
			datastr = datastr[i:]
			break
		}
		prefix += string(r)
	}
	lendatastr := len(datastr)
	for i := lendatastr - 1; i > 0; i-- {
		r := datastr[i]
		if r == '*' || r == '?' {
			datastr = datastr[:i+1]
			break
		}
		suffix = string(r) + suffix
	}
	return prefix, datastr, suffix
}

func combineFixData(prefix string, suffix string, cartesianProduct [][]string) []prepared {
	preparedData := make([]prepared, len(cartesianProduct))
	for i, product := range cartesianProduct {
		p, pattern, s := extractPrefixAndSuffixFromProduct(product)
		preparedData[i] = prepared{
			prefix:          prefix + p,
			prefixLen:       len(prefix + p),
			pattern:         pattern,
			suffix:          suffix + s,
			suffixLen:       len(suffix + s),
			advancedPattern: strings.Contains(pattern, "*"),
		}
	}
	return preparedData
}

func calculateComplexityOfPrepared(p prepared) int {
	// in case the pattern is of length 0 or is * it ts pattern complexity is 0
	patternComplexity := 0
	if p.pattern == "" || p.pattern == "*" {
		patternComplexity = 0
	} else if p.advancedPattern {
		advancedPatternBaseComplexity := 200
		// a advanced pattern is the most complex
		patternComplexity = advancedPatternBaseComplexity + len(p.pattern)
	} else {
		// a simple pattern is not really complex
		simplePatternBaseComplexity := 50
		patternComplexity = simplePatternBaseComplexity + len(p.pattern)
	}
	// if the pattern complexity is greater then 0 its a good thing to have
	// a large prefix. It reduces the amount of hits against the wildcard matches.
	// if the pattern complexity is smaller then 0 a prefix is not exactly expensive but adds a bit overhead
	// because of the constant prefix comparison
	prefixComplexity := 0
	if patternComplexity > 0 {
		// a large prefix is good
		prefixComplexity = len(p.prefix) * -1
	} else {
		// a large prefix is bad
		prefixComplexity = len(p.prefix) * 1
	}
	// the same goes for the suffix complexity
	suffixComplexity := 0
	if patternComplexity > 0 {
		// a large suffix is good
		suffixComplexity = len(p.suffix) * -1
	} else {
		// a large suffix is bad
		suffixComplexity = len(p.suffix) * 1
	}
	// we add up all complexities
	return prefixComplexity + patternComplexity + suffixComplexity
}

func orderPreparedByComplexity(ps []prepared) []prepared {
	// the most complex prepared should be first and the least complex one last
	sort.SliceStable(ps, func(i, j int) bool {
		return calculateComplexityOfPrepared(ps[i]) > calculateComplexityOfPrepared(ps[j])
	})
	return ps
}

func matchSingle(p prepared, data string, dataLen int) bool {
	if !strings.HasPrefix(data, p.prefix) || !strings.HasSuffix(data, p.suffix) {
		return false
	}
	if p.advancedPattern {
		return matchWildcardAdvanced(p.pattern, data[p.prefixLen:dataLen-p.suffixLen])
	}
	return matchWildcardSimple(p.pattern, data[p.prefixLen:dataLen-p.suffixLen])
}

func matchMulti(ps []prepared, data string) bool {
	i := len(ps) - 1
	l := len(data)
	for i != -1 {
		if matchSingle(ps[i], data, l) {
			return true
		}
		i--
	}
	return false
}
