package parser

//
//import (
//	"github.com/stretchr/testify/assert"
//	"gopeg/definition"
//	"testing"
//)
//
//type PegTest struct {
//	R []struct {
//		PegBindRef
//		A *PegBindRef `peg:"A"`
//		B *PegBindRef `peg:"B"`
//	} `peg:"R"`
//}
//
//func TestBind(t *testing.T) {
//	var a PegTest
//	err := Bind(&definition.parsingNode{
//		start: 0,
//		end:   2,
//		children: []definition.ParsingNode{&definition.parsingNode{
//			symbol: "R",
//			start:  0,
//			end:    1,
//			children: []definition.ParsingNode{
//				&definition.parsingNode{
//					symbol:   "A",
//					start:    0,
//					end:      1,
//					children: nil,
//				},
//			},
//		}, &definition.parsingNode{
//			symbol: "R",
//			start:  1,
//			end:    2,
//			children: []definition.ParsingNode{
//				&definition.parsingNode{
//					symbol:   "B",
//					start:    1,
//					end:      2,
//					children: nil,
//				},
//			},
//		}},
//	}, &a)
//	assert.Nil(t, err)
//
//	assert.Len(t, a.R, 2)
//	assert.Equal(t, a.R[0].Start, 0)
//	assert.Equal(t, a.R[0].End, 1)
//	assert.Equal(t, a.R[0].A, &PegBindRef{0, 1})
//	assert.Nil(t, a.R[0].B)
//
//	assert.Equal(t, a.R[1].Start, 1)
//	assert.Equal(t, a.R[1].End, 2)
//	assert.Nil(t, a.R[1].A)
//	assert.Equal(t, a.R[1].B, &PegBindRef{1, 2})
//}
