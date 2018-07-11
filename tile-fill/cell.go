package main

import (
	"sort"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
)

// Cell ...
type Cell struct {
	seedV  pixel.Vec
	points points
	seedX  int
	seedY  int
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
	c.createSpokes(imd)
	c.createOffsetOutline(imd)
	imd.Push(c.seedV)
	imd.Circle(7, 2)
}
