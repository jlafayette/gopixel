package main

import (
	"sort"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

// Cell ...
type Cell struct {
	center pixel.Vec
	points points
	cx     int
	cy     int
	id     int
}

// NewCell ...
func NewCell(x, y, id int) Cell {
	return Cell{
		center: pixel.V(float64(x), float64(y)),
		cx:     x,
		cy:     y,
		id:     id,
	}
}

func (c *Cell) reset(x, y int) {
	c.center = pixel.V(float64(x), float64(y))
	c.cx = x
	c.cy = y
	c.points = nil
}

func (c *Cell) addPoint(v pixel.Vec) {
	c.points = append(c.points, pt{v, &c.center})
}

// order points clockwise around center
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
	v      pixel.Vec
	center *pixel.Vec
}

func (p pt) angle() float64 {
	return p.center.Sub(p.v).Angle()
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
		nudgeV := c.center.Sub(pt.v).Unit().Scaled(2)
		imd.Push(pt.v.Add(nudgeV))
	}
	imd.Polygon(1)
}

func (c *Cell) createSpokes(imd *imdraw.IMDraw) {
	for _, pt := range c.points {
		imd.Push(c.center, pt.v)
		imd.Line(1)
	}
}

func (c *Cell) draw(imd *imdraw.IMDraw) {
	c.createOutline(imd)
	imd.Push(c.center)
	imd.Circle(5, 0)
}

func (c *Cell) drawDebug(imd *imdraw.IMDraw) {
	c.createSpokes(imd)
	c.createOffsetOutline(imd)
	imd.Push(c.center)
	imd.Circle(7, 2)
}
