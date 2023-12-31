package main

import (
	"log"
	"time"

	"github.com/faiface/beep/speaker"
	hook "github.com/robotn/gohook"
)

func main() {
	log.SetFlags(0)

	kbDown, err := newPalette("kd1.mp3", "kd2.mp3")
	if err != nil {
		log.Fatal(err)
	}
	kbUp, err := newPalette("ku1.mp3", "ku2.mp3")
	if err != nil {
		log.Fatal(err)
	}
	// kbSounds, err := newPaletteDownUp("kd1.mp3", "ku1.mp3", "kd2.mp3", "ku2.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	kbDing, err := newPalette("kding.mp3")
	if err != nil {
		log.Fatal(err)
	}
	kbBs, err := newPalette("kbs.mp3")
	if err != nil {
		log.Fatal(err)
	}
	format := kbDown.format
	// format := kbSounds.format

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/200))

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

	kReturn := keychartoRawCode("return")
	kDelete := keychartoRawCode("delete")

	last := time.Now()
	lastCode := uint16(0)
	lastKey := uint16(0)
	lastKind := uint8(0)
	ignoreThis := false

	for e := range events {
		switch e.Kind {
		case hook.KeyDown:
			since := time.Since(last)
			// fmt.Printf("\n%v\n", since)
			if e.Rawcode == lastCode {
				if since < 60*time.Millisecond || e.Kind == lastKind {
					ignoreThis = true
				}
			} else {
				ignoreThis = false
			}
			last = time.Now()
			lastKind = e.Kind

			lastCode = e.Rawcode
			if ignoreThis || ignore[e.Rawcode] {
				continue
			}
			switch e.Rawcode {
			case kReturn:
				if lastKey != kReturn {
					speaker.Play(kbDing.streamer())
				}
			case kDelete:
				speaker.Play(kbBs.streamer())
			default:
				// speaker.Play(kbSounds.streamer())
				speaker.Play(kbDown.streamer())
			}
			lastKey = e.Rawcode
		case hook.KeyUp:
			lastKind = e.Kind
			ignoreThis = false
			// log.Println()
			// log.Println("up: raw", e.Rawcode, "keycode", e.Keycode, "mask", e.Mask, "keychar", e.Keychar, "toKeychar", rawcodetoKeychar(e.Rawcode))
			if ignore[e.Rawcode] {
				continue
			}
			switch e.Rawcode {
			case kReturn, kDelete:
				// no sound
			default:
				speaker.Play(kbUp.streamer())
			}
		}
	}
}
