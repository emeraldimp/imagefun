package main

import (
	"fmt"
	"image"
	"image/png"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const max_x = 1024
const max_y = 1024

func rcolor(n int) uint8 {
	return uint8(rand.Intn(n))
}

func locrand(n int) int {
	return rand.Intn(n)
}

func isotropic_walk(i, j, dist int) (int, int) {
	new_x := (i + locrand(dist * 2) - dist + max_x - 2) % max_x
	new_y := (j + locrand(dist * 2) - dist + max_y - 2) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y
}

func center_isotropic_walk(i, j, dist int) (int, int) {
	move_x := locrand(dist) + max_x - 2
	move_y := locrand(dist) + max_y - 2

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }
	if locrand(10) < 3 {
		move_x = -move_x
	}
	if locrand(10) < 3 {
		move_y = -move_y
	}

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y
}

func neighbor_walk(i, j, dist int) (int, int) {
	new_x := (i + (locrand(dist * 2) - dist)) % max_x
	new_y := (j + (locrand(dist * 2) - dist)) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y
}

func Abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func center_neighbor_walk(i, j, dist int) (int, int) {
	move_x := locrand(dist)
	move_y := locrand(dist)

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y
}

func random_center_neighbor_walk(i, j, dist int) (int, int) {
	move_x := locrand(dist)
	move_y := locrand(dist)

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }

	if locrand(10) < 3 {
		move_x = -move_x
	}
	if locrand(10) < 3 {
		move_y = -move_y
	}

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y
}


func random_attractor_neighbor_walk(i, j, dist int) (int, int) {
	move_x := locrand(dist)
	move_y := locrand(dist)

	var closest Point
	var closest_dist float64
	var second_closest_dist float64
	for k := 0; k < len(attractors); k++ {
		d_x := Abs(attractors[k].X - i)
		d_y := Abs(attractors[k].Y - j)
		dist := math.Hypot(float64(d_x), float64(d_y))
		if dist < closest_dist || closest_dist == 0 {
			second_closest_dist = closest_dist
			closest_dist = dist
			closest = attractors[k]
		}
	}

	rand_close := locrand(100) + 100

	if int(closest_dist) < rand_close + int(second_closest_dist) {
		if i > closest.X { move_x = -move_x }
		if j > closest.Y { move_y = -move_y }
	}

	if int(closest_dist) < 50 {
		move_x = move_x * 2
		move_y = move_y * 2
	}


	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y
}



func simple_color_walk(color uint8) uint8 {
	return color
}

func rand_color_walk(color uint8) uint8 {
	if color == 0 { return color }
	return rcolor(int(color))
}

func gen_color(i, j int, c *[max_x][max_y]uint8, conf *Configuration) uint8 {
	big_dist := conf.BigDist
	small_dist := conf.SmallDist
	color_walk := conf.ColorWalk
	odds := conf.Odds
	if c[i][j] == 0 {
		chance := rand.Intn(999)
		var color uint8
		if chance != 0 {
			var walk func(int, int, int) (int, int)
			var dist int

			if chance < odds["isotropic_walk"] {
				//walk = isotropic_walk
				walk = center_isotropic_walk
				dist = big_dist
			} else {
				//walk = neighbor_walk
				//walk = center_neighbor_walk
				walk = random_attractor_neighbor_walk
				dist = small_dist
			}
			new_x, new_y := walk(i, j, dist)
			color = (gen_color(new_x, new_y, c, conf) + rand_color_walk(color_walk)) % 255
		} else {
			color = rcolor(256)
		}
		c[i][j] = color
	}

	return c[i][j]
}

func write_image(name string, m image.Image) {
  	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
   	defer out.Close()

	png.Encode(out, m)
}

type Configuration struct {
	BigDist int `json:"big_dist"`
	SmallDist	int `json:"small_dist"`
	ColorWalk	uint8 `json:"color_walk"`
	Odds Odds	`json:"odds"`
	Attractors int	`json:"attractors"`
}

type Odds map[string]int 

func getconf(filename string) Configuration {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error:", err)
	}
	decoder := json.NewDecoder(file)
	var configuration Configuration
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

func genAttractors(num int) []Point {
	result := make([]Point, num)

	for i := 0; i < num; i++ {
		result[i] = Point{locrand(max_x), locrand(max_y)}
	}

	return result
}

type Point struct {
	X int
	Y int
}

var attractors []Point

func main() {
	var r [max_x][max_y]uint8
	var g [max_x][max_y]uint8
	var b [max_x][max_y]uint8
	rand.Seed( time.Now().UTC().UnixNano())
	conf := getconf("./conf.json")
	attractors = genAttractors(conf.Attractors)
fmt.Println(attractors)

	m := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_r := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_g := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	m_b := image.NewNRGBA(image.Rect(0, 0, max_x, max_y))
	for y := 0; y < max_y; y++ {
		for x := 0; x < max_x; x++ {
			v_red := gen_color(x, y, &r, &conf)
			v_green := gen_color(x, y, &g, &conf)
			v_blue := gen_color(x, y, &b, &conf)
	       		i := y * m.Stride + x*4

			setPix(m, i, v_red, v_green, v_blue)
			setPix(m_r, i, v_red, 255, 255)
			setPix(m_g, i, 255, v_green, 255)
			setPix(m_b, i, 255, 255, v_blue)
		   }
	   }

	write_image("test.png", m)
	write_image("test_red.png", m_r)
	write_image("test_green.png", m_g)
	write_image("test_blue.png", m_b)

	fmt.Println("done")

}
