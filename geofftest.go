package main

import (
	"fmt"
	"image"
	"image/png"
	"encoding/json"
	"log"
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

func gen_color(i, j int, c *[max_x][max_y]uint8, conf *Configuration) uint8 {
	big_dist := conf.BigDist
	small_dist := conf.SmallDist
	color_walk := conf.ColorWalk
	if c[i][j] == 0 {
		chance := rand.Intn(999)
		var color uint8
		if chance != 0 {
			if chance < 500 {
				new_x := (i + locrand(big_dist * 2) - big_dist + max_x - 2) % max_x
				new_y := (j + locrand(big_dist * 2) - big_dist - max_y - 2) % max_y
				if new_x < 0 { new_x += max_x }
				if new_y < 0 { new_y += max_y }
				color = (gen_color(new_x, new_y, c, conf) + color_walk) % 255
			} else {
				new_x := (i + (locrand(small_dist * 2) - small_dist)) % max_x
				new_y := (j + (locrand(small_dist * 2) - small_dist)) % max_y
				if new_x < 0 { new_x += max_x }
				if new_y < 0 { new_y += max_y }
				color = (gen_color(new_x, new_y, c, conf) + color_walk) % 255
			}
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
}

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


func main() {
	var r [max_x][max_y]uint8
	var g [max_x][max_y]uint8
	var b [max_x][max_y]uint8
	rand.Seed( time.Now().UTC().UnixNano())
	conf := getconf("./conf.json")

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
	       		m.Pix[i] = v_red
			m.Pix[i+1] = v_green
			m.Pix[i+2] = v_blue
			m.Pix[i+3] = 255

			m_r.Pix[i] = v_red
			m_r.Pix[i+1] = 255
			m_r.Pix[i+2] = 255
			m_r.Pix[i+3] = 255
			
			m_g.Pix[i] = v_green % 255
			m_g.Pix[i+1] = 255
			m_g.Pix[i+2] = v_green % 255
			m_g.Pix[i+3] = 255
			
			m_b.Pix[i] = 255
			m_b.Pix[i+1] = 255
			m_b.Pix[i+2] = v_blue
			m_b.Pix[i+3] = 255
		   }
	   }

	write_image("test.png", m)
	write_image("test_red.png", m_r)
	write_image("test_green.png", m_g)
	write_image("test_blue.png", m_b)

	fmt.Println("done")

}
