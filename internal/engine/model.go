package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
)

type Project struct {
	Name       string  `json:"name"`
	SampleRate int     `json:"sample_rate"`
	Channels   int     `json:"channels"`
	DurationMS int     `json:"duration_ms"`
	Tracks     []Track `json:"tracks"`
}

type Track struct {
	ID    string `json:"id"`
	Clips []Clip `json:"clips"`
}

type Clip struct {
	ID         string `json:"id"`
	StartMS    int    `json:"start_ms"`
	DurationMS int    `json:"duration_ms"`
	GainDB     int    `json:"gain_db"`
	FadeInMS   int    `json:"fade_in_ms"`
	FadeOutMS  int    `json:"fade_out_ms"`
	Source     Source `json:"source"`
}

type Source struct {
	Type        string `json:"type"`
	Path        string `json:"path,omitempty"`
	FrequencyHz int    `json:"frequency_hz,omitempty"`
}

type RenderReport struct {
	ProjectSHA256 string `json:"project_sha256"`
	PCMHash       string `json:"pcm_sha256"`
	Frames        int    `json:"frames"`
	SampleRate    int    `json:"sample_rate"`
	Channels      int    `json:"channels"`
	Peak          int16  `json:"peak"`
}

func LoadProject(path string) (Project, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Project{}, err
	}
	var p Project
	if err := json.Unmarshal(b, &p); err != nil {
		return Project{}, err
	}
	return p, Validate(p)
}

func Validate(p Project) error {
	if p.Name == "" {
		return errors.New("project name required")
	}
	if p.SampleRate < 8000 || p.SampleRate > 192000 {
		return fmt.Errorf("sample_rate out of range: %d", p.SampleRate)
	}
	if p.Channels != 1 && p.Channels != 2 {
		return errors.New("channels must be 1 or 2")
	}
	if p.DurationMS <= 0 || p.DurationMS > 3600000 {
		return errors.New("duration_ms must be 1..3600000")
	}
	seen := map[string]bool{}
	for _, t := range p.Tracks {
		if t.ID == "" || seen[t.ID] {
			return fmt.Errorf("invalid/duplicate track id %q", t.ID)
		}
		seen[t.ID] = true
		for _, c := range t.Clips {
			if c.ID == "" {
				return errors.New("clip id required")
			}
			if c.StartMS < 0 || c.DurationMS <= 0 || c.StartMS+c.DurationMS > p.DurationMS {
				return fmt.Errorf("clip %s is outside project bounds", c.ID)
			}
			if c.FadeInMS < 0 || c.FadeOutMS < 0 || c.FadeInMS+c.FadeOutMS > c.DurationMS {
				return fmt.Errorf("invalid fades for %s", c.ID)
			}
			switch c.Source.Type {
			case "tone":
				if c.Source.FrequencyHz < 20 || c.Source.FrequencyHz > 20000 {
					return fmt.Errorf("tone %s frequency out of range", c.ID)
				}
			default:
				return fmt.Errorf("unsupported source type %q", c.Source.Type)
			}
		}
	}
	return nil
}

func CanonicalJSON(p Project) ([]byte, error) {
	q := p
	sort.Slice(q.Tracks, func(i, j int) bool { return q.Tracks[i].ID < q.Tracks[j].ID })
	for i := range q.Tracks {
		sort.Slice(q.Tracks[i].Clips, func(a, b int) bool {
			if q.Tracks[i].Clips[a].StartMS == q.Tracks[i].Clips[b].StartMS {
				return q.Tracks[i].Clips[a].ID < q.Tracks[i].Clips[b].ID
			}
			return q.Tracks[i].Clips[a].StartMS < q.Tracks[i].Clips[b].StartMS
		})
	}
	return json.Marshal(q)
}

func Render(p Project) ([]int16, RenderReport, error) {
	if err := Validate(p); err != nil {
		return nil, RenderReport{}, err
	}
	frames := p.SampleRate * p.DurationMS / 1000
	mix := make([]int64, frames*p.Channels)
	for _, t := range p.Tracks {
		for _, c := range t.Clips {
			start := p.SampleRate * c.StartMS / 1000
			count := p.SampleRate * c.DurationMS / 1000
			gain := math.Pow(10, float64(c.GainDB)/20.0)
			for i := 0; i < count && start+i < frames; i++ {
				amp := gain
				ms := i * 1000 / p.SampleRate
				remaining := c.DurationMS - ms
				if c.FadeInMS > 0 && ms < c.FadeInMS {
					amp *= float64(ms) / float64(c.FadeInMS)
				}
				if c.FadeOutMS > 0 && remaining < c.FadeOutMS {
					amp *= math.Max(0, float64(remaining)/float64(c.FadeOutMS))
				}
				v := int64(math.Round(math.Sin(2*math.Pi*float64(c.Source.FrequencyHz)*float64(i)/float64(p.SampleRate)) * 12000 * amp))
				for ch := 0; ch < p.Channels; ch++ {
					mix[(start+i)*p.Channels+ch] += v
				}
			}
		}
	}
	pcm := make([]int16, len(mix))
	var peak int16
	for i, v := range mix {
		if v > 32767 {
			v = 32767
		}
		if v < -32768 {
			v = -32768
		}
		pcm[i] = int16(v)
		av := pcm[i]
		if av < 0 {
			if av == -32768 {
				av = 32767
			} else {
				av = -av
			}
		}
		if av > peak {
			peak = av
		}
	}
	canonical, _ := CanonicalJSON(p)
	ph := sha256.Sum256(canonical)
	pcmBytes := PCMBytes(pcm)
	mh := sha256.Sum256(pcmBytes)
	return pcm, RenderReport{ProjectSHA256: hex.EncodeToString(ph[:]), PCMHash: hex.EncodeToString(mh[:]), Frames: frames, SampleRate: p.SampleRate, Channels: p.Channels, Peak: peak}, nil
}

func PCMBytes(pcm []int16) []byte {
	b := make([]byte, len(pcm)*2)
	for i, v := range pcm {
		b[i*2] = byte(v)
		b[i*2+1] = byte(uint16(v) >> 8)
	}
	return b
}
