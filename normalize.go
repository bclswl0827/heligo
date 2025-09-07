package heligo

import (
	"errors"
	"math/big"
)

func (h *Helicorder) normalizePlotData(dataArr []PlotData, center float64) ([]PlotData, error) {
	if len(dataArr) == 0 {
		return nil, errors.New("no data to normalize")
	}

	var sum big.Float
	for _, val := range dataArr {
		bigVal := big.NewFloat(val.Value)
		sum.Add(&sum, bigVal)
	}

	avg := new(big.Float).Quo(&sum, big.NewFloat(float64(len(dataArr))))

	minVal, maxVal := 0.0, 0.0
	normalizedData := make([]PlotData, len(dataArr))
	for i, item := range dataArr {
		bigItem := big.NewFloat(item.Value)
		bigCenter := big.NewFloat(center)

		normalizedItem := new(big.Float).Sub(bigItem, avg)
		normalizedItem.Add(normalizedItem, bigCenter)

		data, _ := normalizedItem.Float64()
		normalizedData[i].Value = data
		normalizedData[i].Time = item.Time

		if data < minVal {
			minVal = data
		}
		if data > maxVal {
			maxVal = data
		}
	}

	return normalizedData, nil
}
