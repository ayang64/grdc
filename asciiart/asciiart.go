package asciiart

import (
	"bytes"
	"fmt"
	"image"
	"io"
)

func Encode(w io.Writer, img image.Image) {
	rect := img.Bounds()
	// minor optimization -- store the previous color and avoid emitting escape
	// code if the color hasn't changed.
	prevTop, prevBottom := [3]uint32{0, 0, 0}, [3]uint32{0, 0, 0}

	buf := &bytes.Buffer{}
	buf.Write([]byte("\x1b[;f"))

	for y := rect.Min.Y; y < rect.Max.Y; y += 2 {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, _ := col.RGBA()

			curTop := [3]uint32{r >> 8, g >> 8, b >> 8}

			if y == 0 || curTop != prevTop {
				buf.Write([]byte(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)))
				prevTop = curTop
			}

			col = img.At(x, y+1)
			r, g, b, _ = col.RGBA()
			curBottom := [3]uint32{r >> 8, g >> 8, b >> 8}

			if y == 0 || curBottom != prevBottom {
				buf.Write([]byte(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)))
				prevBottom = curBottom
			}
			buf.WriteRune('▀')
		}
	}

	buf.Write([]byte("\x1b[48;2;0;0;0m"))
	io.Copy(w, buf)
}
