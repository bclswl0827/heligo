package heligo

import (
	"time"

	"gonum.org/v1/plot"
)

type PlotData struct {
	Time  time.Time
	Value float64
}

type DataProvider interface {
	GetPlotName() string
	GetStation() string
	GetNetwork() string
	GetChannel() string
	GetLocation() string
	GetPlotData(start, end time.Time) ([]PlotData, error)
}

type Helicorder struct {
	// Date of the plot in UTC
	// accurate to a specific day
	Date time.Time

	// hoursTickSpan is the number of hours
	// between ticks on the Y-axis
	hoursTickSpan time.Duration
	// minutesTickSpan is the number of minutes
	// between ticks on the X-axis
	minutesTickSpan time.Duration

	// Data provider interface
	// should implemented by user
	dataProvider DataProvider
	// Plot context used to render the plot
	plotCtx *plot.Plot
}

type hourTickMarker struct {
	hoursTickSpan   time.Duration
	minutesTickSpan time.Duration
}
