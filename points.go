package heligo

import (
	"sort"
	"time"

	"gonum.org/v1/plot/plotter"
)

func (h *Helicorder) getPlotPoints(dataArr []PlotData, maxSamples int, scaleFactor, shiftVal float64) plotter.XYs {
	dataLength := len(dataArr)
	fillRatio := float64(dataLength) / float64(h.minutesTickSpan.Seconds()) / 100
	if fillRatio < 1 {
		maxSamples = int(fillRatio * float64(maxSamples))
	}

	// Perform downsampling with time alignment
	if dataLength > maxSamples {
		newDataArr := make([]PlotData, maxSamples)
		timeSpan := dataArr[dataLength-1].Time.Sub(dataArr[0].Time)

		// Interval for downsampled data
		sampleInterval := timeSpan / time.Duration(maxSamples-1)

		for i := 0; i < maxSamples; i++ {
			targetTime := dataArr[0].Time.Add(time.Duration(i) * sampleInterval)
			// Find the closest index where dataArr[j].Time is >= targetTime
			j := sort.Search(dataLength, func(k int) bool {
				return dataArr[k].Time.After(targetTime) || dataArr[k].Time.Equal(targetTime)
			})

			if j > 0 && j < dataLength {
				// Linear interpolation between dataArr[j-1] and dataArr[j]
				t1, t2 := dataArr[j-1].Time, dataArr[j].Time
				v1, v2 := dataArr[j-1].Value, dataArr[j].Value
				weight := targetTime.Sub(t1).Seconds() / t2.Sub(t1).Seconds()
				newDataArr[i] = PlotData{
					Time:  targetTime,
					Value: v1*(1-weight) + v2*weight,
				}
			} else if j == 0 {
				newDataArr[i] = dataArr[0]
			} else {
				newDataArr[i] = dataArr[dataLength-1]
			}
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
