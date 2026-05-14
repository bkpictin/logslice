package grok

import (
	"fmt"
	"regexp"
	"strings"
)

// Predefined named patterns that can be referenced inside templates.
var builtins = map[string]string{
	"INT":        `[+-]?(?:[0-9]+)`,
	"NUMBER":     `[+-]?(?:[0-9]+\.?[0-9]*)`,
	"WORD":       `\b\w+\b`,
	"NOTSPACE":   `\S+`,
	"SPACE":      `\s+`,
	"DATA":       `.*?`,
	"GREEDYDATA": `.*`,
	"IP":         `(?:[0-9]{1,3}\.){3}[0-9]{1,3}`,
	"TIMESTAMP":  `%{YEAR}-%{MONTHNUM}-%{MONTHDAY}[T ]%{HOUR}:%{MINUTE}:%{SECOND}`,
	"YEAR":       `[0-9]{4}`,
	"MONTHNUM":   `0?[1-9]|1[0-2]`,
	"MONTHDAY":   `(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9]`,
	"HOUR":       `2[0123]|[01]?[0-9]`,
	"MINUTE":     `[0-5][0-9]`,
	"SECOND":     `(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?`,
	"LOGLEVEL":   `[Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn(?:ing)?|WARN(?:ING)?|[Ee]rr(?:or)?|ERR(?:OR)?|[Cc]rit(?:ical)?|CRIT(?:ICAL)?|[Ff]atal|FATAL|[Ss]evere|SEVERE`,
}

// Pattern compiles a grok-style pattern string into a regexp and exposes
// named capture groups as a field map.
type Pattern struct {
	re     *regexp.Regexp
	fields []string
}

// New compiles the grok template, expanding %{NAME} and %{NAME:field}
// references recursively using the built-in dictionary plus any extras
// provided via WithPattern options.
func New(template string, opts ...Option) (*Pattern, error) {
	cfg := &config{patterns: make(map[string]string)}
	for _, o := range opts {
		o(cfg)
	}

	resolved, err := expand(template, cfg.patterns, 0)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(resolved)
	if err != nil {
		return nil, fmt.Errorf("grok: compile %w", err)
	}

	p := &Pattern{re: re}
	for _, n := range re.SubexpNames() {
		if n != "" {
			p.fields = append(p.fields, n)
		}
	}
	return p, nil
}

// Match returns a map of named captures for line, or nil if no match.
func (p *Pattern) Match(line string) map[string]string {
	m := p.re.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(p.fields))
	for i, name := range p.re.SubexpNames() {
		if name != "" && i < len(m) {
			out[name] = m[i]
		}
	}
	return out
}

// Fields returns the named capture group names in pattern order.
func (p *Pattern) Fields() []string { return p.fields }

// expand resolves %{NAME} and %{NAME:field} tokens up to maxDepth levels deep.
func expand(s string, extra map[string]string, depth int) (string, error) {
	if depth > 16 {
		return "", fmt.Errorf("grok: pattern expansion depth exceeded")
	}
	re := regexp.MustCompile(`%\{(\w+)(?::(\w+))?\}`)
	var expandErr error
	result := re.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return ""
		}
		subs := re.FindStringSubmatch(match)
		name, field := subs[1], subs[2]
		raw, ok := extra[name]
		if !ok {
			raw, ok = builtins[name]
		}
		if !ok {
			expandErr = fmt.Errorf("grok: unknown pattern %q", name)
			return ""
		}
		inner, err := expand(raw, extra, depth+1)
		if err != nil {
			expandErr = err
			return ""
		}
		if field != "" {
			return fmt.Sprintf("(?P<%s>%s)", field, inner)
		}
		return fmt.Sprintf("(?:%s)", inner)
	})
	if expandErr != nil {
		return "", expandErr
	}
	_ = strings.Contains // keep import
	return result, nil
}
