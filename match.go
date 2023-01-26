package match

type Matcher interface {
	Matches(data string) bool
}

type matcher struct {
	prepared []prepared
}

// Compile takes a pattern compiles it and if its valid reuturns a matcher if not it returns an error
//
// # A pattern can look like the following:
//
// djakwndaw[ m* | t* ]tewaljdm[ test | zetto ]rest
//
// ...
func Compile(pattern string) (Matcher, error) {
	parts, err := parseQueryIntoParts(pattern)
	if err != nil {
		return matcher{}, err
	}
	parts = parsePatterns(parts)

	prefix, patterns, suffix := extractPreAndSuffixFromParts(parts)
	cartesianProduct := generateCartesianProduct(patterns)
	preparedMatcher := combineFixData(prefix, suffix, cartesianProduct)
	preparedMatcher = orderPreparedByComplexity(preparedMatcher)

	return matcher{
		prepared: preparedMatcher,
	}, nil
}

func (m matcher) Matches(data string) bool {
	return matchMulti(m.prepared, data)
}
