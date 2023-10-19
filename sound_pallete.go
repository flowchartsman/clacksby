package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
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
		var (
			s      beep.StreamSeekCloser
			format beep.Format
		)

		ext := strings.ToLower(filepath.Ext(f.Name()))

		switch ext {
		case ".wav":
			s, format, err = wav.Decode(f)
		case ".mp3":
			s, format, err = mp3.Decode(f)
		default:
			return nil, fmt.Errorf("unknown/unsupported sample extension %s", ext)
		}

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

func newPaletteDownUp(filenames ...string) (*palette, error) {
	if len(filenames)%2 != 0 {
		return nil, fmt.Errorf("must be an even number of down-up pairs")
	}
	bufs := make([]*beep.Buffer, 0, len(filenames)/2)
	var sfmt beep.Format
	for i := 0; i < len(filenames); i += 2 {
		down := filenames[i]
		up := filenames[i+1]
		var buf *beep.Buffer
		for _, name := range []string{down, up} {
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
			if buf == nil {
				buf = beep.NewBuffer(format)
			}
			buf.Append(s)
			if err := s.Close(); err != nil {
				return nil, fmt.Errorf("sample close: %v", err)
			}
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
