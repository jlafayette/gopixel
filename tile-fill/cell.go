package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

// Cell ...
type Cell struct {
	seedV     pixel.Vec
	centroidV pixel.Vec
	points    points
	seedX     int
	seedY     int
}

// NewCell ...
func NewCell(x, y int) Cell {
	return Cell{
		seedV: pixel.V(float64(x), float64(y)),
		seedX: x,
		seedY: y,
	}
}

func (c *Cell) reset(x, y int) {
	c.seedV = pixel.V(float64(x), float64(y))
	c.seedX = x
	c.seedY = y
	c.points = nil
}

func (c *Cell) update() {
	c.orderPoints()
	c.computeCentroid()
}

func (c *Cell) computeCentroid() {

}

func (c *Cell) addPoint(v pixel.Vec) {
	c.points = append(c.points, pt{v, &c.seedV})
}

// order points clockwise around seed vector
func (c *Cell) orderPoints() {
	sort.Sort(c.points)
}

type points []pt

func (p points) Len() int {
	return len(p)
}
func (p points) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p points) Less(i, j int) bool {
	return p[i].angle() < p[j].angle()
}

type pt struct {
	v     pixel.Vec
	seedV *pixel.Vec
}

func (p pt) angle() float64 {
	return p.seedV.Sub(p.v).Angle()
}

func (c *Cell) createOutline(imd *imdraw.IMDraw) {
	c.orderPoints()
	for _, pt := range c.points {
		imd.Push(pt.v)
	}
	imd.Polygon(2)
}

func (c *Cell) createOffsetOutline(imd *imdraw.IMDraw) {
	c.orderPoints()
	for _, pt := range c.points {
		nudgeV := c.seedV.Sub(pt.v).Unit().Scaled(2)
		imd.Push(pt.v.Add(nudgeV))
	}
	imd.Polygon(1)
}

func (c *Cell) createSpokes(imd *imdraw.IMDraw) {
	for _, pt := range c.points {
		imd.Push(c.seedV, pt.v)
		imd.Line(1)
	}
}

func (c *Cell) draw(imd *imdraw.IMDraw) {
	c.createOutline(imd)
	imd.Push(c.seedV)
	imd.Circle(2, 0)
}

func (c *Cell) drawDebug(imd *imdraw.IMDraw) {
	// c.createSpokes(imd)
	// c.createOffsetOutline(imd)
	// imd.Push(c.seedV)
	// imd.Circle(7, 2)

	imd.Push(c.seedV)
	imd.Circle(3, 0)

	// figuring out centroid visually here first

	type triangle struct {
		centroid pixel.Vec
		area     float64
	}
	var tris []triangle

	end := 2
	a := c.points[0].v
	for end < len(c.points) {
		// centroid
		b := c.points[end-1].v
		c := c.points[end].v
		x := (a.X + b.X + c.X) / 3
		y := (a.Y + b.Y + c.Y) / 3

		// area
		// 1/2 * | x1y2 + x2y3 + x3y1 - x2y1 - x3y2 - x1y3 |
		// 1/2 * | a.X*b.Y + b.X*c.Y + c.X*a.Y - b.X*a.Y - c.X*b.Y - a.X*c.Y |
		area := .5 * math.Abs(a.X*b.Y+b.X*c.Y+c.X*a.Y-b.X*a.Y-c.X*b.Y-a.X*c.Y)

		t := triangle{pixel.V(x, y), area}
		tris = append(tris, t)

		imd.Push(t.centroid)
		imd.Circle(5, 1)

		imd.Push(a)
		imd.Push(b)
		imd.Push(c)
		imd.Polygon(1)
		end++
	}

	mainCentroid := tris[0].centroid
	totalArea := tris[0].area
	// start with main centroid being tri[0].centroid
	// for each additional triangle, lerp it by the amount of area normalized?
	for i := 1; i < len(tris); i++ {

		// solve for z
		// 0 <= z <= 1
		// a > 0
		// b > 0

		// a = 1, b = 1, z = .5
		// a = 2, b = 2, z = .5

		// z = .5 * a / b

		// what times the bigger = smaller (always between 0-1)
		// small = big * f
		// f = small / big
		var f float64
		if totalArea >= tris[i].area {
			big := totalArea
			small := tris[i].area
			f = .5 * (small / big)
			mainCentroid = pixel.Lerp(mainCentroid, tris[i].centroid, f)
		} else {
			big := tris[i].area
			small := totalArea
			f = .5 * (small / big)
			mainCentroid = pixel.Lerp(tris[i].centroid, mainCentroid, f)
		}
		fmt.Println(f)
		totalArea = totalArea + tris[i].area
	}

	imd.Push(mainCentroid)
	imd.Circle(8, 0)
}
