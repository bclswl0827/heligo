package heligo

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"time"
)

func (h *Helicorder) Plot(date time.Time, maxSamples int, scaleFactor, lineWidth float64, colorScheme ColorScheme) error {
	if colorScheme == nil {
		colorScheme = &defaultColorScheme{}
	}

	if maxSamples < 100 {
		return errors.New("maxSamples must be greater than 100")
	}

	if lineWidth < 0.1 {
		return errors.New("lineWidth must be greater than 0.1")
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

	// Plot the data
	groupRows := int(time.Hour.Minutes() / h.minutesTickSpan.Minutes())
	totalRows := groupRows * int(h.hoursTickSpan.Hours())

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	// Set concurrency jobs depending on number of CPUs
	sem := make(chan struct{}, runtime.NumCPU()*10)

	for row := totalRows; row >= 1; row-- {
		currentCol := totalRows - row
		startTime := plotDate.Add(time.Duration(currentCol) * h.minutesTickSpan)
		endTime := startTime.Add(h.minutesTickSpan)

		// Get plot data from data provider
		plotData, err := h.dataProvider.GetPlotData(startTime, endTime)
		if err != nil {
			return err
		}
		if len(plotData) == 0 {
			continue
		}
		sort.SliceStable(plotData, func(i, j int) bool { return plotData[i].Time.Before(plotData[j].Time) })

		sem <- struct{}{}
		wg.Add(1)
		go func(row int, lineData []PlotData, currentCol int) {
			defer func() {
				<-sem
				defer wg.Done()
			}()

			lineColor := colorScheme.GetColor(groupRows, currentCol%groupRows)
			segments := h.getPlotSegments(lineData, maxSamples, row, scaleFactor, lineWidth, lineColor)

			for _, segment := range segments {
				mu.Lock()
				h.plotCtx.Add(segment)
				mu.Unlock()
			}
		}(row, plotData, currentCol)
	}

	wg.Wait()
	return nil
}
