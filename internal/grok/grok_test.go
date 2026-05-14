package grok

import (
	"testing"
)

func TestNew_SimpleExpansion(t *testing.T) {
	p, err := New(`%{IP:client} - %{WORD:method}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := p.Match("10.0.0.1 - GET")
	if fields == nil {
		t.Fatal("expected match, got nil")
	}
	if fields["client"] != "10.0.0.1" {
		t.Errorf("client: got %q", fields["client"])
	}
	if fields["method"] != "GET" {
		t.Errorf("method: got %q", fields["method"])
	}
}

func TestNew_NoMatch(t *testing.T) {
	p, err := New(`%{IP:client}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Match("not an ip address here") != nil {
		t.Error("expected nil for non-matching line")
	}
}

func TestNew_UnknownPattern(t *testing.T) {
	_, err := New(`%{UNKNOWN_PATTERN:field}`)
	if err == nil {
		t.Fatal("expected error for unknown pattern")
	}
}

func TestNew_InvalidRegexp(t *testing.T) {
	_, err := New(`[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestWithPattern_CustomPattern(t *testing.T) {
	p, err := New(`%{MYTOKEN:tok}`, WithPattern("MYTOKEN", `[A-Z]{3}\d{3}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := p.Match("prefix ABC123 suffix")
	if fields == nil {
		t.Fatal("expected match")
	}
	if fields["tok"] != "ABC123" {
		t.Errorf("tok: got %q", fields["tok"])
	}
}

func TestWithPatterns_MultipleCustom(t *testing.T) {
	p, err := New(`%{A:a} %{B:b}`, WithPatterns(map[string]string{
		"A": `foo`,
		"B": `bar`,
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := p.Match("foo bar")
	if fields["a"] != "foo" || fields["b"] != "bar" {
		t.Errorf("got %v", fields)
	}
}

func TestFields_ReturnsCaptureNames(t *testing.T) {
	p, err := New(`%{IP:host} %{INT:port}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := p.Fields()
	if len(names) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(names))
	}
}

func TestExpand_DepthLimit(t *testing.T) {
	// Build a chain of patterns that reference each other to trigger depth limit.
	opts := make([]Option, 0)
	prev := `x`
	for i := 0; i < 20; i++ {
		curr := prev
		name := string(rune('A' + i))
		opts = append(opts, WithPattern(name, `%{`+prev+`}`))
		_ = curr
		prev = name
	}
	_, err := New(`%{`+prev+`:f}`, opts...)
	if err == nil {
		t.Fatal("expected depth-limit error")
	}
}

func TestMatch_LogLevel(t *testing.T) {
	p, err := New(`%{LOGLEVEL:level}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, lvl := range []string{"INFO", "WARN", "ERROR", "DEBUG", "FATAL"} {
		fields := p.Match("2024-01-01 " + lvl + " message")
		if fields == nil || fields["level"] != lvl {
			t.Errorf("level %q: got %v", lvl, fields)
		}
	}
}
