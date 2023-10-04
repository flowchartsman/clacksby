package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

type palette struct {
	bufs   []*beep.Buffer
	format beep.Format
}

func newPalette(filenames ...string) (*palette, error) {
	bufs := make([]*beep.Buffer, 0, len(filenames))
	var sfmt beep.Format
	for i, name := range filenames {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		s, format, err := mp3.Decode(f)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			sfmt = format
		}
		if format.SampleRate != sfmt.SampleRate {
			return nil, fmt.Errorf("format mismatch, had samplerate: %d, got: %d", sfmt.SampleRate, format.SampleRate)
		}
		buf := beep.NewBuffer(format)
		buf.Append(s)
		if err := s.Close(); err != nil {
			return nil, fmt.Errorf("sample close: %v", err)
		}
		bufs = append(bufs, buf)
	}
	return &palette{
		bufs:   bufs,
		format: sfmt,
	}, nil
}

func (p *palette) streamer() beep.StreamSeeker {
	b := p.bufs[rand.Intn(len(p.bufs))]
	return b.Streamer(0, b.Len())
}
