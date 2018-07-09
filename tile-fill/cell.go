package main

import (
	"fmt"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Cell ...
type Cell struct {
	center pixel.Vec
	points points
	imd    imdraw.IMDraw
	cx     int
	cy     int
	id     int
}

// NewCell ...
func NewCell(x, y, id int) Cell {
	imd := *imdraw.New(nil)
	imd.Color = randomColor()
	imd.EndShape = imdraw.NoEndShape
	return Cell{
		center: pixel.V(float64(x), float64(y)),
		imd:    imd,
		cx:     x,
		cy:     y,
		id:     id,
	}
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

func (c *Cell) createOutline() {
	c.orderPoints()
	for _, pt := range c.points {
		nudgeV := c.center.Sub(pt.v).Unit()
		c.imd.Push(pt.v.Add(nudgeV).Add(nudgeV))
	}
	c.imd.Polygon(2)
}

func (c *Cell) createSpokes() {
	for _, pt := range c.points {
		c.imd.Push(c.center, pt.v)
		c.imd.Line(1)
	}
}

func (c *Cell) draw(tgt pixel.Target) {
	fmt.Printf("calling draw for pt: %v\n", c.center)
	for _, pt := range c.points {
		fmt.Printf("    pt: %v\n", pt)
	}
	c.createSpokes()
	c.createOutline()
	c.imd.Push(c.center)
	c.imd.Circle(3, 0)
	c.imd.Draw(tgt)
}
