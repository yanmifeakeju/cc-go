package wc

import (
	"io"
	"unicode"
	"unicode/utf8"
)

type WCopts struct {
	CountBytes bool
	CountLines bool
	CountWords bool
	CountChars bool
}

type WC struct {
	Reader io.Reader
}

type Counts struct {
	Bytes      int64
	Words      int
	Lines      int
	Characters int
}

func (wc *WC) Count(opts WCopts) (Counts, error) {
	counts, err := wc.countAll()
	if err != nil {
		return Counts{}, err
	}

	return counts, nil
}

func (wc *WC) countAll() (Counts, error) {
	var counts Counts
	var inWord bool

	buf := make([]byte, 32*1024)

	for {
		n, err := wc.Reader.Read(buf)
		if err != nil && err != io.EOF {
			return counts, err
		}

		counts.Bytes += int64(n)

		// Process the buffer
		p := buf[:n]
		for len(p) > 0 {
			r, size := utf8.DecodeRune(p)
			counts.Characters++

			if r == '\n' {
				counts.Lines++
			}

			if unicode.IsSpace(r) {
				inWord = false
			} else {
				if !inWord {
					counts.Words++
					inWord = true
				}
			}

			p = p[size:]
		}

		if err == io.EOF {
			break
		}
	}

	return counts, nil
}
