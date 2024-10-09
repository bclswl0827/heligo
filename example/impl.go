package main

import (
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
	for _, v := range d.mseed.Series {
		if start.Before(v.FixedSection.StartTime) && end.After(v.FixedSection.StartTime) {
			duration := v.FixedSection.SamplesNumber / v.FixedSection.SampleFactor
			sampleRate := float64(v.FixedSection.SamplesNumber) / float64(duration)
			for i, vv := range v.DataSection.Decoded {
				plotData = append(plotData, heligo.PlotData{
					Time:  v.FixedSection.StartTime.Add(time.Duration(i*int(1000/sampleRate)) * time.Millisecond),
					Value: float64(vv.(int32)),
				})
			}
		}
	}

	return plotData, nil
}
