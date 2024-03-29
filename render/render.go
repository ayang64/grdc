package render

import (
	"image"
	"image/color"

	"golang.org/x/crypto/ssh/terminal"
)

func TextToBitmap(s string) *image.RGBA {
	// figure out the number of runes in s
	runes := []rune(s)
	l := len(runes)

	// each glyph is 12x22 so we'll need (12 * l)x22 bitmap
	img := image.NewRGBA(image.Rect(0, 0, l*12, 22))

	for i, r := range runes {
		y := 0
		for _, v := range font[r] {
			for cur, end := 16, 3; cur > end; cur-- {
				color := func() color.RGBA {
					if (v>>cur)&0x1 == 1 {
						return color.RGBA{
							R: 0xff,
							G: 0xff,
							B: 0xff,
							A: 0xff,
						}
					}
					return color.RGBA{
						R: 0x00,
						G: 0x00,
						B: 0x00,
						A: 0xff,
					}
				}

				x := (16 - cur) + (12 * i)
				img.Set(x, y, color())
			}
			y++
		}
	}

	termWidth, termHeight, _ := terminal.GetSize(0)
	termHeight *= 2

	out := image.NewRGBA(image.Rect(0, 0, termWidth, termHeight))

	outBounds := out.Bounds()

	textBounds := img.Bounds()

	// this type is used purely for scaling.  i'm trying an averaging algorithm
	// but i think it is either overkill or doesn't work.  i didn't look this up
	// -- it is just something stupid i came up with -- do not reuse.
	type point struct {
		r     uint32
		g     uint32
		b     uint32
		a     uint32
		count uint32
	}

	bitmap := make([]point, termWidth*termHeight)

	// calculate scaling factors.
	xf := float64(textBounds.Max.X-textBounds.Min.X) / float64(outBounds.Max.X-outBounds.Min.X)
	yf := float64(textBounds.Max.Y-textBounds.Min.Y) / float64(outBounds.Max.Y-outBounds.Min.Y)

	// simple translate convenience functon.
	xlate := func(x int, y int) (int, int) {
		return int(float64(x) * xf), int(float64(y) * yf)
	}

	maxcount := uint32(0)
	mincount := ^uint32(0)
	for stride, y := termWidth, outBounds.Min.Y; y < outBounds.Max.Y; y++ {
		for x := outBounds.Min.X; x < outBounds.Max.X; x++ {
			ox, oy := xlate(x, y)
			c := img.At(ox, oy)
			r, g, b, a := c.RGBA()
			offs := x + (stride)*y
			bitmap[offs].r += r
			bitmap[offs].g += g
			bitmap[offs].b += b
			bitmap[offs].a += a
			bitmap[offs].count++

			if bitmap[offs].count > maxcount {
				maxcount = bitmap[offs].count
			}
			if bitmap[offs].count < mincount {
				mincount = bitmap[offs].count
			}
		}
	}

	avg := func(sum uint32, count uint32) uint8 {
		if count == 0 {
			return 0
		}
		return uint8(sum / count)
	}

	for i := range bitmap {
		newColor := color.RGBA{
			R: avg(bitmap[i].r, maxcount),
			G: 0, // avg(bitmap[i].g, bitmap[i].count),
			B: 0, // avg(bitmap[i].b, bitmap[i].count),
			A: avg(bitmap[i].a, maxcount),
		}

		x := i % termWidth
		y := i / termWidth

		out.Set(x, y, newColor)
	}

	return out
}
