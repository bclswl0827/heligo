package heligo

func (h *Helicorder) getScaleRatio(arr []PlotData, targetDiff, minVal, maxVal float64) float64 {
	if len(arr) == 0 {
		return 1
	}

	diff := maxVal - minVal
	if diff == 0 {
		return 1
	}

	return targetDiff / diff
}
