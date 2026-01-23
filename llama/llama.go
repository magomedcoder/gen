//go:build llama
// +build llama

package llama

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR} -lllama -lm -lstdc++
#include "llama.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

var errLoadModel = errors.New("не удалось загрузить модель")

type LLama struct {
	state unsafe.Pointer
}

func New(model string, opts ...ModelOption) (*LLama, error) {
	mo := NewModelOptions(opts...)
	modelPath := C.CString(model)
	defer C.free(unsafe.Pointer(modelPath))

	loraBase := C.CString(mo.LoraBase)

	defer C.free(unsafe.Pointer(loraBase))

	loraAdapter := C.CString(mo.LoraAdapter)

	defer C.free(unsafe.Pointer(loraAdapter))

	result := C.load_model(
		modelPath,
		C.int(mo.ContextSize),
		C.int(mo.Seed),
		C.bool(mo.F16Memory),
		C.bool(mo.MLock),
		C.bool(mo.Embeddings),
		C.bool(mo.MMap),
		C.bool(mo.LowVRAM),
		C.int(mo.NGPULayers),
		C.int(mo.NBatch),
		C.CString(mo.MainGPU),
		C.CString(mo.TensorSplit),
		C.bool(mo.NUMA),
		C.float(mo.FreqRopeBase),
		C.float(mo.FreqRopeScale),
		loraAdapter, loraBase,
	)
	if result == nil {
		return nil, errLoadModel
	}

	return &LLama{
		state: result,
	}, nil
}

func (l *LLama) Free() {
	C.llama_binding_free_model(l.state)
}
