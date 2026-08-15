package heligo

import "errors"

func (h *Helicorder) plotDataAverage(dataArr []PlotData) (float64, error) {
	if len(dataArr) == 0 {
		return 0, errors.New("no data to normalize")
	}

	// big.NewFloat uses the same 53-bit precision as float64 by default, so the
	// previous big.Float calculation added allocation cost without increasing
	// the precision of these float64 inputs and outputs.
	var sum float64
	for _, item := range dataArr {
		sum += item.Value
	}
	return sum / float64(len(dataArr)), nil
}
