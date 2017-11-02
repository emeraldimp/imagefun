package main

import (
	"fmt"
	"image"
	"image/png"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"imagefun/walks"
	"imagefun/util"
	"strings"
	"github.com/lucasb-eyer/go-colorful"
	//"image/color"
)

const max_x = util.MaxX
const max_y = util.MaxY

func gen_color_kite(i, j int, dim_min, dim_max float64, apply_multiplier bool, c *[max_x][max_y]float64, conf *util.Configuration, attractors *[]util.Point) float64 {
	big_dist := conf.BigDist
	small_dist := conf.SmallDist
	color_walk := conf.ColorWalk
	odds := conf.Odds

	if c[i][j] != 0 {
		return c[i][j]
	}

	var walk func(util.Walk) (int, int, int)

	chance := rand.Intn(999)

	if chance == 0 {
		return rcolorf64(dim_min, dim_max)
	}

	var dist int

	if chance < odds["big_walk"] {
		walk = get_walk("big", conf)
		dist = big_dist
	} else {
		walk = get_walk("small", conf)
		dist = small_dist
	}

	walkParams := util.Walk{
		i,
		j,
		dist,
		*attractors,
	}
	new_x, new_y, attr_dist := walk(walkParams)

	color_multiplier := 1.0

	if apply_multiplier {
		color_multiplier = get_color_multiplier(conf, attr_dist)
	}

	color_walk_value := rand_color_walk(color_walk, odds) * color_multiplier

	if !apply_multiplier {
		color_walk_value = 0
	}

	//color := (gen_color_kite(new_x, new_y, in_color, c, conf, attractors) + color_walk_value) % 255
	color := gen_color_kite(new_x, new_y, dim_min, dim_max, apply_multiplier, c, conf, attractors) + color_walk_value

	for color > dim_max {
		color -= dim_max
	}
	for color < dim_min {
		color += dim_min
	}

	x_range := conf.KiteRange["x"]
	y_range := conf.KiteRange["y"]

	c[i][j] = color

	for y := -y_range; y < y_range; y++ {
		for x := -x_range; x < x_range; x++ {
			if x == y {
				continue;
			}

			new_i := (i + x) % max_x
			new_j := (j + y) % max_y

			if new_i < 0 {
				new_i += max_x
			}
			if new_j < 0 {
				new_j += max_y
			}

			if c[new_i][new_j] == 0 {
				c[new_i][new_j] = color
			}
		}
	}

	if util.Abs(conf.Mirror) == 2 {
		apply_mirror(i, j, c, color, conf)
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

func gen_color(i, j int, dim_min, dim_max float64, apply_multiplier bool, c *[max_x][max_y]float64, conf *util.Configuration, attractors *[]util.Point) float64 {
	big_dist := conf.BigDist
	small_dist := conf.SmallDist
	color_walk := conf.ColorWalk
	odds := conf.Odds

	if c[i][j] == 0 {
		chance := rand.Intn(999)
		var color float64
		if chance != 0 {
			var walk func(util.Walk) (int, int, int)
			var dist int

			if chance < odds["big_walk"] {
					walk = get_walk("big", conf)
					dist = big_dist
				} else {
					walk = get_walk("small", conf)
					dist = small_dist
			}
			walkParams := util.Walk{
				i,
				j,
				dist,
				*attractors,
			}
			new_x, new_y, attr_dist := walk(walkParams)
			if new_x == i && new_y == j {
				color = rcolorf64(dim_min, dim_max)
			} else {
				color_multiplier := 1.0

				if apply_multiplier == true	{
					color_multiplier = get_color_multiplier(conf, attr_dist)
				}

				if apply_multiplier == true {
					color = gen_color(new_x, new_y, dim_min, dim_max, apply_multiplier, c, conf, attractors) + (rand_color_walk(color_walk, odds) * color_multiplier  )
				} else {
					color = gen_color(new_x, new_y, dim_min, dim_max, apply_multiplier, c, conf, attractors)
				}

				for color > dim_max {
					color -= dim_max
				}
				for color < dim_min {
					color += dim_min
				}


				//color = (gen_color(new_x, new_y, in_color, c, conf, attractors) + (rand_color_walk(color_walk, odds) * color_multiplier  ) ) % 255
			}
		} else {
			color = rcolorf64(dim_min, dim_max)
		}

		c[i][j] = color

		if util.Abs(conf.Mirror) == 2 {
			apply_mirror(i, j, c, color, conf)
		} else if conf.Nup == 2 {
			apply_nup(i, j,c, color)
		}
	}

	return c[i][j]
}

func get_walk(walk_type string, conf *util.Configuration) (func(util.Walk)(int,int,int)) {
	if val, ok := conf.Walks[walk_type]; ok {
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

func apply_nup(i, j int, c *[max_x][max_y]float64, color float64) {
	var new_i int
	half_x := max_x / 2
	if i > half_x {
		new_i = i - half_x
	} else {
		new_i = i + half_x
	}
	if new_i > 0 && new_i < max_x && c[new_i][j] == 0 {
		c[new_i][j] = color
	}
}
func apply_mirror(i, j int, c *[max_x][max_y]float64, color float64, conf *util.Configuration) {
	var new_i int
	half_x := max_x / 2
	if i < half_x {
		new_i = max_x - i
	} else if i > half_x {
		new_i = max_x - i
	} else {
		new_i = 0
	}
	if new_i >= 0 && new_i < max_x && c[new_i][j] == 0 && color > 0 {
		if conf.Mirror > 0 {
			c[new_i][j] = color
		} else {
			c[new_i][j] =float64((uint8(color) + 128) % 255)
		}
	}
}

func write_image(name string, m image.Image) {
  	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
   	defer out.Close()

	png.Encode(out, m)
}



func getconf(filename string) util.Configuration {
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
		result[i] = util.Point{util.Locrand(max_x), util.Locrand(max_y)}
	}

	return result
}


func getRgb(x,y int, r, g, b *[max_x][max_y]float64, setup util.RgbSetup) colorful.Color {
	vRed := setup.ColorGen(x, y, 0, 255, true, r, setup.Conf, setup.Attractors)
	vGreen := setup.ColorGen(x, y, 0, 255, true, g, setup.Conf, setup.Attractors)
	vBlue := setup.ColorGen(x, y, 0, 255, true, b, setup.Conf, setup.Attractors)
	return colorful.Color{vRed/255.0, vGreen/255.0, vBlue/255.0}
}

func getHcl(x,y int, h, c, l *[max_x][max_y]float64, setup util.RgbSetup) colorful.Color {

	hue := setup.ColorGen(x, y, 0, 360, true, h, setup.Conf, setup.Attractors)
	chroma := setup.ColorGen(x, y, -1, 1, false, c, setup.Conf, setup.Attractors)
	lightness := setup.ColorGen(x, y, 0, 1, false, l, setup.Conf, setup.Attractors)

	return colorful.Hcl(hue, chroma, lightness)
}

func main() {
	var r [max_x][max_y]float64
	var g [max_x][max_y]float64
	var b [max_x][max_y]float64

	conf := getconf("./conf.json")
	attractors := genAttractors(conf.Attractors)

	conf.Seed = util.GetSeed(conf)
	rand.Seed(conf.Seed)

	setup := util.RgbSetup{
		BaseColor:  rcolor(255),
		ColorGen:   get_color_gen_func(conf),
		Conf:       &conf,
		Attractors: &attractors,
	}

	getColorFunc := determineColorSpaceFunc(conf.ColorSpace)

	m := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_r := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_g := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_b := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))

	m_rg := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_rb := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_gb := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))

	for y := 0; y < max_y; y++ {
		for x := 0; x < max_x; x++ {

			rgb := getColorFunc(x, y, &r, &g, &b, setup)
			v_red, v_green, v_blue := rgb.Clamped().RGB255()

			i := y*m.Stride + x*4

			setPix(m, i, v_red, v_green, v_blue)
			setPix(m_r, i, v_red, 0, 0)
			setPix(m_g, i, 0, v_green, 0)
			setPix(m_b, i, 0, 0, v_blue)

			setPix(m_rg, i, v_red, v_green, 0)
			setPix(m_rb, i, v_red, 0, v_blue)
			setPix(m_gb, i, 0, v_green, v_blue)
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
		max_x,
		max_y)

	fmt.Println("Writing config")
	write_config("output/"+filename+".json", conf)
	fmt.Println("Writing composite view")
	write_image("test.png", m)
	fmt.Println("Writing composite")
	write_image("output/"+filename+".png", m)

	if conf.OutputComponent {
		fmt.Println("Writing red")
		write_image("output/"+filename+"-red.png", m_r)
		fmt.Println("Writing green")
		write_image("output/"+filename+"-green.png", m_g)
		fmt.Println("Writing blue")
		write_image("output/"+filename+"-blue.png", m_b)

		fmt.Println("Writing red-green")
		write_image("output/"+filename+"-red-green.png", m_rg)
		fmt.Println("Writing red-blue")
		write_image("output/"+filename+"-red-blue.png", m_rb)
		fmt.Println("Writing green-blue")
		write_image("output/"+filename+"-green-blue.png", m_gb)
	}

	fmt.Println(filename)

	fmt.Println("done")
}

func determineColorSpaceFunc(colorSpace string) func (x,y int, r, g, b *[max_x][max_y]float64, setup util.RgbSetup) colorful.Color {
	if colorSpace == "hcl" {
		return getHcl
	}
	return getRgb
}

func get_color_gen_func(conf util.Configuration) (color_gen func(int, int, float64, float64, bool, *[max_x][max_y]float64, *util.Configuration, *[]util.Point) float64) {
	if conf.ColorGen == "kite" {
		color_gen = gen_color_kite
	} else {
		color_gen = gen_color
	}
	return color_gen
}

func write_config(filename string, conf util.Configuration) {
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
