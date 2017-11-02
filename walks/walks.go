package walks

import (
	"math/rand"
	"imagefun/util"
)

const max_x = util.MaxX
const max_y = util.MaxY


func Random(walk util.Walk) (int, int, int) {
	return rand.Intn(walk.Dist), rand.Intn(walk.Dist), rand.Intn(max_x/2)
}

func Isotropic(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist

	new_x := (i + util.Locrand(dist * 2) - dist + max_x - 2) % max_x
	new_y := (j + util.Locrand(dist * 2) - dist + max_y - 2) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y, -1
}

func CenterIsotropic(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist

	move_x := util.Locrand(dist) + max_x - 2
	move_y := util.Locrand(dist) + max_y - 2

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }
	if util.Locrand(10) < 3 {
		move_x = -move_x
	}
	if util.Locrand(10) < 3 {
		move_y = -move_y
	}

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	center_dist := util.CalcDist(max_x / 2, max_y / 2, i, j)
	return new_x, new_y, center_dist
}

func Neighbor(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist

	new_x := (i + (util.Locrand(dist * 2) - dist)) % max_x
	new_y := (j + (util.Locrand(dist * 2) - dist)) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	return new_x, new_y, -1
}

func Center_neighbor_walk(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist

	var move_x, move_y int
	if dist == 1 {
		move_x = 1
		move_y = 1
	} else {
		move_x = util.Locrand(dist)
		move_y = util.Locrand(dist)
	}

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }
	center_dist := util.CalcDist(max_x / 2, max_y / 2, i, j)
	return new_x, new_y, center_dist
}

func Random_center_neighbor_walk(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist

	move_x := util.Locrand(dist)
	move_y := util.Locrand(dist)

	if i > max_x / 2 { move_x = -move_x }
	if j > max_y / 2 { move_y = -move_y }

	if util.Locrand(10) < 3 {
		move_x = -move_x
	}
	if util.Locrand(10) < 3 {
		move_y = -move_y
	}

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 { new_x += max_x }
	if new_y < 0 { new_y += max_y }

	center_dist := util.CalcDist(max_x / 2, max_y / 2, i, j)

	return new_x, new_y, center_dist
}

func Isotropic_attractor_neighbor_walk(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist
	attractors := walk.Attractors

	var move_x, move_y int

	for move_x == 0 && move_y == 0 {
		if dist == 1 {
			move_x = 1
			move_y = 1
		} else {
			move_x = (util.Locrand(dist * 2) - dist)
			move_y = (util.Locrand(dist * 2) - dist)
		}
	}

	var closest *util.Point
	var second_closest *util.Point
	var closest_dist int
	var second_closest_dist int
	for k := 0; k < len(attractors); k++ {
		dist := util.CalcDist(attractors[k].X, attractors[k].Y, i, j)
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

	closest_rand := util.Locrand(int(closest_dist))
	second_closest_rand := util.Locrand(int(second_closest_dist))

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


func Random_attractor_neighbor_walk(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist
	attractors := walk.Attractors

	var move_x, move_y int

	for move_x == 0 && move_y == 0 {
		if dist == 1 {
			move_x = 1
			move_y = 1
		} else {
			move_x = util.Locrand(dist)
			move_y = util.Locrand(dist)
		}

	}

	var closest *util.Point
	var second_closest *util.Point
	var closest_dist int
	var second_closest_dist int
	for k := 0; k < len(attractors); k++ {
		dist := util.CalcDist(attractors[k].X, attractors[k].Y, i, j)
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

	closest_rand := util.Locrand(int(closest_dist))
	second_closest_rand := util.Locrand(int(second_closest_dist))

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

func Random_attractor_neighbor_twist(walk util.Walk) (int, int, int) {
	i := walk.I
	j := walk.J
	dist := walk.Dist
	attractors := walk.Attractors

	var move_x, move_y, twist_x, twist_y, rand_dist int

	for move_x == 0 && move_y == 0 {

		if dist == 1 {
			rand_dist = 1
		} else {
			rand_dist = util.Locrand(dist)
		}

		if 50 < rand.Intn(99) {
			twist_x = 1
			twist_y = -1
		} else {
			twist_x = -1
			twist_y = 1
		}

		move_x = rand_dist + twist_x
		move_y = rand_dist + twist_y
	}

	var closest *util.Point
	var second_closest *util.Point
	var closest_dist int
	var second_closest_dist int
	for k := 0; k < len(attractors); k++ {
		dist := util.CalcDist(attractors[k].X, attractors[k].Y, i, j)
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

	if closest_dist == 0 {
		closest_dist = 1
	}
	if second_closest_dist == 0 {
		second_closest_dist = closest_dist
	}

	closest_rand := util.Locrand(int(closest_dist))
	second_closest_rand := util.Locrand(int(second_closest_dist))

	var attr_dist int
	if second_closest_rand < closest_rand {
		if i > second_closest.X {
			move_x = -move_x
		}
		if j > second_closest.Y {
			move_y = -move_y
		}
		attr_dist = int(second_closest_dist)
	} else {
		if i > closest.X {
			move_x = -move_x
		}
		if j > closest.Y {
			move_y = -move_y
		}
		attr_dist = int(closest_dist)
	}

	new_x := (i + move_x) % max_x
	new_y := (j + move_y) % max_y
	if new_x < 0 {
		new_x += max_x
	}
	if new_y < 0 {
		new_y += max_y
	}
	return new_x, new_y, attr_dist
}
