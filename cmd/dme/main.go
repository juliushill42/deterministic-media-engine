package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/juliushill42/deterministic-media-engine/internal/engine"
	"github.com/juliushill42/deterministic-media-engine/internal/server"
	"github.com/juliushill42/deterministic-media-engine/internal/wav"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "render":
		render(os.Args[2:])
	case "inspect":
		inspect(os.Args[2:])
	case "serve":
		serve(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}
func render(args []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	in := fs.String("in", "examples/demo.json", "project json")
	out := fs.String("out", "dist/render.wav", "wav output")
	_ = fs.Parse(args)
	p, err := engine.LoadProject(*in)
	must(err)
	pcm, report, err := engine.Render(p)
	must(err)
	f, err := os.Create(*out)
	must(err)
	defer f.Close()
	must(wav.WritePCM16(f, pcm, p.SampleRate, p.Channels))
	b, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(b))
}
func inspect(args []string) {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)
	in := fs.String("in", "examples/demo.json", "project json")
	_ = fs.Parse(args)
	p, err := engine.LoadProject(*in)
	must(err)
	_, report, err := engine.Render(p)
	must(err)
	b, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(b))
}
func serve(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", "127.0.0.1:8789", "listen address")
	_ = fs.Parse(args)
	must(server.Listen(*addr))
}
func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}
func usage() { fmt.Println("dme render|inspect|serve") }
