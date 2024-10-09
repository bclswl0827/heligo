package heligo

import (
	"gonum.org/v1/plot/plotter"
)

func (h *Helicorder) getPlotPoints(dataArr []PlotData, sampleRate, downSampleFactor int, scaleFactor, shiftVal float64) plotter.XYs {
	dataLength := len(dataArr)

	// Get downsampled data
	maxSamples := int(float64(dataLength) / float64(sampleRate*int(h.minutesTickSpan.Seconds())) * float64(downSampleFactor))
	if dataLength > maxSamples {
		newDataArr := make([]PlotData, maxSamples)
		for i := range newDataArr {
			newDataArr[i] = dataArr[int(float64(dataLength)/float64(maxSamples))*i]
		}

		dataArr = newDataArr
		dataLength = maxSamples
	}

	normalizedDataArr, minVal, maxVal := h.normalizePlotData(dataArr, 0)
	scaleRatio := h.getScaleRatio(normalizedDataArr, scaleFactor, minVal, maxVal)
	points := make(plotter.XYs, dataLength)
	for i := 0; i < dataLength; i++ {
		minutes := normalizedDataArr[i].Time.Minute() - (normalizedDataArr[i].Time.Minute()/int(h.minutesTickSpan.Minutes()))*int(h.minutesTickSpan.Minutes())
		seconds := float64(normalizedDataArr[i].Time.Second()) + float64(normalizedDataArr[i].Time.Nanosecond())/1000000000
		points[i].X = float64(minutes) + seconds/60
		points[i].Y = shiftVal + normalizedDataArr[i].Value*scaleRatio
	}

	return points
}
