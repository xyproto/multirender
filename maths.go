package multirender

// Sort2 sorts two numbers
func Sort2(a, b int32) (int32, int32) {
	if a < b {
		return a, b
	}
	return b, a
}

// Sort3 sorts three numbers
func Sort3(a, b, c int32) (int32, int32, int32) {
	if a < b && a < c {
		x, y := Sort2(b, c)
		return a, x, y
	}
	if b < a && b < c {
		x, y := Sort2(a, c)
		return b, x, y
	}
	x, y := Sort2(a, b)
	return c, x, y
}

// MinMax3 finds the smallest and largest of three numbers
func MinMax3(a, b, c int32) (int32, int32) {
	if a < b && a < c {
		if b > c {
			return a, b
		}
		return a, c
	}
	if b < a && b < c {
		if a > c {
			return b, a
		}
		return b, c
	}
	if a > b {
		return c, a
	}
	return c, b
}

// Min3 finds the smallest of three numbers
func Min3(a, b, c int32) int32 {
	if a < b && a < c {
		return a
	}
	if b < a && b < c {
		return b
	}
	return c
}

// Max3 finds the largest of three numbers
func Max3(a, b, c int32) int32 {
	if c > a && c > b {
		return c
	}
	if b > a && b > c {
		return b
	}
	return a
}

// MinMax find the smallest and greatest of two given numbers
func MinMax(a, b int32) (int32, int32) {
	if a < b {
		return a, b
	}
	return b, a
}

// CorrespondingY takes an Y and three coordinate pairs.
// The X of the coordinate pair that matches Y is returned.
// May panic if misused!
func CorrespondingY(ySelector, y1, y2, y3, x1, x2, x3 int32) int32 {
	switch ySelector {
	case y1:
		return x1
	case y2:
		return x2
	case y3:
		return x3
	default:
		panic("No corresponding Y!")
	}
}

// CorrespondingX takes an X and three coordinate pairs.
// The Y of the coordinate pair that matches X is returned.
// May panic if misused!
func CorrespondingX(xSelector, y1, y2, y3, x1, x2, x3 int32) int32 {
	switch xSelector {
	case x1:
		return y1
	case x2:
		return y2
	case x3:
		return y3
	default:
		panic("No corresponding X!")
	}
}

// Lengths returns the x direction and y direction distance between the two points
func Lengths(p1, p2 *Pos) (int32, int32) {
	var xlength, ylength int32
	if p1.x > p2.x {
		xlength = p1.x - p2.x
	} else {
		xlength = p2.x - p1.x
	}
	if p1.y > p2.y {
		ylength = p1.y - p2.y
	} else {
		ylength = p2.y - p1.y
	}
	return xlength, ylength
}

// Interpolate interpolates between two points, with the number of steps equal to the length of the longest stretch
func Interpolate(p1, p2 *Pos) []*Pos {
	var points []*Pos
	xlength, ylength := Lengths(p1, p2)
	if xlength > ylength {
		xstart, xstop := MinMax(p1.x, p2.x)
		ystart, _ := MinMax(p1.y, p2.y)
		y := float32(ystart)
		ystep := float32(ylength) / float32(xlength)
		for x := xstart; x < xstop; x++ {
			y += ystep
			points = append(points, &Pos{x, int32(y)})
		}
		return points
	}
	// ylength >= xlength
	ystart, ystop := MinMax(p1.y, p2.y)
	xstart, _ := MinMax(p1.x, p2.x)
	x := float32(xstart)
	xstep := float32(xlength) / float32(ylength)
	for y := ystart; y < ystop; y++ {
		x += xstep
		points = append(points, &Pos{int32(x), y})
	}
	return points
}
