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

	ignore := map[uint16]bool{}
	for _, ignoreKey := range []string{
		"command",
		"right command",
		"shift",
		"right shift",
		"option",
		"right option",
		"escape",
		"fn",
		"tab",
		"control",
		"left arrow",
		"up arrow",
		"down arrow",
		"right arrow",
	} {
		ignore[keychartoRawCode(ignoreKey)] = true
	}

	kEnter := uint16(36)
	kBS := uint16(51)

	for e := range events {
		switch e.Kind {
		case hook.KeyDown:
			if ignore[e.Rawcode] {
				continue
			}
			switch e.Rawcode {
			case kEnter:
				speaker.Play(kbDing.streamer())
			case kBS:
				speaker.Play(kbBs.streamer())
			default:
				speaker.Play(kbDown.streamer())
			}
		case hook.KeyUp:
			// log.Println()
			// log.Println("up: raw", e.Rawcode, "keycode", e.Keycode, "mask", e.Mask, "keychar", e.Keychar, "toKeychar", rawcodetoKeychar(e.Rawcode))
			if ignore[e.Rawcode] {
				continue
			}
			switch e.Rawcode {
			case kEnter, kBS:
				// no sound
			default:
				speaker.Play(kbUp.streamer())
			}
		}
	}
}
