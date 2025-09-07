package heligo

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gonum.org/v1/plot/plotter"
)

func (h *Helicorder) resampleMaxSamples(in []PlotData, maxSamples int) ([]PlotData, error) {
	if len(in) == 0 {
		return nil, errors.New("input slice is empty")
	}
	if maxSamples <= 0 {
		return nil, errors.New("maxSamples must be > 0")
	}

	n := len(in)
	if n <= maxSamples {
		return in, nil
	}

	out := make([]PlotData, maxSamples)
	ratio := float64(n-1) / float64(maxSamples-1)

	for i := 0; i < maxSamples; i++ {
		pos := float64(i) * ratio
		l := int(pos)
		if l >= n-1 {
			out[i] = in[n-1]
			continue
		}
		frac := pos - float64(l)

		timeDelta := in[l+1].Time.Sub(in[l].Time)
		newTime := in[l].Time.Add(time.Duration(frac * float64(timeDelta)))

		newValue := in[l].Value*(1-frac) + in[l+1].Value*frac

		out[i] = PlotData{
			Time:  newTime,
			Value: newValue,
		}
	}

	return out, nil
}

func (h *Helicorder) getPlotPoints(dataArr []PlotData, maxSamples, currentRow int, scaleFactor float64) (plotter.XYs, error) {
	dataArr, err := h.resampleMaxSamples(dataArr, maxSamples)
	if err != nil {
		return nil, fmt.Errorf("failed to resample data: %w", err)
	}

	// Normalize data to make it easier to plot
	normalizedDataArr, err := h.normalizePlotData(dataArr, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize data: %w", err)
	}
	scaleRatio := scaleFactor / math.MaxInt32

	minuteSteps := int(time.Hour.Minutes() / h.minutesTickSpan.Minutes())
	totalRows := minuteSteps * int(h.hoursTickSpan.Hours())
	currentCarry := (totalRows - currentRow) % minuteSteps

	var points plotter.XYs
	for idx := 0; idx < len(normalizedDataArr); idx++ {
		// Check carries to prevent overlapping lines
		calcCarry := int(normalizedDataArr[idx].Time.Minute()) / int(h.minutesTickSpan.Minutes())
		if calcCarry != currentCarry {
			continue
		}

		minutes := normalizedDataArr[idx].Time.Minute() - calcCarry*int(h.minutesTickSpan.Minutes())
		seconds := float64(normalizedDataArr[idx].Time.Second()) + float64(normalizedDataArr[idx].Time.Nanosecond())/1000000000
		points = append(points, plotter.XY{
			X: float64(minutes) + seconds/60,
			Y: float64(currentRow) + normalizedDataArr[idx].Value*scaleRatio,
		})
	}

	return points, nil
}
