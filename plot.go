package heligo

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (h *Helicorder) Plot(date time.Time, maxSamples int, scaleFactor, lineWidth float64) error {
	if maxSamples < 1 {
		return errors.New("maxSamples must be greater than 0")
	}
	plotDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	// Set plot title
	h.plotCtx.Title.Text = fmt.Sprintf(
		"%s / %s / %s %s %s %s",
		h.dataProvider.GetPlotName(),
		plotDate.UTC().Format("Jan 2, 2006"),
		h.dataProvider.GetStation(),
		h.dataProvider.GetChannel(),
		h.dataProvider.GetNetwork(),
		h.dataProvider.GetLocation(),
	)

	// Get plot data from data provider
	plotData, err := h.dataProvider.GetPlotData(plotDate, plotDate.Add(h.hoursTickSpan))
	if err != nil {
		return err
	}
	if len(plotData) == 0 {
		return errors.New("no data found for plot")
	}
	sort.SliceStable(plotData, func(i, j int) bool { return plotData[i].Time.Before(plotData[j].Time) })

	// Plot the data
	groupRows := int(time.Hour.Minutes() / h.minutesTickSpan.Minutes())
	totalRows := groupRows * int(h.hoursTickSpan.Hours())

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	numCPU := runtime.NumCPU()
	sem := make(chan struct{}, numCPU)

	for row := totalRows; row >= 1; row-- {
		currentCol := totalRows - row
		startTime := plotDate.Add(time.Duration(currentCol) * h.minutesTickSpan)
		endTime := startTime.Add(h.minutesTickSpan)

		// Get slice within time range
		startIndex := sort.Search(len(plotData), func(i int) bool { return plotData[i].Time.After(startTime) || plotData[i].Time.Equal(startTime) })
		endIndex := sort.Search(len(plotData), func(i int) bool { return plotData[i].Time.After(endTime) })
		if startIndex == endIndex {
			continue
		}
		lineData := plotData[startIndex:endIndex]

		sem <- struct{}{}
		wg.Add(1)
		go func(row int, lineData []PlotData, currentCol int) {
			defer func() {
				<-sem
				defer wg.Done()
			}()

			points := h.getPlotPoints(lineData, maxSamples, scaleFactor, float64(row))
			line, _, _ := plotter.NewLinePoints(points)
			line.LineStyle.Width = vg.Length(lineWidth)
			line.Color = h.getColor(groupRows, currentCol%groupRows)

			mu.Lock()
			h.plotCtx.Add(line)
			mu.Unlock()
		}(row, lineData, currentCol)
	}

	wg.Wait()
	return nil
}
