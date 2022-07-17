package main

import (
	"./util"
	"./walks"
	"encoding/json"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"
)

const maxX = util.MaxX
const maxY = util.MaxY

func genColorKite(i, j int, dimMin, dimMax float64, applyMultiplier bool, c *[maxX][maxY]float64, conf *util.Configuration, attractors *[]util.Point) float64 {
	bigDist := conf.BigDist
	smallDist := conf.SmallDist
	colorWalk := conf.ColorWalk
	odds := conf.Odds

	if c[i][j] != 0 {
		return c[i][j]
	}

	var walk func(util.Walk) (int, int, int)

	chance := rand.Intn(999)

	if chance == 0 {
		return rcolorf64(dimMin, dimMax)
	}

	var dist int

	if chance < odds["big_walk"] {
		walk = getWalk("big", conf)
		dist = bigDist
	} else {
		walk = getWalk("small", conf)
		dist = smallDist
	}

	walkParams := util.Walk{
		I:          i,
		J:          j,
		Dist:       dist,
		Attractors: *attractors,
	}
	newX, newY, attrDist := walk(walkParams)

	colorMultiplier := 1.0

	if applyMultiplier {
		colorMultiplier = get_color_multiplier(conf, attrDist)
	}

	colorWalkValue := rand_color_walk(colorWalk, odds) * colorMultiplier

	if !applyMultiplier {
		colorWalkValue = 0
	}

	//color := (gen_color_kite(new_x, new_y, in_color, c, conf, attractors) + color_walk_value) % 255
	color := genColorKite(newX, newY, dimMin, dimMax, applyMultiplier, c, conf, attractors) + colorWalkValue

	for color > dimMax {
		color -= dimMax
	}
	for color < dimMin {
		color += dimMin
	}

	xRange := conf.KiteRange["x"]
	yRange := conf.KiteRange["y"]

	c[i][j] = color

	for y := -yRange; y < yRange; y++ {
		for x := -xRange; x < xRange; x++ {
			if x == y {
				continue
			}

			newI := (i + x) % maxX
			newJ := (j + y) % maxY

			if newI < 0 {
				newI += maxX
			}
			if newJ < 0 {
				newJ += maxY
			}

			if c[newI][newJ] == 0 {
				c[newI][newJ] = color
			}
		}
	}

	if util.Abs(conf.Mirror) == 2 {
		applyMirror(i, j, c, color, conf)
	} else if conf.Nup == 2 {
		apply_nup(i, j, c, color)
	}

	return color
}

func get_color_multiplier(conf *util.Configuration, attr_dist int) float64 {
	var color_multiplier float64
	if !conf.ApplyMultiplier || attr_dist == -1 || attr_dist > 200 {
		color_multiplier = 1
	} else if attr_dist < 50 {
		color_multiplier = 4
	} else if attr_dist < 100 {
		color_multiplier = 3
	} else if attr_dist < 200 {
		color_multiplier = 2
	}
	return color_multiplier
}

func genColor(i, j int, dimMin, dimMax float64, applyMultiplier bool, c *[maxX][maxY]float64, conf *util.Configuration, attractors *[]util.Point) float64 {
	bigDist := conf.BigDist
	smallDist := conf.SmallDist
	colorWalk := conf.ColorWalk
	odds := conf.Odds

	if c[i][j] == 0 {
		chance := rand.Intn(999)
		var color float64
		if chance != 0 {
			var walk func(util.Walk) (int, int, int)
			var dist int

			if chance < odds["big_walk"] {
				walk = getWalk("big", conf)
				dist = bigDist
			} else {
				walk = getWalk("small", conf)
				dist = smallDist
			}
			walkParams := util.Walk{
				I:          i,
				J:          j,
				Dist:       dist,
				Attractors: *attractors,
			}
			newX, newY, attrDist := walk(walkParams)
			if newX == i && newY == j {
				color = rcolorf64(dimMin, dimMax)
			} else {
				color_multiplier := 1.0

				if applyMultiplier == true {
					color_multiplier = get_color_multiplier(conf, attrDist)
				}

				if applyMultiplier == true {
					color = genColor(newX, newY, dimMin, dimMax, applyMultiplier, c, conf, attractors) + (rand_color_walk(colorWalk, odds) * color_multiplier)
				} else {
					color = genColor(newX, newY, dimMin, dimMax, applyMultiplier, c, conf, attractors)
				}

				for color > dimMax {
					color -= dimMax
				}
				for color < dimMin {
					color += dimMin
				}

				//color = (gen_color(new_x, new_y, in_color, c, conf, attractors) + (rand_color_walk(color_walk, odds) * color_multiplier  ) ) % 255
			}
		} else {
			color = rcolorf64(dimMin, dimMax)
		}

		c[i][j] = color

		if util.Abs(conf.Mirror) == 2 {
			applyMirror(i, j, c, color, conf)
		} else if conf.Nup == 2 {
			apply_nup(i, j, c, color)
		}
	}

	return c[i][j]
}

func getWalk(walkType string, conf *util.Configuration) func(util.Walk) (int, int, int) {
	if val, ok := conf.Walks[walkType]; ok {
		switch strings.ToLower(val) {
		case "random":
			return walks.Random
		case "isotropic":
			return walks.Isotropic
		case "center_isotropic":
			return walks.CenterIsotropic
		case "neighbor":
			return walks.Neighbor
		case "center_neighbor":
			return walks.Center_neighbor_walk
		case "random_center_neighbor":
			return walks.Random_center_neighbor_walk
		case "random_attractor_neighbor":
			return walks.Random_attractor_neighbor_walk
		case "random_attractor_neighbor_twist":
			return walks.Random_attractor_neighbor_twist
		case "isotropic_attractor_neighbor":
			return walks.Isotropic_attractor_neighbor_walk
		}
	}

	return walks.Isotropic
}

func apply_nup(i, j int, c *[maxX][maxY]float64, color float64) {
	var new_i int
	half_x := maxX / 2
	if i > half_x {
		new_i = i - half_x
	} else {
		new_i = i + half_x
	}
	if new_i > 0 && new_i < maxX && c[new_i][j] == 0 {
		c[new_i][j] = color
	}
}
func applyMirror(i, j int, c *[maxX][maxY]float64, color float64, conf *util.Configuration) {
	var new_i int
	halfX := maxX / 2
	if i < halfX {
		new_i = maxX - i
	} else if i > halfX {
		new_i = maxX - i
	} else {
		new_i = 0
	}
	if new_i >= 0 && new_i < maxX && c[new_i][j] == 0 && color > 0 {
		if conf.Mirror > 0 {
			c[new_i][j] = color
		} else {
			c[new_i][j] = float64((uint8(color) + 128) % 255)
		}
	}
}

func writeImage(name string, m image.Image) {
	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	png.Encode(out, m)
}

func getConf(filename string) util.Configuration {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error:", err)
	}
	decoder := json.NewDecoder(file)
	var configuration util.Configuration
	err2 := decoder.Decode(&configuration)
	if err2 != nil {
		fmt.Println("error:", err2)
	}
	return configuration
}

func setPix(m *image.NRGBA, i int, v_red, v_green, v_blue uint8) {
	m.Pix[i] = v_red
	m.Pix[i+1] = v_green
	m.Pix[i+2] = v_blue
	m.Pix[i+3] = 255
}

func genAttractors(num int) []util.Point {
	result := make([]util.Point, num)

	for i := 0; i < num; i++ {
		result[i] = util.Point{X: util.Locrand(maxX), Y: util.Locrand(maxY)}
	}

	return result
}

func getRgb(x, y int, r, g, b *[maxX][maxY]float64, setup util.RgbSetup) colorful.Color {
	vRed := setup.ColorGen(x, y, 0, 255, true, r, setup.Conf, setup.Attractors)
	vGreen := setup.ColorGen(x, y, 0, 255, true, g, setup.Conf, setup.Attractors)
	vBlue := setup.ColorGen(x, y, 0, 255, true, b, setup.Conf, setup.Attractors)
	return colorful.Color{vRed / 255.0, vGreen / 255.0, vBlue / 255.0}
}

func getHcl(x, y int, h, c, l *[maxX][maxY]float64, setup util.RgbSetup) colorful.Color {

	hue := setup.ColorGen(x, y, 0, 360, true, h, setup.Conf, setup.Attractors)
	chroma := setup.ColorGen(x, y, -1, 1, false, c, setup.Conf, setup.Attractors)
	lightness := setup.ColorGen(x, y, 0, 1, false, l, setup.Conf, setup.Attractors)

	return colorful.Hcl(hue, chroma, lightness)
}

func main() {
	var r [maxX][maxY]float64
	var g [maxX][maxY]float64
	var b [maxX][maxY]float64

	conf := getConf("./conf.json")
	attractors := genAttractors(conf.Attractors)

	conf.Seed = util.GetSeed(conf)
	rand.Seed(conf.Seed)

	setup := util.RgbSetup{
		BaseColor:  rcolor(255),
		ColorGen:   getColorGenFunc(conf),
		Conf:       &conf,
		Attractors: &attractors,
	}

	getColorFunc := determineColorSpaceFunc(conf.ColorSpace)

	m := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))
	mR := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))
	mG := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))
	mB := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))

	mRg := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))
	mRb := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))
	mGb := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {

			rgb := getColorFunc(x, y, &r, &g, &b, setup)
			v_red, v_green, v_blue := rgb.Clamped().RGB255()

			i := y*m.Stride + x*4

			setPix(m, i, v_red, v_green, v_blue)
			setPix(mR, i, v_red, 0, 0)
			setPix(mG, i, 0, v_green, 0)
			setPix(mB, i, 0, 0, v_blue)

			setPix(mRg, i, v_red, v_green, 0)
			setPix(mRb, i, v_red, 0, v_blue)
			setPix(mGb, i, 0, v_green, v_blue)
		}
	}

	filename := fmt.Sprintf("%v_%v_%v_%v_%v-%v-%v-%vx%v",
		conf.BigDist,
		conf.SmallDist,
		conf.ColorWalk,
		conf.Attractors,
		conf.Mirror,
		conf.Seed,
		conf.ColorSpace,
		maxX,
		maxY)

	fmt.Println("Writing config")
	writeConfig("output/"+filename+".json", conf)
	fmt.Println("Writing composite view")
	writeImage("test.png", m)
	fmt.Println("Writing composite")
	writeImage("output/"+filename+".png", m)

	if conf.OutputComponent {
		fmt.Println("Writing red")
		writeImage("output/"+filename+"-red.png", mR)
		fmt.Println("Writing green")
		writeImage("output/"+filename+"-green.png", mG)
		fmt.Println("Writing blue")
		writeImage("output/"+filename+"-blue.png", mB)

		fmt.Println("Writing red-green")
		writeImage("output/"+filename+"-red-green.png", mRg)
		fmt.Println("Writing red-blue")
		writeImage("output/"+filename+"-red-blue.png", mRb)
		fmt.Println("Writing green-blue")
		writeImage("output/"+filename+"-green-blue.png", mGb)
	}

	fmt.Println(filename)

	fmt.Println("done")
}

func determineColorSpaceFunc(colorSpace string) func(x, y int, r, g, b *[maxX][maxY]float64, setup util.RgbSetup) colorful.Color {
	if colorSpace == "hcl" {
		return getHcl
	}
	return getRgb
}

func getColorGenFunc(conf util.Configuration) (colorGen func(int, int, float64, float64, bool, *[maxX][maxY]float64, *util.Configuration, *[]util.Point) float64) {
	if conf.ColorGen == "kite" {
		colorGen = genColorKite
	} else {
		colorGen = genColor
	}
	return colorGen
}

func writeConfig(filename string, conf util.Configuration) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("error:", err)
	}
	encoder := json.NewEncoder(file)
	err2 := encoder.Encode(&conf)
	if err2 != nil {
		fmt.Println("error:", err2)
	}
}
