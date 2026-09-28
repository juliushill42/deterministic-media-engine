package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/juliushill42/deterministic-media-engine/internal/engine"
	"github.com/juliushill42/deterministic-media-engine/internal/ui"
	"github.com/juliushill42/deterministic-media-engine/internal/wav"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html; charset=utf-8")
		_, _ = w.Write(ui.Index)
	})
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"ok": true, "engine": "dme", "watermark": "\":\""})
	})
	mux.HandleFunc("POST /api/render", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
		var p engine.Project
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid project: "+err.Error(), 400)
			return
		}
		pcm, report, err := engine.Render(p)
		if err != nil {
			http.Error(w, err.Error(), 422)
			return
		}
		w.Header().Set("content-type", "audio/wav")
		w.Header().Set("x-pcm-sha256", report.PCMHash)
		w.Header().Set("x-render-frames", strconv.Itoa(report.Frames))
		if err := wav.WritePCM16(w, pcm, p.SampleRate, p.Channels); err != nil {
			return
		}
	})
	return http.TimeoutHandler(mux, 15*time.Second, "timeout")
}

func Listen(addr string) error {
	fmt.Println("DME listening on http://" + addr)
	return http.ListenAndServe(addr, Handler())
}
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
