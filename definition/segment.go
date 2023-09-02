package definition

import (
	"slices"
)

type (
	Segment  struct{ Start, End int }
	Segments struct {
		segments    []Segment
		totalLength []int
	}
)

func (s Segment) Length() int { return s.End - s.Start }

func BuildSegments(segments ...Segment) Segments {
	totalLength := make([]int, 0, len(segments))
	length := 0
	for _, segment := range segments {
		length += segment.End - segment.Start
		totalLength = append(totalLength, length)
	}
	return Segments{segments: segments, totalLength: totalLength}
}

func JoinSegments(segments ...Segments) Segments {
	all := make([]Segment, 0)
	for _, segment := range segments {
		all = append(all, segment.segments...)
	}
	return BuildSegments(all...)
}

func (s Segments) Cut(start, end int) Segments {
	if len(s.segments) == 0 {
		return Segments{}
	}
	if len(s.segments) == 1 {
		return BuildSegments(Segment{Start: s.segments[0].Start + start, End: s.segments[0].End})
	}

	lower, _ := slices.BinarySearchFunc(s.totalLength, start, func(end int, pos int) int { return end - (pos + 1) })
	upper, _ := slices.BinarySearchFunc(s.totalLength, end, func(end int, pos int) int { return end - pos })
	segments := make([]Segment, 0, upper-lower+1)
	for i := lower; i <= upper; i++ {
		segment := s.segments[i]
		previous := 0
		if i > 0 {
			previous = s.totalLength[i-1]
		}
		segments = append(segments, Segment{
			Start: segment.Start + max(0, start-previous),
			End:   segment.End - max(0, s.totalLength[i]-end),
		})
	}
	return BuildSegments(segments...)
}
