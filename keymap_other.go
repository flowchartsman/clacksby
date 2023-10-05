//go:build !darwin

package main

import (
	hook "github.com/robotn/gohook"
)

func rawcodetoKeychar(r uint16) string {
	return hook.RawcodetoKeychar(r)
}

func keychartoRawCode(kc string) uint16 {
	return hook.KeychartoRawcode(kc)
}
