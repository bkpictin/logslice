package grok

// config holds optional settings for pattern compilation.
type config struct {
	patterns map[string]string
}

// Option is a functional option for New.
type Option func(*config)

// WithPattern registers a custom named pattern that can be referenced inside
// templates as %{NAME} or %{NAME:field}. Custom patterns take precedence over
// built-ins with the same name.
func WithPattern(name, raw string) Option {
	return func(c *config) {
		c.patterns[name] = raw
	}
}

// WithPatterns registers multiple custom patterns at once.
func WithPatterns(m map[string]string) Option {
	return func(c *config) {
		for k, v := range m {
			c.patterns[k] = v
		}
	}
}
