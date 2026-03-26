package llama

func WithGPULayers(n int) ModelOption {
	return func(c *modelConfig) {
		c.gpuLayers = n
	}
}

func WithMLock() ModelOption {
	return func(c *modelConfig) {
		c.mlock = true
	}
}

func WithMMap(enabled bool) ModelOption {
	return func(c *modelConfig) {
		c.mmap = enabled
	}
}

func WithMainGPU(gpu string) ModelOption {
	return func(c *modelConfig) {
		c.mainGPU = gpu
	}
}

func WithTensorSplit(split string) ModelOption {
	return func(c *modelConfig) {
		c.tensorSplit = split
	}
}

func WithSilentLoading() ModelOption {
	return func(c *modelConfig) {
		c.disableProgressCallback = true
	}
}

type ProgressCallback func(progress float32) bool

func WithProgressCallback(cb ProgressCallback) ModelOption {
	return func(c *modelConfig) {
		c.progressCallback = cb
	}
}
