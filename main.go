package main

import (
	"log"
	"time"

	"github.com/faiface/beep/speaker"
	hook "github.com/robotn/gohook"
)

func main() {
	log.SetFlags(0)

	// smpDown, err := newPalette(
	// 	"fallout/01.wav",
	// 	"fallout/02.wav",
	// 	"fallout/03.wav",
	// 	"fallout/04.wav",
	// 	"fallout/05.wav",
	// 	"fallout/06.wav",
	// )
	// smpEnter, err := newPalette(
	// 	"fallout/charenter_01.wav",
	// 	"fallout/charenter_02.wav",
	// 	"fallout/charenter_03.wav",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// smpDelete, err := newPalette("fallout/charenter_01.wav")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// smpDown, err := newPalette(
	// 	"nk-cream/A.mp3",
	// 	"nk-cream/B.mp3",
	// 	"nk-cream/C.mp3",
	// 	"nk-cream/D.mp3",
	// 	"nk-cream/E.mp3",
	// 	"nk-cream/F.mp3",
	// 	"nk-cream/G.mp3",
	// 	"nk-cream/H.mp3",
	// 	"nk-cream/I.mp3",
	// 	"nk-cream/J.mp3",
	// 	"nk-cream/K.mp3",
	// 	"nk-cream/L.mp3",
	// 	"nk-cream/M.mp3",
	// 	"nk-cream/N.mp3",
	// 	"nk-cream/O.mp3",
	// 	"nk-cream/P.mp3",
	// 	"nk-cream/Q.mp3",
	// 	"nk-cream/R.mp3",
	// 	"nk-cream/S.mp3",
	// 	"nk-cream/T.mp3",
	// 	"nk-cream/U.mp3",
	// 	"nk-cream/V.mp3",
	// 	"nk-cream/W.mp3",
	// 	"nk-cream/X.mp3",
	// 	"nk-cream/Y.mp3",
	// 	"nk-cream/Z.mp3",
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// smpDelete, err := newPalette("nk-cream/BACKSPACE.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// smpEnter, err := newPalette("nk-cream/ENTER.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	smpDown, err := newPalette("typewriter/kd1.mp3", "typewriter/kd2.mp3")
	if err != nil {
		log.Fatal(err)
	}
	smpUp, err := newPalette("typewriter/ku1.mp3", "typewriter/ku2.mp3")
	if err != nil {
		log.Fatal(err)
	}
	// smp, err := newPaletteDownUp("kd1.mp3", "ku1.mp3", "kd2.mp3", "ku2.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	smpEnter, err := newPalette("typewriter/kcr.mp3")
	if err != nil {
		log.Fatal(err)
	}
	smpDelete, err := newPalette("typewriter/kbs.mp3")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

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

	format := smpDown.format
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/200))

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
					speaker.Play(smpEnter.streamer())
				}
			case kDelete:
				speaker.Play(smpDelete.streamer())
			default:
				// speaker.Play(smp.streamer())
				speaker.Play(smpDown.streamer())
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
				speaker.Play(smpUp.streamer())
			}
		}
	}
}
