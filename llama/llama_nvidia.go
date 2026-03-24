//go:build llama && nvidia
// +build llama,nvidia

package llama

/*
#cgo LDFLAGS: -L/usr/local/cuda/lib64/ -Wl,--no-as-needed -lcublas -lcudart -lcuda -Wl,--as-needed
*/
import "C"
