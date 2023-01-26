package match

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQueryIntoParts(t *testing.T) {
	query := "[ * ]test[ wild1 | wild2 | wil?4 | wi*ld ]next[ * ]"
	expected := []part{
		{
			static:  false,
			content: "*",
		},
		{
			static:  true,
			content: "test",
		},
		{
			static:  false,
			content: "wild1 | wild2 | wil?4 | wi*ld",
		},
		{
			static:  true,
			content: "next",
		},
		{
			static:  false,
			content: "*",
		},
	}
	actual, err := parseQueryIntoParts(query)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	notParseableQuery := "test wild1 | wild2 | wil?4 | wi*ld ]next[ * ]"
	_, err = parseQueryIntoParts(notParseableQuery)
	assert.Error(t, err)
}

func TestParsePatterns(t *testing.T) {
	query := "test[ wild1 | wild2 | wil?4 | wi*ld ]next[ * ]"
	parts, err := parseQueryIntoParts(query)
	assert.NoError(t, err)

	expected := []part{
		{
			static:  true,
			content: "test",
		},
		{
			static:  false,
			content: "wild1 | wild2 | wil?4 | wi*ld",
			patterns: []string{
				"wild1",
				"wild2",
				"wil?4",
				"wi*ld",
			},
		},
		{
			static:  true,
			content: "next",
		},
		{
			static:   false,
			content:  "*",
			patterns: []string{"*"},
		},
	}

	actual := parsePatterns(parts)
	assert.Equal(t, expected, actual)
}

func TestExtractPreAndSuffix(t *testing.T) {
	{
		query := "test[ wild1 | wild2 | wil?4 | wi*ld ]next[ * ]suff"
		parts, err := parseQueryIntoParts(query)
		assert.NoError(t, err)
		parts = parsePatterns(parts)

		pre, rest, suf := extractPreAndSuffixFromParts(parts)

		assert.Equal(t, "test", pre)
		assert.Equal(t, []part{
			{
				static:  false,
				content: "wild1 | wild2 | wil?4 | wi*ld",
				patterns: []string{
					"wild1",
					"wild2",
					"wil?4",
					"wi*ld",
				},
			},
			{
				static:  true,
				content: "next",
			},
			{
				static:   false,
				content:  "*",
				patterns: []string{"*"},
			},
		}, rest)
		assert.Equal(t, "suff", suf)
	}
	{
		query := "[ wild1 | wild2 | wil?4 | wi*ld ]next[ * ]"
		parts, err := parseQueryIntoParts(query)
		assert.NoError(t, err)
		parts = parsePatterns(parts)

		pre, rest, suf := extractPreAndSuffixFromParts(parts)

		assert.Equal(t, "", pre)
		assert.Equal(t, []part{
			{
				static:  false,
				content: "wild1 | wild2 | wil?4 | wi*ld",
				patterns: []string{
					"wild1",
					"wild2",
					"wil?4",
					"wi*ld",
				},
			},
			{
				static:  true,
				content: "next",
			},
			{
				static:   false,
				content:  "*",
				patterns: []string{"*"},
			},
		}, rest)
		assert.Equal(t, "", suf)
	}

}

func TestGenerateCartesianProduct(t *testing.T) {
	query := "test[ wild1 | wild2 | wil?4 | wi*ld ]next[ * ]"
	parts, err := parseQueryIntoParts(query)
	assert.NoError(t, err)
	parts = parsePatterns(parts)
	_, rest, _ := extractPreAndSuffixFromParts(parts)

	actual := generateCartesianProduct(rest)

	expected := [][]string{
		{
			"wild1",
			"next",
			"*",
		},
		{
			"wild2",
			"next",
			"*",
		},
		{
			"wil?4",
			"next",
			"*",
		},
		{
			"wi*ld",
			"next",
			"*",
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestExtractPrefixAndSuffixFromProduct(t *testing.T) {
	{
		product := []string{
			"wil?4",
			"next",
			"*",
		}

		pre, rest, suf := extractPrefixAndSuffixFromProduct(product)

		assert.Equal(t, "wil", pre)
		assert.Equal(t, "?4next*", rest)
		assert.Equal(t, "", suf)
	}
	{
		product := []string{
			"wil?4",
			"next",
			"*",
			"suf",
		}

		pre, rest, suf := extractPrefixAndSuffixFromProduct(product)

		assert.Equal(t, "wil", pre)
		assert.Equal(t, "?4next*", rest)
		assert.Equal(t, "suf", suf)
	}
}

func TestCombineFixData(t *testing.T) {
	query := "test[ wild1 | wild2 | wil?4 | wi*ld ]next[ * ]"
	parts, err := parseQueryIntoParts(query)
	assert.NoError(t, err)
	parts = parsePatterns(parts)
	prefix, rest, suffix := extractPreAndSuffixFromParts(parts)
	product := generateCartesianProduct(rest)
	actual := combineFixData(prefix, suffix, product)
	expected := []prepared{
		{
			prefix:          "testwild1next",
			prefixLen:       len("testwild1next"),
			pattern:         "*",
			suffix:          "",
			suffixLen:       len(""),
			advancedPattern: true,
		},
		{
			prefix:          "testwild2next",
			prefixLen:       len("testwild2next"),
			pattern:         "*",
			suffix:          "",
			suffixLen:       len(""),
			advancedPattern: true,
		},
		{
			prefix:          "testwil",
			prefixLen:       len("testwil"),
			pattern:         "?4next*",
			suffix:          "",
			suffixLen:       len(""),
			advancedPattern: true,
		},
		{
			prefix:          "testwi",
			prefixLen:       len("testwi"),
			pattern:         "*ldnext*",
			suffix:          "",
			suffixLen:       len(""),
			advancedPattern: true,
		},
	}
	assert.ElementsMatch(t, expected, actual)
}

func TestCalculateComplexityOfPrepared(t *testing.T) {
	{
		p := prepared{
			prefix:          "testwild1next",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		}
		complexity := calculateComplexityOfPrepared(p)
		assert.Equal(t, 13, complexity)
	}
	{
		p := prepared{
			prefix:          "testwild1next",
			pattern:         "?",
			suffix:          "",
			advancedPattern: false,
		}
		complexity := calculateComplexityOfPrepared(p)
		assert.Equal(t, 51-13, complexity)
	}
}

func TestOrderByPreparedComplexity(t *testing.T) {
	expected := []prepared{
		{
			prefix:          "testwi",
			pattern:         "*ldnext*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwil",
			pattern:         "?4next*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwild2next2",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwild1next",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		},
	}

	shuffled := []prepared{
		{
			prefix:          "testwil",
			pattern:         "?4next*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwild1next",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwi",
			pattern:         "*ldnext*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwild2next2",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		},
	}

	assert.Equal(t, expected, orderPreparedByComplexity(shuffled))
}

func TestHasPrefix(t *testing.T) {
	assert.True(t, strings.HasPrefix("djawmkd", ""))
}

func TestMatchSingle(t *testing.T) {
	{
		p := prepared{
			prefix:          "testwild1next",
			prefixLen:       len("testwild1next"),
			pattern:         "*",
			suffix:          "",
			suffixLen:       len(""),
			advancedPattern: true,
		}

		assert.Equal(t, true, matchSingle(p, "testwild1nextblablabla", len("testwild1nextblablabla")))
		assert.Equal(t, false, matchSingle(p, "test1wild1nextblablabla", len("test1wild1nextblablabla")))
	}
	{
		p := prepared{
			prefix:          "testwild1next",
			prefixLen:       len("testwild1next"),
			pattern:         "?",
			suffix:          "",
			suffixLen:       len(""),
			advancedPattern: false,
		}

		assert.Equal(t, true, matchSingle(p, "testwild1next1", len("testwild1next1")))
		assert.Equal(t, false, matchSingle(p, "test1wild1nextblablabla", len("test1wild1nextblablabla")))
	}
}

func TestMatchMulti(t *testing.T) {
	ps := []prepared{
		{
			prefix:          "testwild1next",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwild2next",
			pattern:         "*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwil",
			pattern:         "?4next*",
			suffix:          "",
			advancedPattern: true,
		},
		{
			prefix:          "testwi",
			pattern:         "*ldnext*",
			suffix:          "",
			advancedPattern: true,
		},
	}
	assert.Equal(t, true, matchMulti(ps, "testwild1nextblablabla"))
	assert.Equal(t, false, matchMulti(ps, "test1wild1nextblablabla"))
}
