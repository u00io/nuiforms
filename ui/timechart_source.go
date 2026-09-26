package ui

import (
	"math"
	"sort"
	"sync"
	"time"
)

// TimeChartMemorySource is a TimeChartDataSource over an in-memory slice of
// points sorted by DT. It is safe to add points from another goroutine.
type TimeChartMemorySource struct {
	mtx    sync.RWMutex
	points []TimeChartPoint
}

func NewTimeChartMemorySource() *TimeChartMemorySource {
	return &TimeChartMemorySource{}
}

// SetPoints replaces the data. The points must be sorted by DT.
func (c *TimeChartMemorySource) SetPoints(points []TimeChartPoint) {
	c.mtx.Lock()
	c.points = points
	c.mtx.Unlock()
}

// AddPoint appends a point. Its DT must not be earlier than the last one.
func (c *TimeChartMemorySource) AddPoint(p TimeChartPoint) {
	c.mtx.Lock()
	c.points = append(c.points, p)
	c.mtx.Unlock()
}

func (c *TimeChartMemorySource) GetData(from, to time.Time, groupDuration time.Duration) []TimeChartPoint {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return TimeChartDownsample(c.points, from, to, groupDuration)
}

// TimeChartDownsample returns the points of the sorted slice that fall into
// [from, to], plus the nearest point on each side, merged into buckets of
// groupDuration: First/Last of a bucket come from its first and last good
// point, High/Low are the extremes of the good points, HasGood is set when
// the bucket has any good point and HasBad when it has any bad one - so a
// short outage stays visible however far the chart is zoomed out. Buckets
// are aligned to multiples of groupDuration since the Unix epoch, so panning
// does not change how points are grouped. It is the same aggregation a server should do before sending
// data to the chart.
func TimeChartDownsample(points []TimeChartPoint, from, to time.Time, groupDuration time.Duration) []TimeChartPoint {
	i0 := sort.Search(len(points), func(i int) bool { return !points[i].DT.Before(from) })
	i1 := sort.Search(len(points), func(i int) bool { return points[i].DT.After(to) })
	if i0 > 0 {
		i0--
	}
	if i1 < len(points) {
		i1++
	}
	src := points[i0:i1]

	if groupDuration <= 0 {
		res := make([]TimeChartPoint, len(src))
		copy(res, src)
		return res
	}

	g := int64(groupDuration)
	res := make([]TimeChartPoint, 0, min(len(src), int(to.Sub(from)/groupDuration)+3))
	var bucket int64
	for _, p := range src {
		ns := p.DT.UnixNano()
		b := ns / g
		if ns < 0 && ns%g != 0 {
			b--
		}
		good := p.HasValue()
		if len(res) > 0 && b == bucket {
			last := &res[len(res)-1]
			last.HasBad = last.HasBad || p.HasBad
			switch {
			case !good:
			case !last.HasGood:
				last.First, last.Last, last.High, last.Low = p.First, p.Last, p.High, p.Low
				last.HasGood = true
			default:
				last.Last = p.Last
				last.High = math.Max(last.High, p.High)
				last.Low = math.Min(last.Low, p.Low)
			}
			continue
		}
		bucket = b
		q := TimeChartPoint{DT: time.Unix(0, b*g), HasBad: p.HasBad}
		if good {
			q.First, q.Last, q.High, q.Low = p.First, p.Last, p.High, p.Low
			q.HasGood = true
		}
		res = append(res, q)
	}
	return res
}
