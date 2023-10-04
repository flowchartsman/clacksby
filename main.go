package main

import (
	"log"
	"time"

	"github.com/faiface/beep/speaker"
	hook "github.com/robotn/gohook"
)

func main() {
	log.SetFlags(0)

	kbDown, err := newPalette("kd1.mp3", "kd2.mp3", "kd3.mp3")
	if err != nil {
		log.Fatal(err)
	}
	kbUp, err := newPalette("ku2.mp3", "ku3.mp3")
	if err != nil {
		log.Fatal(err)
	}
	kbDing, err := newPalette("kding.mp3")
	if err != nil {
		log.Fatal(err)
	}
	kbBs, err := newPalette("kbs.mp3")
	if err != nil {
		log.Fatal(err)
	}
	format := kbDown.format

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/100))

	l := newListener()
	events := l.start()
	defer l.stop()

	kEnter := uint16(36)
	kBS := uint16(51)

	for e := range events {
		switch e.Kind {
		case hook.KeyDown:
			// log.Println("down: raw", e.Rawcode, "keycode", e.Keycode, "mask", e.Mask, "keychar", e.Keychar)
			switch e.Rawcode {
			case kEnter:
				speaker.Play(kbDing.streamer())
			case kBS:
				speaker.Play(kbBs.streamer())
			default:
				speaker.Play(kbDown.streamer())
			}
		case hook.KeyUp:
			// log.Println("up: raw", e.Rawcode, "keycode", e.Keycode, "mask", e.Mask, "keychar", e.Keychar)
			switch e.Rawcode {
			case kEnter, kBS:
				//
			default:
				speaker.Play(kbUp.streamer())
			}
		}
	}
}
