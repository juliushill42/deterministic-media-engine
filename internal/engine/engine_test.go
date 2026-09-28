package engine

import "testing"

func TestDeterministicRender(t *testing.T) {
	p := Project{Name: "x", SampleRate: 48000, Channels: 2, DurationMS: 1000, Tracks: []Track{{ID: "a", Clips: []Clip{{ID: "c", StartMS: 0, DurationMS: 1000, GainDB: -3, FadeInMS: 10, FadeOutMS: 20, Source: Source{Type: "tone", FrequencyHz: 440}}}}}}
	a, ra, e := Render(p)
	if e != nil {
		t.Fatal(e)
	}
	b, rb, e := Render(p)
	if e != nil {
		t.Fatal(e)
	}
	if ra.PCMHash != rb.PCMHash {
		t.Fatal("hash mismatch")
	}
	if len(a) != len(b) || len(a) == 0 {
		t.Fatal("bad pcm")
	}
}
func TestRejectBounds(t *testing.T) {
	p := Project{Name: "x", SampleRate: 48000, Channels: 2, DurationMS: 100, Tracks: []Track{{ID: "a", Clips: []Clip{{ID: "c", StartMS: 90, DurationMS: 20, Source: Source{Type: "tone", FrequencyHz: 440}}}}}}
	if Validate(p) == nil {
		t.Fatal("expected bounds error")
	}
}
