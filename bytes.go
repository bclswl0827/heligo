package heligo

import (
	"bytes"

	"gonum.org/v1/plot/vg"
)

func (h *Helicorder) Bytes(size int, format string) ([]byte, error) {
	writer, err := h.plotCtx.WriterTo(vg.Length(size), vg.Length(size), format)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = writer.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
