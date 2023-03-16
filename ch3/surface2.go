package main

import (
	"errors"
	"fmt"
	"image/color"
	"math"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	var minZ, maxZ float64
	var minColor, maxColor color.RGBA

	// Find minimum and maximum z values for coloring the peaks and valleys
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			_, _, color := corner(i, j)
			z, _ := f(xyrange*(float64(i)/cells-0.5), xyrange*(float64(j)/cells-0.5))
			if z < minZ {
				minZ = z
				minColor = color
			}
			if z > maxZ {
				maxZ = z
				maxColor = color
			}
		}
	}

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, acolor := corner(i+1, j)
			bx, by, bcolor := corner(i, j)
			cx, cy, ccolor := corner(i, j+1)
			dx, dy, dcolor := corner(i+1, j+1)

			// Use the minimum and maximum z values to color the peaks and valleys
			z, _ := f(xyrange*(float64(i)/cells-0.5), xyrange*(float64(j)/cells-0.5))
			if z < 0 {
				acolor = averageColor(acolor, minColor)
				bcolor = averageColor(bcolor, minColor)
				ccolor = averageColor(ccolor, minColor)
				dcolor = averageColor(dcolor, minColor)
			} else {
				acolor = averageColor(acolor, maxColor)
				bcolor = averageColor(bcolor, maxColor)
				ccolor = averageColor(ccolor, maxColor)
				dcolor = averageColor(dcolor, maxColor)
			}

			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' style='fill:#%02x%02x%02x'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, acolor.R, acolor.G, acolor.B)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, color.RGBA) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z for the four vertices of the cell.
	z, err := f(x, y)
	if err != nil {
		fmt.Println(err)
		return 0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}

	zTopLeft, _ := f(x-0.5*xyrange/cells, y-0.5*xyrange/cells)
	zTopRight, _ := f(x+0.5*xyrange/cells, y-0.5*xyrange/cells)
	zBottomLeft, _ := f(x-0.5*xyrange/cells, y+0.5*xyrange/cells)
	zBottomRight, _ := f(x+0.5*xyrange/cells, y+0.5*xyrange/cells)

	// Determine whether each vertex is a peak or a valley.
	color := color.RGBA{R: 255, G: 255, B: 255, A: 255} // default color

	if z >= zTopLeft && z >= zTopRight && z >= zBottomLeft && z >= zBottomRight {
		color.R = 255 // Red for peaks
		color.G = 0
		color.B = 0
	} else if z <= zTopLeft && z <= zTopRight && z <= zBottomLeft && z <= zBottomRight {
		color.R = 0 // Blue for valleys
		color.G = 0
		color.B = 255
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, color
}

func f(x, y float64) (float64, error) {
	r := math.Hypot(x, y) // distance from (0,0)
	result := math.Sin(r) / r
	if math.IsInf(result, 0) {
		return 0, errors.New("El resultado tiene infinitos decimales")
	}
	return result, nil
}

func averageColor(vertices ...color.RGBA) color.RGBA {
	var rSum, gSum, bSum int
	for _, vertex := range vertices {
		rSum += int(vertex.R)
		gSum += int(vertex.G)
		bSum += int(vertex.B)
	}
	count := len(vertices)
	return color.RGBA{
		uint8(rSum / count),
		uint8(gSum / count),
		uint8(bSum / count),
		255, // alpha value
	}
}
