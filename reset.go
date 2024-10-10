package heligo

func (h *Helicorder) Reset() error {
	h.plotCtx = nil

	_, err := New(h.dataProvider, h.hoursTickSpan, h.minutesTickSpan)
	if err != nil {
		return err
	}

	return nil
}
