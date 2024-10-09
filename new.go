package heligo

import (
	"errors"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func New(dataProvider DataProvider, plotDate time.Time, hoursTickSpan, minutesTickSpan time.Duration) (Helicorder, error) {
	if minutesTickSpan.Minutes() == 0 {
		return Helicorder{}, errors.New("set minutesTickSpan field to zero is not allowed")
	}
	if minutesTickSpan.Minutes() > 60 {
		return Helicorder{}, errors.New("set minutesTickSpan field to greater than 60 minutes is not allowed")
	}
	if int(hoursTickSpan.Minutes())%int(minutesTickSpan.Minutes()) != 0 {
		return Helicorder{}, errors.New("only ticks of a multiple of minutesTickSpan field are allowed for hoursTickSpan field")
	}
	if hoursTickSpan.Hours() > 24 {
		return Helicorder{}, errors.New("set hoursTickSpan field to greater than 24 hours is not allowed")
	}

	heli := Helicorder{
		Date:            time.Date(plotDate.Year(), plotDate.Month(), plotDate.Day(), 0, 0, 0, 0, time.UTC),
		hoursTickSpan:   hoursTickSpan,
		minutesTickSpan: minutesTickSpan,
		dataProvider:    dataProvider,
		plotCtx:         plot.New(),
	}
	heli.plotCtx.Add(plotter.NewGrid())

	// Set title styles
	heli.plotCtx.Title.TextStyle.Font.Size = 20
	heli.plotCtx.Title.TextStyle.YAlign = -1.2
	heli.plotCtx.Title.TextStyle.Font.Variant = "Mono"

	// Set Y-axis styles
	heli.plotCtx.Y.Label.Text = "TIME (UTC)"
	heli.plotCtx.Y.Label.TextStyle.Font.Size = 16
	heli.plotCtx.Y.Label.TextStyle.YAlign = -0.2
	heli.plotCtx.Y.Label.Padding = -1
	heli.plotCtx.Y.Label.TextStyle.Font.Variant = "Mono"
	heli.plotCtx.Y.Min = 0
	heli.plotCtx.Y.Max = (time.Hour.Minutes() / minutesTickSpan.Minutes()) * hoursTickSpan.Hours()
	heli.plotCtx.Y.Tick.Marker = hourTickMarker{hoursTickSpan, minutesTickSpan}

	// Set X-axis styles
	heli.plotCtx.X.Label.Text = "TIME (MINUTES)"
	heli.plotCtx.X.Label.TextStyle.Font.Size = 16
	heli.plotCtx.X.Label.TextStyle.YAlign = -0.1
	heli.plotCtx.X.Label.Padding = 5
	heli.plotCtx.X.Label.TextStyle.Font.Variant = "Mono"
	heli.plotCtx.X.Min = 0
	heli.plotCtx.X.Max = minutesTickSpan.Minutes()

	return heli, nil
}
