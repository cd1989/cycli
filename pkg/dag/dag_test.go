package dag

import (
	"strings"
	"testing"

	"github.com/kyokomi/emoji"
)

func TestAsciDAG(t *testing.T) {
	dag := AsciDAG{
		nodes: []*Node{
			{
				Name: "A",
			},
			{
				Name: "B",
			},
			{
				Name: "C",
			},
			{
				Name: "D",
			},
		},
		edges: []*Edge{
			{
				From: "A",
				To:   "B",
			},
			{
				From: "A",
				To:   "C",
			},
			{
				From: "B",
				To:   "D",
			},
			{
				From: "C",
				To:   "D",
			},
		},
		nodeElement: strings.TrimSpace(emoji.Sprint(":large_blue_circle:")),
	}

	dag.Render()
}
