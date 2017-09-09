package main

import (
	"math/rand"
	"imagefun/util"
)

func simple_color_walk(color uint8) uint8 {
	return color
}

func rcolor(n int) uint8 {
	return uint8(rand.Intn(n))
}

func rand_color_walk(color_walk uint8, odds util.Odds) uint8 {

	chance := rand.Intn(99)

	if chance < int(odds["color_no_walk"]) {
		return 0
	}

	if color_walk == 0 {
		return 0
	}
	if color_walk == 1 {
		return 1
	}

	return rcolor(int(color_walk))
}