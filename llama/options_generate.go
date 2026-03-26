package llama

func WithMaxTokens(n int) GenerateOption {
	return func(c *generateConfig) {
		c.maxTokens = n
	}
}

func WithTemperature(t float32) GenerateOption {
	return func(c *generateConfig) {
		c.temperature = t
	}
}

func WithTopP(p float32) GenerateOption {
	return func(c *generateConfig) {
		c.topP = p
	}
}

func WithTopK(k int) GenerateOption {
	return func(c *generateConfig) {
		c.topK = k
	}
}

func WithSeed(seed int) GenerateOption {
	return func(c *generateConfig) {
		c.seed = seed
	}
}

func WithStopWords(words ...string) GenerateOption {
	return func(c *generateConfig) {
		c.stopWords = words
	}
}

func WithDraftTokens(n int) GenerateOption {
	return func(c *generateConfig) {
		c.draftTokens = n
	}
}

func WithDebug() GenerateOption {
	return func(c *generateConfig) {
		c.debug = true
	}
}

func WithMinP(p float32) GenerateOption {
	return func(c *generateConfig) {
		c.minP = p
	}
}

func WithTypicalP(p float32) GenerateOption {
	return func(c *generateConfig) {
		c.typP = p
	}
}

func WithTopNSigma(sigma float32) GenerateOption {
	return func(c *generateConfig) {
		c.topNSigma = sigma
	}
}

func WithMinKeep(n int) GenerateOption {
	return func(c *generateConfig) {
		c.minKeep = n
	}
}

func WithRepeatPenalty(penalty float32) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyRepeat = penalty
	}
}

func WithFrequencyPenalty(penalty float32) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyFreq = penalty
	}
}

func WithPresencePenalty(penalty float32) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyPresent = penalty
	}
}

func WithPenaltyLastN(n int) GenerateOption {
	return func(c *generateConfig) {
		c.penaltyLastN = n
	}
}

func WithDRYMultiplier(mult float32) GenerateOption {
	return func(c *generateConfig) {
		c.dryMultiplier = mult
	}
}

func WithDRYBase(base float32) GenerateOption {
	return func(c *generateConfig) {
		c.dryBase = base
	}
}

func WithDRYAllowedLength(length int) GenerateOption {
	return func(c *generateConfig) {
		c.dryAllowedLength = length
	}
}

func WithDRYPenaltyLastN(n int) GenerateOption {
	return func(c *generateConfig) {
		c.dryPenaltyLastN = n
	}
}

func WithDRYSequenceBreakers(breakers ...string) GenerateOption {
	return func(c *generateConfig) {
		c.drySequenceBreakers = breakers
	}
}

func WithDynamicTemperature(tempRange, exponent float32) GenerateOption {
	return func(c *generateConfig) {
		c.dynatempRange = tempRange
		c.dynatempExponent = exponent
	}
}

func WithXTC(probability, threshold float32) GenerateOption {
	return func(c *generateConfig) {
		c.xtcProbability = probability
		c.xtcThreshold = threshold
	}
}

func WithMirostat(version int) GenerateOption {
	return func(c *generateConfig) {
		c.mirostat = version
	}
}

func WithMirostatTau(tau float32) GenerateOption {
	return func(c *generateConfig) {
		c.mirostatTau = tau
	}
}

func WithMirostatEta(eta float32) GenerateOption {
	return func(c *generateConfig) {
		c.mirostatEta = eta
	}
}

func WithNPrev(n int) GenerateOption {
	return func(c *generateConfig) {
		c.nPrev = n
	}
}

func WithNProbs(n int) GenerateOption {
	return func(c *generateConfig) {
		c.nProbs = n
	}
}

func WithIgnoreEOS(ignore bool) GenerateOption {
	return func(c *generateConfig) {
		c.ignoreEOS = ignore
	}
}
