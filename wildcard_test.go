package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchWildcardAdvanced(t *testing.T) {
	testCases := []struct {
		pattern string
		text    string
		matched bool
	}{
		{
			pattern: "*",
			text:    "s3:GetObject",
			matched: true,
		},
		{
			pattern: "",
			text:    "s3:GetObject",
			matched: false,
		},
		{
			pattern: "",
			text:    "",
			matched: true,
		},
		{
			pattern: "s3:*",
			text:    "s3:ListMultipartUploadParts",
			matched: true,
		},
		{
			pattern: "someRandomLongSentence",
			text:    "s3:Listblabla",
			matched: false,
		},
		{
			pattern: "s3:Listblabla",
			text:    "s3:Listblabla",
			matched: true,
		},
		{
			pattern: "someRandomLongSentence",
			text:    "someRandomLongSentence",
			matched: true,
		},
		{
			pattern: "someStringoo*",
			text:    "someStringoo",
			matched: true,
		},
		{
			pattern: "someStringIn*",
			text:    "someStringIndia/Karnataka/",
			matched: true,
		},
		{
			pattern: "someStringIn*",
			text:    "someStringKarnataka/India/",
			matched: false,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Ban/Ban/Ban/Ban/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Area1/Area2/Area3/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/State1/State2/Karnataka/Area1/Area2/Area3/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Bangalore",
			matched: false,
		},
		{
			pattern: "someStringIn*/Ka*/Ban*",
			text:    "someStringIndia/Karnataka/Bangalore",
			matched: true,
		},
		{
			pattern: "someString*",
			text:    "someStringIndia",
			matched: true,
		},
		{
			pattern: "someStringoo*",
			text:    "someStringodo",
			matched: false,
		},
		{
			pattern: "my-blabla?/abc*",
			text:    "myblabla/abc",
			matched: false,
		},
		{
			pattern: "my-blabla?/abc*",
			text:    "my-blabla1/abc",
			matched: true,
		},
		{
			pattern: "my-?-blabla/abc*",
			text:    "my--blabla/abc",
			matched: false,
		},
		{
			pattern: "my-?-blabla/abc*",
			text:    "my-1-blabla/abc",
			matched: true,
		},
		{
			pattern: "my-?-blabla/abc*",
			text:    "my-k-blabla/abc",
			matched: true,
		},
		{
			pattern: "my??blabla/abc*",
			text:    "myblabla/abc",
			matched: false,
		},
		{
			pattern: "my??blabla/abc*",
			text:    "my4ablabla/abc",
			matched: true,
		},
		{
			pattern: "my-blabla?abc*",
			text:    "someStringabc",
			matched: false,
		},
		{
			pattern: "someStringabc?efg",
			text:    "someStringabcdefg",
			matched: true,
		},
		{
			pattern: "someStringabc?efg",
			text:    "someStringabc/efg",
			matched: true,
		},
		{
			pattern: "someStringabc????",
			text:    "someStringabc",
			matched: false,
		},
		{
			pattern: "someStringabc????",
			text:    "someStringabcde",
			matched: false,
		},
		{
			pattern: "someStringabc????",
			text:    "someStringabcdefg",
			matched: true,
		},
		{
			pattern: "someStringabc?",
			text:    "someStringabc",
			matched: false,
		},
		{
			pattern: "someStringabc?",
			text:    "someStringabcd",
			matched: true,
		},
		{
			pattern: "someStringabc?",
			text:    "someStringabcde",
			matched: false,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnop",
			matched: false,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnopqrst/mnopqr",
			matched: true,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnopqrst/mnopqrs",
			matched: true,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnop",
			matched: false,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnopq",
			matched: true,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnopqr",
			matched: true,
		},
		{
			pattern: "someStringmnop*?and",
			text:    "someStringmnopqand",
			matched: true,
		},
		{
			pattern: "someStringmnop*?and",
			text:    "someStringmnopand",
			matched: false,
		},
		{
			pattern: "someStringmnop*?and",
			text:    "someStringmnopqand",
			matched: true,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmn",
			matched: false,
		},
		{
			pattern: "someStringmnop*?",
			text:    "someStringmnopqrst/mnopqrs",
			matched: true,
		},
		{
			pattern: "someStringmnop*??",
			text:    "someStringmnopqrst",
			matched: true,
		},
		{
			pattern: "someStringmnop*qrst",
			text:    "someStringmnopabcdegqrst",
			matched: true,
		},
		{
			pattern: "someStringmnop*?and",
			text:    "someStringmnopqand",
			matched: true,
		},
		{
			pattern: "someStringmnop*?and",
			text:    "someStringmnopand",
			matched: false,
		},
		{
			pattern: "someStringmnop*?and?",
			text:    "someStringmnopqanda",
			matched: true,
		},
		{
			pattern: "someStringmnop*?and",
			text:    "someStringmnopqanda",
			matched: false,
		},
		{
			pattern: "my-?-blabla/abc*",
			text:    "someStringmnopqanda",
			matched: false,
		},
	}
	for i, testCase := range testCases {
		actualResult := matchWildcardAdvanced(testCase.pattern, testCase.text)
		assert.Equal(t, testCase.matched, actualResult, "Test %v failed: pattern=%v, text=%v", i+1, testCase.pattern, testCase.text)
	}
}

func TestMatchWildcardSimple(t *testing.T) {
	testCases := []struct {
		pattern string
		text    string
		matched bool
	}{
		{
			pattern: "*",
			text:    "s3:GetObject",
			matched: true,
		},
		{
			pattern: "",
			text:    "s3:GetObject",
			matched: false,
		},
		{
			pattern: "",
			text:    "",
			matched: true,
		},
		{
			pattern: "s3:*",
			text:    "s3:ListMultipartUploadParts",
			matched: true,
		},
		{
			pattern: "someRandomLongSentence",
			text:    "s3:Listblabla",
			matched: false,
		},
		{
			pattern: "s3:Listblabla",
			text:    "s3:Listblabla",
			matched: true,
		},
		{
			pattern: "someRandomLongSentence",
			text:    "someRandomLongSentence",
			matched: true,
		},
		{
			pattern: "someStringoo*",
			text:    "someStringoo",
			matched: true,
		},
		{
			pattern: "someStringIn*",
			text:    "someStringIndia/Karnataka/",
			matched: true,
		},
		{
			pattern: "someStringIn*",
			text:    "someStringKarnataka/India/",
			matched: false,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Ban/Ban/Ban/Ban/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Area1/Area2/Area3/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/State1/State2/Karnataka/Area1/Area2/Area3/Ban",
			matched: true,
		},
		{
			pattern: "someStringIn*/Ka*/Ban",
			text:    "someStringIndia/Karnataka/Bangalore",
			matched: false,
		},
		{
			pattern: "someStringIn*/Ka*/Ban*",
			text:    "someStringIndia/Karnataka/Bangalore",
			matched: true,
		},
		{
			pattern: "someString*",
			text:    "someStringIndia",
			matched: true,
		},
		{
			pattern: "someStringoo*",
			text:    "someStringodo",
			matched: false,
		},
		{
			pattern: "someStringoo?*",
			text:    "someStringoo???",
			matched: true,
		},
		{
			pattern: "someStringoo??*",
			text:    "someStringodo",
			matched: false,
		},
		{
			pattern: "?h?*",
			text:    "?h?hello",
			matched: true,
		},
	}
	for i, testCase := range testCases {
		actualResult := matchWildcardSimple(testCase.pattern, testCase.text)
		assert.Equal(t, testCase.matched, actualResult, "Test %v failed: pattern=%v, text=%v", i+1, testCase.pattern, testCase.text)
	}
}
