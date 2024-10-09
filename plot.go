package heligo

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (h *Helicorder) Plot(sampleRate, downSampleFactor int, scaleFactor, lineWidth float64) error {
	if sampleRate < 1 {
		return errors.New("sampleRate must be greater than 0")
	}
	if downSampleFactor < 1 {
		return errors.New("downSampleFactor must be greater than 0")
	}

	// Set plot title
	h.plotCtx.Title.Text = fmt.Sprintf(
		"%s / %s / %s %s %s %s",
		h.dataProvider.GetPlotName(),
		h.Date.UTC().Format("Jan 2, 2006"),
		h.dataProvider.GetStation(),
		h.dataProvider.GetChannel(),
		h.dataProvider.GetNetwork(),
		h.dataProvider.GetLocation(),
	)

	// Get plot data from data provider
	plotData, err := h.dataProvider.GetPlotData(h.Date, h.Date.Add(h.hoursTickSpan))
	if err != nil {
		return err
	}
	if len(plotData) == 0 {
		return errors.New("no data found")
	}
	sort.SliceStable(plotData, func(i, j int) bool { return plotData[i].Time.Before(plotData[j].Time) })

	// Plot the data
	groupRows := int(time.Hour.Minutes() / h.minutesTickSpan.Minutes())
	totalRows := groupRows * int(h.hoursTickSpan.Hours())

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	for row := totalRows; row >= 1; row-- {
		currentCol := totalRows - row
		startTime := h.Date.Add(time.Duration(currentCol) * h.minutesTickSpan)
		endTime := startTime.Add(h.minutesTickSpan)

		// Get slice within time range
		startIndex := sort.Search(len(plotData), func(i int) bool { return plotData[i].Time.After(startTime) || plotData[i].Time.Equal(startTime) })
		endIndex := sort.Search(len(plotData), func(i int) bool { return plotData[i].Time.After(endTime) })
		if startIndex == endIndex {
			continue
		}
		lineData := plotData[startIndex:endIndex]

		// Get final plot points concurrently
		wg.Add(1)
		go func(row int, lineData []PlotData, currentCol int) {
			points := h.getPlotPoints(lineData, sampleRate, downSampleFactor, scaleFactor, float64(row))
			line, _, _ := plotter.NewLinePoints(points)
			line.LineStyle.Width = vg.Length(lineWidth)
			line.Color = h.getColor(groupRows, currentCol%groupRows)

			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()
			h.plotCtx.Add(line)
		}(row, lineData, currentCol)
	}

	wg.Wait()
	return nil
}
