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

func isotropic_walk(i, j, dist int) (int, int, int) {
	new_x := (i + locrand(dist * 2) - dist + max_x - 2) % max_x
	new_y := (j + locrand(dist * 2) - dist + max_y - 2) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y, -1
}

func center_isotropic_walk(i, j, dist int) (int, int, int) {
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
	center_dist := calc_dist(max_x / 2, max_y / 2, i, j)
	return new_x, new_y, center_dist
}

func neighbor_walk(i, j, dist int) (int, int, int) {
	new_x := (i + (locrand(dist * 2) - dist)) % max_x
	new_y := (j + (locrand(dist * 2) - dist)) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y, -1
}

func Abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func center_neighbor_walk(i, j, dist int) (int, int, int) {
	var move_x, move_y int
	if dist == 1 {
		move_x = 1
		move_y = 1
	} else {
		move_x = locrand(dist)
		move_y = locrand(dist)
	}	

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	center_dist := calc_dist(max_x / 2, max_y / 2, i, j)
	return new_x, new_y, center_dist
}

func random_center_neighbor_walk(i, j, dist int) (int, int, int) {
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

	center_dist := calc_dist(max_x / 2, max_y / 2, i, j)

	return new_x, new_y, center_dist
}

func calc_dist(x, y, i, j int) int {
	d_x := Abs(x - i)
	d_y := Abs(y - j)
	dist := math.Hypot(float64(d_x), float64(d_y))
	return int(dist)
}


func random_attractor_neighbor_walk(i, j, dist int) (int, int, int) {
	var move_x, move_y int

	for move_x == 0 && move_y == 0 {
		if dist == 1 {
			move_x = 1
			move_y = 1
		} else {
			move_x = locrand(dist)
			move_y = locrand(dist)
		}

	}

	var closest *Point
	var second_closest *Point
	var closest_dist int
	var second_closest_dist int
	for k := 0; k < len(attractors); k++ {
		dist := calc_dist(attractors[k].X, attractors[k].Y, i, j)
		if dist < closest_dist || closest_dist == 0 {
			second_closest = closest
			second_closest_dist = closest_dist
			closest_dist = dist
			closest = &attractors[k]
			if (second_closest == nil) {
				second_closest = closest
				second_closest_dist = dist
			}
		}
	}

	if closest_dist == 0 { closest_dist = 1 }
	if second_closest_dist == 0 { second_closest_dist = closest_dist }

	closest_rand := locrand(int(closest_dist))
	second_closest_rand := locrand(int(second_closest_dist))

	var attr_dist int
	if second_closest_rand < closest_rand {
		if i > second_closest.X { move_x = -move_x }
		if j > second_closest.Y { move_y = -move_y }
		attr_dist = int(second_closest_dist)
	} else {
		if i > closest.X { move_x = -move_x }
		if j > closest.Y { move_y = -move_y }
		attr_dist = int(closest_dist)
	}

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y, attr_dist
}



func simple_color_walk(color uint8) uint8 {
	return color
}

func rand_color_walk(color uint8) uint8 {
	if color == 0 { return color }
	if color == 1 { return color + 1 }
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
			var walk func(int, int, int) (int, int, int)
			var dist int

			if chance < odds["isotropic_walk"] {
				//walk = isotropic_walk
				walk = center_isotropic_walk
				dist = big_dist
			} else {
				//walk = neighbor_walk
				if len(attractors) > 0 {
					walk = random_attractor_neighbor_walk
				} else {
					walk = center_neighbor_walk
				}
				dist = small_dist
			}
			new_x, new_y, attr_dist := walk(i, j, dist)
			if new_x == i && new_y == j { 
				color = rcolor(256) 
			} else {
				var color_multiplier uint8
				if attr_dist == -1 || attr_dist > 200 {
					color_multiplier = 1
				} else if attr_dist < 50 {
					color_multiplier = 4
				} else if attr_dist < 100 {
					color_multiplier = 3
				} else if attr_dist < 200 {
					color_multiplier = 2
				}
				color = (gen_color(new_x, new_y, c, conf) + (rand_color_walk(color_walk) * color_multiplier  ) ) % 255
			}
		} else {
			color = rcolor(256)
		}

		c[i][j] = color

		if Abs(conf.Mirror) == 2 {
			var new_i int
			half_x := max_x / 2

			if i < half_x { new_i = max_x - i } else if i > half_x { new_i = max_x - i } else { new_i = 0}

			if new_i >= 0 && new_i < max_x && c[new_i][j] == 0 && color > 0 {
				if (conf.Mirror > 0) { 
					c[new_i][j] = color
				} else { 
					c[new_i][j] = (color + 128) % 255
				}
			}
		} else if conf.Nup == 2 {
			var new_i int
			half_x := max_x / 2

			if i > half_x { new_i = i - half_x } else { new_i = i + half_x }

			if new_i > 0 && new_i < max_x && c[new_i][j] == 0 {
				c[new_i][j] = color
			}
		}
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
	Mirror 	int	`json:"mirror"`
	Nup int	`json:"nup"`
	Seed int64 `json:"seed"`
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

func getSeed(conf Configuration) int64 {
	if conf.Seed != 0 { return conf.Seed }
	return time.Now().UTC().UnixNano()
}

func main() {
	var r [max_x][max_y]uint8
	var g [max_x][max_y]uint8
	var b [max_x][max_y]uint8
	
	conf := getconf("./conf.json")
	attractors = genAttractors(conf.Attractors)
	
	seed := getSeed(conf)
	rand.Seed(seed) 

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

	filename := fmt.Sprintf("%v_%v_%v_%v_%v-%v.png", 
		conf.BigDist,
		conf.SmallDist,
		conf.ColorWalk,
		conf.Attractors,
		conf.Mirror,
		seed)

	fmt.Println("Writing composite view")
	write_image("test.png", m)
	fmt.Println("Writing composite")
	write_image("output/" + filename, m)
	fmt.Println("Writing red")
	write_image("output/red-" + filename, m_r)
	fmt.Println("Writing green")
	write_image("output/green-" + filename, m_g)
	fmt.Println("Writing blue")
	write_image("output/blue-" + filename, m_b)

	fmt.Println("done")

}
