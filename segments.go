package heligo

import (
	"fmt"
	"image/color"
	"time"

	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func (h *Helicorder) getPlotSegments(dataArr []PlotData, maxSamples, currentRow int, scaleFactor, lineWidth float64, lineColor color.Color) ([]*plotter.Line, error) {
	var lines []*plotter.Line

	// Helper function to create and append a line segment
	appendLine := func(lineData []PlotData) error {
		if len(lineData) > 0 {
			xyPoints, err := h.getPlotPoints(lineData, maxSamples, currentRow, scaleFactor)
			if err != nil {
				return fmt.Errorf("failed to get plot points: %w", err)
			}
			line := &plotter.Line{
				XYs:       xyPoints,
				LineStyle: plotter.DefaultLineStyle,
			}
			line.LineStyle.Width = vg.Length(lineWidth)
			line.LineStyle.Color = lineColor
			lines = append(lines, line)
		}
		return nil
	}

	segmentStart := 0
	for i := 1; i < len(dataArr); i++ {
		// Create new segment if the time difference is greater than 1 second
		if dataArr[i].Time.Sub(dataArr[i-1].Time) >= time.Second {
			if err := appendLine(dataArr[segmentStart:i]); err != nil {
				return nil, err
			}
			segmentStart = i
		}
	}

	if err := appendLine(dataArr[segmentStart:]); err != nil {
		return nil, err
	}

	return lines, nil
}
