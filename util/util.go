package util

import (
	"math"
	"math/rand"
	"time"
)


//const MaxX = 1024
//const MaxY = 1024

const MaxX = 1920
const MaxY = 1080

type Point struct {
	X int
	Y int
}

type Configuration struct {
	BigDist int `json:"big_dist"`
	SmallDist	int `json:"small_dist"`
	ColorWalk	uint8 `json:"color_walk"`
	Odds Odds	`json:"odds"`
	Attractors int	`json:"attractors"`
	Mirror 	int	`json:"mirror"`
	Nup int	`json:"nup"`
	KiteRange Range	`json:"kite"`
	Seed int64 `json:"seed"`
	ColorGen string `json:"color_gen"`
	ColorSpace string `json:"color_space"`
	ApplyMultiplier bool `json:"apply_multiplier"`
	OutputComponent bool `json:"output_component"`
	Walks map[string]string `json:"walks"`
}

type Odds map[string]int

type Range map[string]int

type Walk struct {
	I int
	J int
	Dist int
	Attractors []Point
}

type RgbSetup struct {
	BaseColor  uint8
	ColorGen   func(int, int, float64, float64, bool, *[MaxX][MaxY]float64, *Configuration, *[]Point) float64
	Conf       *Configuration
	Attractors *[]Point
}

func Locrand(n int) int {
	return rand.Intn(n)
}


func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CalcDist(x, y, i, j int) int {
	d_x := Abs(x - i)
	d_y := Abs(y - j)
	dist := math.Hypot(float64(d_x), float64(d_y))
	return int(dist)
}

func GetSeed(conf Configuration) int64 {
	if conf.Seed != 0 { return conf.Seed }
	return time.Now().UTC().UnixNano()
}