package definition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildSegment(t *testing.T) {
	r := BuildSegments([]Segment{{Start: 0, End: 5}, {Start: 10, End: 12}, {Start: 100, End: 107}}...)
	require.Equal(t, []int{5, 7, 14}, r.totalLength)
}

func TestCutSegment(t *testing.T) {
	r := BuildSegments([]Segment{{Start: 0, End: 5}, {Start: 10, End: 12}, {Start: 100, End: 107}}...)
	require.Equal(t, []Segment{{Start: 0, End: 5}}, r.Cut(0, 5).segments)
	require.Equal(t, []Segment{{Start: 10, End: 12}}, r.Cut(5, 7).segments)
	require.Equal(t, []Segment{{Start: 10, End: 12}, {Start: 100, End: 101}}, r.Cut(5, 8).segments)
	require.Equal(t, []Segment{{Start: 10, End: 12}, {Start: 100, End: 102}}, r.Cut(5, 9).segments)
	require.Equal(t, []Segment{{Start: 100, End: 102}}, r.Cut(7, 9).segments)
	require.Equal(t, []Segment{{Start: 101, End: 102}}, r.Cut(8, 9).segments)
	require.Equal(t, []Segment{{Start: 4, End: 5}, {Start: 10, End: 12}, {Start: 100, End: 101}}, r.Cut(4, 8).segments)
}
