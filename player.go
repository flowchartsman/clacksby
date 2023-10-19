package main

type Player interface {
	KeyDown(key string)
	KeyUp(key string)
}
