package merge_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/merge"
)

func makeChan(lines []merge.Line) <-chan merge.Line {
	ch := make(chan merge.Line, len(lines))
	for _, l := range lines {
		ch <- l
	}
	close(ch)
	return ch
}

func ts(sec int) time.Time {
	return time.Date(2024, 1, 1, 0, 0, sec, 0, time.UTC)
}

func TestMerge_SingleSource(t *testing.T) {
	lines := []merge.Line{
		{Text: "a", Timestamp: ts(1)},
		{Text: "b", Timestamp: ts(2)},
		{Text: "c", Timestamp: ts(3)},
	}
	m := merge.New([]<-chan merge.Line{makeChan(lines)})
	var got []merge.Line
	for l := range m.Merge() {
		got = append(got, l)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	for i, l := range got {
		if l.Text != lines[i].Text {
			t.Errorf("line %d: want %q, got %q", i, lines[i].Text, l.Text)
		}
	}
}

func TestMerge_TwoSources_Interleaved(t *testing.T) {
	src1 := []merge.Line{
		{Text: "s1-1", Timestamp: ts(1), Source: "s1"},
		{Text: "s1-3", Timestamp: ts(3), Source: "s1"},
		{Text: "s1-5", Timestamp: ts(5), Source: "s1"},
	}
	src2 := []merge.Line{
		{Text: "s2-2", Timestamp: ts(2), Source: "s2"},
		{Text: "s2-4", Timestamp: ts(4), Source: "s2"},
		{Text: "s2-6", Timestamp: ts(6), Source: "s2"},
	}
	m := merge.New([]<-chan merge.Line{makeChan(src1), makeChan(src2)})
	var got []merge.Line
	for l := range m.Merge() {
		got = append(got, l)
	}
	if len(got) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i].Timestamp.Before(got[i-1].Timestamp) {
			t.Errorf("out of order at index %d: %v before %v", i, got[i].Timestamp, got[i-1].Timestamp)
		}
	}
}

func TestMerge_EmptySource(t *testing.T) {
	empty := make(chan merge.Line)
	close(empty)
	src := []merge.Line{{Text: "only", Timestamp: ts(1)}}
	m := merge.New([]<-chan merge.Line{empty, makeChan(src)})
	var got []merge.Line
	for l := range m.Merge() {
		got = append(got, l)
	}
	if len(got) != 1 || got[0].Text != "only" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestMerge_NoSources(t *testing.T) {
	m := merge.New(nil)
	count := 0
	for range m.Merge() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 lines from empty merge, got %d", count)
	}
}

func TestMerge_SourceField(t *testing.T) {
	src := []merge.Line{
		{Text: fmt.Sprintf("line-%d", 1), Timestamp: ts(1), Source: "file.log"},
	}
	m := merge.New([]<-chan merge.Line{makeChan(src)})
	for l := range m.Merge() {
		if l.Source != "file.log" {
			t.Errorf("expected source 'file.log', got %q", l.Source)
		}
	}
}
