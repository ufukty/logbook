package lec

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

type Infinite struct {
	start  time.Time
	res    time.Duration
	levels [][]int
}

func NewInfinite(start time.Time, resolution time.Duration) *Infinite {
	return &Infinite{
		start:  start,
		res:    resolution,
		levels: [][]int{},
	}
}

func (c *Infinite) grow(to time.Time) {
	neededlevels := minpow(int(to.Sub(c.start)/c.res)+1) + 1
	missinglevels := neededlevels - len(c.levels)
	if missinglevels > 0 {
		c.levels = append(c.levels, make([][]int, missinglevels)...)
	}

	neededcells := pow(2, neededlevels)
	for l := 0; l < len(c.levels); l++ {
		missingcells := (neededcells+1)/2 - len(c.levels[l])
		if missingcells > 0 {
			c.levels[l] = append(c.levels[l], make([]int, missingcells)...)
			if l > 0 {
				for j := len(c.levels[l]) - missingcells; j < len(c.levels[l]); j++ {
					if j%4 == 0 {
						c.levels[l][j] = c.levels[l-1][j*2] // no virtual child cell
					}
				}
			}
		}
		neededcells /= 2
	}
}

func (c *Infinite) Save(t time.Time, v int) {
	c.grow(t)
	cell := int(float64(t.Sub(c.start).Nanoseconds()) / float64(c.res)) // virtual index in first layer
	for level := 0; level < len(c.levels); level++ {
		if cell%2 == 0 { // stored cell
			c.levels[level][cell/2] += v
		}
		cell /= 2 // floor
	}
}

func (c *Infinite) cellvalue(level, cell int) int {
	if len(c.levels) <= level {
		return 0
	}
	if cell%2 == 0 {
		return c.levels[level][cell/2] // virtual index -> array index
	} else if level != len(c.levels)-1 {
		parent := c.cellvalue(level+1, cell/2) // floor
		sibling := c.cellvalue(level, cell-1)
		return parent - sibling
	}
	return 0
}

func (c *Infinite) query(from, to, level int) (int, error) {
	if from > to {
		return 0, fmt.Errorf("'from' must be less than or equal to 'to'")
	}
	total := 0
	res := pow(2, level)
	l := int(math.Ceil(float64(from) / float64(res)))
	r := int(math.Floor(float64(to) / float64(res)))
	if r < l {
		if level == 0 {
			return -1, fmt.Errorf("bad developer") // FIXME:
		}
		q, err := c.query(from, to, level-1)
		if err != nil {
			return -1, fmt.Errorf("c.query on lower level: %w", err)
		}
		return q, nil
	}
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
	if from == to {
		return -1, fmt.Errorf("values 'from' and 'to' are same")
	}
	if to.Before(c.start) {
		return -1, fmt.Errorf("ends before series start: %s", c.start)
	}
	if from.Before(c.start) {
		return -1, fmt.Errorf("starts before the series: %s", c.start)
	}
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

func (c *Infinite) dump() string {
	b := bytes.NewBuffer([]byte{})
	g := -1
	for i := 0; i < len(c.levels); i++ {
		lastci := len(c.levels[i]) - 1
		lastreal := pow(2, i) * (2*lastci + 1)
		if g < lastreal {
			g = lastreal
		}
	}
	d := twodim(len(c.levels), g+1, -1)

	for level := len(c.levels) - 1; level >= 0; level-- {
		for ci := 0; ci < len(c.levels[level]); ci++ {
			d[level][(2*ci)*pow(2, level)] = c.cellvalue(level, ci*2)
			if level != len(c.levels)-1 {
				d[level][(2*ci+1)*pow(2, level)] = c.cellvalue(level, ci*2+1)
			}
		}
	}
	d = transpose(d)
	for i := 0; i < len(d); i++ {
		if i != 0 {
			fmt.Fprintln(b)
		}
		for j := len(d[i]) - 1; j >= 0; j-- {
			if d[i][j] != -1 {
				fmt.Fprintf(b, "%-5d", d[i][j])
			} else {
				fmt.Fprintf(b, "     ")
			}
		}
	}
	return b.String()
}
