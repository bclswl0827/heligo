package heligo

import (
	"time"

	"gonum.org/v1/plot"
)

func (c hourTickMarker) Ticks(min, max float64) []plot.Tick {
	intMin, intMax := int(min), int(max)
	steps := int(time.Hour.Minutes() / c.minutesTickSpan.Minutes())
	ticks := make([]plot.Tick, (intMax-intMin)/steps+1)

	for tick, idx := intMin, 0; tick <= intMax; tick += steps {
		hourValue := int(c.hoursTickSpan.Hours()) - tick/steps
		timeStr := time.Date(0, 0, 0, hourValue, 0, 0, 0, time.UTC).Format("15:04")
		if tick == intMin {
			timeStr = "0"
		}

		ticks[idx] = plot.Tick{Value: float64(tick), Label: timeStr}
		idx++
	}

	return ticks
}
