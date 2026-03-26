package llama

import (
	"runtime"
)

func WithContext(size int) ContextOption {
	return func(c *contextConfig) {
		c.contextSize = size
	}
}

func WithBatch(size int) ContextOption {
	return func(c *contextConfig) {
		c.batchSize = size
	}
}

func WithThreads(n int) ContextOption {
	return func(c *contextConfig) {
		c.threads = n
	}
}

func WithThreadsBatch(n int) ContextOption {
	return func(c *contextConfig) {
		c.threadsBatch = n
	}
}

func WithF16Memory() ContextOption {
	return func(c *contextConfig) {
		c.f16Memory = true
	}
}

func WithEmbeddings() ContextOption {
	return func(c *contextConfig) {
		c.embeddings = true
	}
}

func WithKVCacheType(cacheType string) ContextOption {
	return func(c *contextConfig) {
		switch cacheType {
		case
			"f16",
			"q8_0",
			"q4_0":
			c.kvCacheType = cacheType
		default:
		}
	}
}

func WithFlashAttn(mode string) ContextOption {
	return func(c *contextConfig) {
		switch mode {
		case
			"auto",
			"enabled",
			"disabled":
			c.flashAttn = mode
		default:
		}
	}
}

func WithParallel(n int) ContextOption {
	return func(c *contextConfig) {
		if n < 1 {
			n = 1
		}

		c.nParallel = n
	}
}

func WithPrefixCaching(enabled bool) ContextOption {
	return func(c *contextConfig) {
		c.prefixCaching = enabled
	}
}

func init() {
	defaultContextConfig.threads = runtime.NumCPU()
}
