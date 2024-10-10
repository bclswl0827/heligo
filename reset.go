package heligo

func (h *Helicorder) Reset() error {
	_, err := New(h.dataProvider, h.hoursTickSpan, h.minutesTickSpan)
	if err != nil {
		return err
	}

	return nil
}
