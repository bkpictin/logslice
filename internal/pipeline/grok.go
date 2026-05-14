package pipeline

import (
	"github.com/user/logslice/internal/grok"
)

// GrokConfig holds the configuration for grok-based field extraction within
// the pipeline.  When Pattern is non-empty the pipeline will attempt to match
// every log line and attach the extracted fields to the entry metadata.
type GrokConfig struct {
	// Pattern is a grok template string, e.g.
	// "%{IP:client} %{WORD:method} %{NOTSPACE:path} %{INT:status}"
	Pattern string

	// ExtraPatterns allows callers to register additional named patterns that
	// can be referenced inside Pattern.
	ExtraPatterns map[string]string
}

// grokExtractor wraps a compiled grok.Pattern and applies it to log lines.
type grokExtractor struct {
	p *grok.Pattern
}

// newGrokExtractor compiles the grok pattern from cfg.  Returns nil and no
// error when cfg.Pattern is empty (feature disabled).
func newGrokExtractor(cfg GrokConfig) (*grokExtractor, error) {
	if cfg.Pattern == "" {
		return nil, nil
	}
	opts := []grok.Option{}
	if len(cfg.ExtraPatterns) > 0 {
		opts = append(opts, grok.WithPatterns(cfg.ExtraPatterns))
	}
	p, err := grok.New(cfg.Pattern, opts...)
	if err != nil {
		return nil, err
	}
	return &grokExtractor{p: p}, nil
}

// applyGrok runs the compiled pattern against line and returns the extracted
// fields.  If the extractor is nil or the line does not match, nil is returned
// so callers can decide whether to drop or pass through the line.
func applyGrok(g *grokExtractor, line string) map[string]string {
	if g == nil {
		return nil
	}
	return g.p.Match(line)
}
