package lec

import (
	"fmt"
	"math"
	"slices"
	"time"
)

func pow(number, power int) int {
	if power == 0 {
		return 1
	}
	result := 1
	for i := 0; i < power; i++ {
		result *= number
	}
	return result
}

type Infinite struct {
	start       time.Time
	length, res time.Duration // TODO: Make another version with fixed length, new data overwritting old data
	levels      [][]int
}

func New(start time.Time, length, resolution time.Duration) *Infinite {
	i := int(math.Ceil(float64(length.Nanoseconds()) / float64(resolution.Nanoseconds())))
	levels := 1
	for j := 1; j < i; j *= 2 {
		levels += 1
	}
	return &Infinite{
		start:  start,
		length: length,
		res:    resolution,
		levels: make([][]int, levels),
	}
}

func (c *Infinite) Save(t time.Time, v int) {
	cell := int(float64(t.Sub(c.start).Nanoseconds()) / float64(c.res)) // virtual index in first layer
	for level := 0; level < len(c.levels); level++ {
		if cell%2 == 0 { // stored cell
			missingcells := cell/2 - (len(c.levels[level]) - 1)
			if missingcells > 0 { // not initialized yet
				m := make([]int, missingcells)
				for j := range missingcells {
					m[j] = 0
				}
				c.levels[level] = slices.Concat(c.levels[level], m)
			}
			c.levels[level][cell/2] += v
		}
		cell /= 2 // floor
	}
}

func (c *Infinite) cellvalue(level, cell int) int {
	if cell%2 == 0 {
		if len(c.levels[level]) < cell/2 { // no data
			return 0
		}
		return c.levels[level][cell/2] // virtual index -> array index
	} else {
		parent := c.cellvalue(level+1, cell/2) // floor
		sibling := c.cellvalue(level, cell-1)
		return parent - sibling
	}
}

func (c *Infinite) query(from, to, level int) (int, error) {
	total := 0
	res := pow(2, level)
	l := int(math.Ceil(float64(from) / float64(res)))
	r := int(math.Floor(float64(to) / float64(res)))
	for i := l; i < r; i++ {
		total += c.cellvalue(level, i)
	}
	lx := l * res
	rx := r * res
	if from != lx {
		if level == 0 {
			return -1, fmt.Errorf("bad developer") // FIXME: there should not be excess area resolved beneath first layer
		}
		q, err := c.query(from, lx, level-1)
		if err != nil {
			return -1, fmt.Errorf("c.query on left: %w", err)
		}
		total += q
	}
	if to != rx {
		if level == 0 {
			return -1, fmt.Errorf("bad developer") // FIXME:
		}
		q, err := c.query(rx, to, level-1)
		if err != nil {
			return -1, fmt.Errorf("c.query on right: %w", err)
		}
		total += q
	}
	return total, nil
}

func (c *Infinite) Query(from, to time.Time) (int, error) {
	if from.Sub(c.start).Nanoseconds()%c.res.Nanoseconds() != 0 {
		return -1, fmt.Errorf("value 'from' doesn't align with resolution")
	} else if to.Sub(c.start).Nanoseconds()%c.res.Nanoseconds() != 0 {
		return -1, fmt.Errorf("value 'to' doesn't align with resolution")
	}

	return c.query(
		int(from.Sub(c.start).Nanoseconds()/c.res.Nanoseconds()),
		int(to.Sub(c.start).Nanoseconds()/c.res.Nanoseconds()),
		len(c.levels)-1,
	)
}
