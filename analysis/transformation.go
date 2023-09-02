package analysis

type (
	Transformation struct{ Forward, Backward map[string]string }
)

func NewTransformation(forward map[string]string) Transformation {
	backward := make(map[string]string, len(forward))
	for a, b := range forward {
		backward[b] = a
	}
	return Transformation{Forward: forward, Backward: backward}
}
