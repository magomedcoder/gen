//go:build llama
// +build llama

package llama

type ModelOptions struct {
	ContextSize   int
	Seed          int
	NBatch        int
	F16Memory     bool
	MLock         bool
	MMap          bool
	LowVRAM       bool
	Embeddings    bool
	NUMA          bool
	NGPULayers    int
	MainGPU       string
	TensorSplit   string
	FreqRopeBase  float32
	FreqRopeScale float32
	LoraBase      string
	LoraAdapter   string
}

type ModelOption func(p *ModelOptions)
