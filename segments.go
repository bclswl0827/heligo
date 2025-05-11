package heligo

import (
	"image/color"
	"time"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (h *Helicorder) getPlotSegments(dataArr []PlotData, maxSamples, currentRow int, scaleFactor, lineWidth float64, lineColor color.Color) []*plotter.Line {
	var (
		lines    []*plotter.Line
		lineData = make([]PlotData, 0, len(dataArr))
	)

	// Helper function to create and append a line segment
	appendLine := func() {
		if len(lineData) > 0 {
			line := &plotter.Line{
				XYs:       h.getPlotPoints(lineData, maxSamples, currentRow, scaleFactor),
				LineStyle: plotter.DefaultLineStyle,
			}
			line.LineStyle.Width = vg.Length(lineWidth)
			line.LineStyle.Color = lineColor
			lines = append(lines, line)
			lineData = []PlotData{}
		}
	}

	for i := 1; i < len(dataArr); i++ {
		lineData = append(lineData, dataArr[i-1])

		// Create new segment if the time difference is greater than 2 second
		if dataArr[i].Time.Sub(dataArr[i-1].Time) > 2*time.Second {
			appendLine()
		}
	}

	lineData = append(lineData, dataArr[len(dataArr)-1])
	appendLine()

	return lines
}
