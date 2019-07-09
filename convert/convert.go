package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Scanner struct {
	r   io.RuneScanner // rune Scanner
	w   io.Writer      //
	hex int
}

func NewScanner(r io.Reader) *Scanner {
	runeReader := func() io.RuneScanner {
		switch v := r.(type) {
		case io.RuneScanner:
			return v
		default:
			return bufio.NewReader(r)
		}
	}

	rc := Scanner{
		r: runeReader(),
		w: os.Stdout,
	}

	return &rc
}

func (s *Scanner) Comment() error {
	b := strings.Builder{}

	for {
		r, _, err := s.r.ReadRune()
		if err != nil {
			return err
		}

		b.WriteRune(r)

		if v := b.String(); strings.HasSuffix(v, "*/") {
			return nil
		}
	}
}

func (s *Scanner) Hex() (uint8, error) {
	b := strings.Builder{}

	for {
		r, _, err := s.r.ReadRune()
		if err != nil {
			return 0, err
		}

		if !(r == 'x' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			s.r.UnreadRune()
			break
		}

		b.WriteRune(r)
	}

	var rc uint8

	fmt.Sscanf(b.String()[2:], "%x", &rc)
	return rc, nil
}

func (s *Scanner) Peek() (rune, error) {
	r, _, err := s.r.ReadRune()

	s.r.UnreadRune()

	if err != nil {
		return r, err
	}

	return r, nil
}

type bitmap struct {
	bits []uint8
}

func (b *bitmap) dump() {
	fmt.Printf("{\n")
	for bits, stride := b.bits, 2; len(bits) > 0; bits = bits[stride:] {
		fmt.Printf("0x%04[1]x, // %016[1]b\n", uint16(bits[0])<<8|uint16(bits[1]))
	}
	fmt.Printf("},\n")
}

func main() {
	scn := NewScanner(os.Stdin)

	data := []uint8{}

	for {
		r, err := scn.Peek()

		if err != nil {
			break
		}

		switch r {
		case '/':
			scn.Comment()
		case '0':
			scn.hex += (scn.hex + 1) % 44
			h, _ := scn.Hex()
			data = append(data, h)
		default:
			scn.r.ReadRune()
		}
	}

	fmt.Printf("font = [][]uint16{\n")
	for pos, stride := 0, 44; pos < len(data); pos += stride {
		cur := data[pos : pos+stride]
		b := bitmap{
			bits: cur,
		}
		b.dump()
	}
	fmt.Printf("}\n")
}
