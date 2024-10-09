package heligo

import "gonum.org/v1/plot/font"

func (h *Helicorder) Save(size int, filePath string) error {
	return h.plotCtx.Save(font.Length(size), font.Length(size), filePath)
}
