package heligo

import (
	"bytes"
	"io"

	"gonum.org/v1/plot/vg"
)

func (h *Helicorder) Bytes(size int, format string) ([]byte, error) {
	var buf bytes.Buffer
	if err := h.WriteTo(size, format, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *Helicorder) WriteTo(size int, format string, dst io.Writer) error {
	writer, err := h.plotCtx.WriterTo(vg.Length(size), vg.Length(size), format)
	if err != nil {
		return err
	}

	_, err = writer.WriteTo(dst)
	return err
}
