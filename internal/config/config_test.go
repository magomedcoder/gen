package config

import (
	"testing"

	"github.com/magomedcoder/gen/internal/domain"
)

func TestRAGConfigEffectiveHydeDefaultsAndBounds(t *testing.T) {
	var c RAGConfig
	if !c.EffectiveHydeEnabled() {
		t.Fatal("default hyde_enabled should be true (nil)")
	}

	if c.EffectiveHydeMaxTokens() != ragDefaultHydeMaxTokens {
		t.Fatalf("default hyde_max_tokens: %d", c.EffectiveHydeMaxTokens())
	}

	if c.EffectiveHydeTimeoutSeconds() != ragDefaultHydeTimeoutSec {
		t.Fatalf("default hyde_timeout_seconds: %d", c.EffectiveHydeTimeoutSeconds())
	}

	f := false
	c.HydeEnabled = &f
	if c.EffectiveHydeEnabled() {
		t.Fatal("hyde_enabled false")
	}

	c.HydeEnabled = nil
	c.HydeMaxTokens = 4
	if c.EffectiveHydeMaxTokens() != 32 {
		t.Fatalf("min hyde_max_tokens clamp: %d", c.EffectiveHydeMaxTokens())
	}

	c.HydeMaxTokens = 9999
	if c.EffectiveHydeMaxTokens() != 768 {
		t.Fatalf("max hyde_max_tokens clamp: %d", c.EffectiveHydeMaxTokens())
	}

	c.HydeTimeoutSeconds = 9999
	if c.EffectiveHydeTimeoutSeconds() != 300 {
		t.Fatalf("max hyde_timeout_seconds clamp: %d", c.EffectiveHydeTimeoutSeconds())
	}
}

func TestRAGConfigEffectiveRerankDefaults(t *testing.T) {
	var c RAGConfig
	if !c.EffectiveRerankEnabled() {
		t.Fatal("default rerank_enabled should be true (nil)")
	}

	if c.EffectiveRerankMaxCandidates() != ragDefaultRerankMaxCandidates {
		t.Fatalf("rerank_max_candidates: %d", c.EffectiveRerankMaxCandidates())
	}

	if c.EffectiveRerankMaxTokens() != ragDefaultRerankMaxTokens {
		t.Fatalf("rerank_max_tokens: %d", c.EffectiveRerankMaxTokens())
	}

	if c.EffectiveRerankTimeoutSeconds() != ragDefaultRerankTimeoutSec {
		t.Fatalf("rerank_timeout: %d", c.EffectiveRerankTimeoutSeconds())
	}

	if c.EffectiveRerankPassageMaxRunes() != ragDefaultRerankPassageMaxRunes {
		t.Fatalf("rerank_passage_max_runes: %d", c.EffectiveRerankPassageMaxRunes())
	}

	c.RerankMaxCandidates = 99
	if c.EffectiveRerankMaxCandidates() != 32 {
		t.Fatalf("rerank_max_candidates clamp: %d", c.EffectiveRerankMaxCandidates())
	}
}

func TestRAGConfigEffectiveDeepRAGDefaults(t *testing.T) {
	var c RAGConfig
	if !c.EffectiveDeepRAGEnabled() {
		t.Fatal("default deep_rag_enabled should be true (nil)")
	}

	if c.EffectiveDeepRAGMaxMapCalls() != ragDefaultDeepRAGMaxMapCalls {
		t.Fatalf("deep_rag_max_map_calls default: %d", c.EffectiveDeepRAGMaxMapCalls())
	}

	if c.EffectiveDeepRAGChunksPerMap() != ragDefaultDeepRAGChunksPerMap {
		t.Fatalf("deep_rag_chunks_per_map default: %d", c.EffectiveDeepRAGChunksPerMap())
	}

	if c.EffectiveDeepRAGMapMaxTokens() != ragDefaultDeepRAGMapMaxTokens {
		t.Fatalf("deep_rag_map_max_tokens default: %d", c.EffectiveDeepRAGMapMaxTokens())
	}

	if c.EffectiveDeepRAGMapTimeoutSeconds() != ragDefaultDeepRAGMapTimeoutSec {
		t.Fatalf("deep_rag_map_timeout default: %d", c.EffectiveDeepRAGMapTimeoutSeconds())
	}

	if c.EffectiveDeepRAGMaxMapOutputRunes() != ragDefaultDeepRAGMaxMapOutputRunes {
		t.Fatalf("deep_rag_max_map_output_runes default: %d", c.EffectiveDeepRAGMaxMapOutputRunes())
	}

	c.DeepRAGMaxMapCalls = 99
	if c.EffectiveDeepRAGMaxMapCalls() != 16 {
		t.Fatalf("deep_rag_max_map_calls clamp: %d", c.EffectiveDeepRAGMaxMapCalls())
	}
}

func TestRAGConfigEffectiveAdaptiveKAndMinScore(t *testing.T) {
	var c RAGConfig
	if !c.EffectiveAdaptiveKEnabled() {
		t.Fatal("default adaptive_k_enabled should be true (nil)")
	}

	if c.EffectiveAdaptiveKMultiplier() != ragDefaultAdaptiveKMultiplier {
		t.Fatalf("default multiplier: %d", c.EffectiveAdaptiveKMultiplier())
	}

	if c.EffectiveMinSimilarityScore() != -1 {
		t.Fatalf("default min_similarity_score: %v", c.EffectiveMinSimilarityScore())
	}

	c.AdaptiveKMultiplier = 99
	if c.EffectiveAdaptiveKMultiplier() != 6 {
		t.Fatalf("max multiplier clamp: %d", c.EffectiveAdaptiveKMultiplier())
	}

	c.MinSimilarityScore = 5
	if c.EffectiveMinSimilarityScore() != 1 {
		t.Fatalf("max min_similarity_score clamp: %v", c.EffectiveMinSimilarityScore())
	}
}

func TestValidateMCPServerHTTPIgnoredForNonHTTPTransport(t *testing.T) {
	cfg := &Config{
		MCP: MCPConfig{},
	}
	s := &domain.MCPServer{
		Transport: "invalid",
	}
	if err := cfg.ValidateMCPServerHTTP(s); err != nil {
		t.Fatalf("non-http транспорт не должен валидироваться как HTTP: %v", err)
	}
}
