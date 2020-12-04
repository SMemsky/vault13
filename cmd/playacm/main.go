package main

import (
	"os"
	"time"

	"github.com/dexter3k/vault13/file/acm"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func main() {
	f, err := os.Open(os.Args[1])
	check(err)
	defer f.Close()

	var s acm.SoundStreamer
	s.Start(f)

	sr := beep.SampleRate(s.SampleRate())
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(beep.StreamerFunc(func(out [][2]float64) (n int, ok bool) {
		for i := range out {
			l, r := s.NextSample()
			out[i][0], out[i][1] = float64(l), float64(r)
		}
		return len(out), true
	}))
	select {}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
