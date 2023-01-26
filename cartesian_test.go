package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartesian(t *testing.T) {

	in1 := []string{"a", "b", "c"}
	in2 := []string{"1", "2", "3"}

	actual := cartesian(in1, in2)

	expected := [][]string{
		{
			"a", "1",
		},
		{
			"a", "2",
		},
		{
			"a", "3",
		},
		{
			"b", "1",
		},
		{
			"b", "2",
		},
		{
			"b", "3",
		},
		{
			"c", "1",
		},
		{
			"c", "2",
		},
		{
			"c", "3",
		},
	}

	assert.ElementsMatch(t, expected, actual)

}
