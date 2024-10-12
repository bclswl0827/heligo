package main

import (
	"sort"
	"time"

	"github.com/bclswl0827/heligo"
)

func (d *dataProviderImpl) GetPlotName() string {
	return "AnyShake Observer"
}

func (d *dataProviderImpl) GetStation() string {
	return d.mseed.Series[0].FixedSection.StationCode
}

func (d *dataProviderImpl) GetNetwork() string {
	return d.mseed.Series[0].FixedSection.NetworkCode
}

func (d *dataProviderImpl) GetChannel() string {
	return d.mseed.Series[0].FixedSection.ChannelCode
}

func (d *dataProviderImpl) GetLocation() string {
	return d.mseed.Series[0].FixedSection.LocationCode
}

func (d *dataProviderImpl) GetPlotData(start, end time.Time) ([]heligo.PlotData, error) {
	var plotData []heligo.PlotData

	startIndex := sort.Search(len(d.mseed.Series), func(i int) bool {
		return !d.mseed.Series[i].FixedSection.StartTime.Before(start)
	})

	endIndex := sort.Search(len(d.mseed.Series), func(i int) bool {
		return d.mseed.Series[i].FixedSection.StartTime.After(end)
	})

	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex > len(d.mseed.Series) {
		endIndex = len(d.mseed.Series)
	}

	for _, v := range d.mseed.Series[startIndex:endIndex] {
		if start.Before(v.FixedSection.StartTime) && end.After(v.FixedSection.StartTime) {
			duration := v.FixedSection.SamplesNumber / v.FixedSection.SampleFactor
			sampleRate := int(v.FixedSection.SamplesNumber / duration)
			for i, vv := range v.DataSection.Decoded {
				plotData = append(plotData, heligo.PlotData{
					Time:  v.FixedSection.StartTime.Add(time.Duration(i*1000/sampleRate) * time.Millisecond),
					Value: float64(vv.(int32)),
				})
			}
		}
	}

	return plotData, nil
}
