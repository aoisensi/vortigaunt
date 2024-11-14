package mdl

import (
	"fmt"
)

var errFailedToReadAnimValueForBadFrames = fmt.Errorf("vortigaunt: failed to read anim value for bad frames")

type AnimValuePtr struct {
	Offsets [3]int16
	X, Y, Z []*int16
}

type AnimValue struct {
	Valid, Total byte
	Values       []int16
	Next         *AnimValue
}

func (d *Decoder) decodeAnimValuePtr(frames int) (*AnimValuePtr, error) {
	ptr := new(AnimValuePtr)
	if err := d.read(&ptr.Offsets); err != nil {
		return nil, fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
	}
	for in := range 3 {
		offset := ptr.Offsets[in]
		store := []*[]*int16{&ptr.X, &ptr.Y, &ptr.Z}[in]
		if offset <= 0 {
			continue
		}
		*store = make([]*int16, frames)
		if err := d.ppush(
			int32(offset)-6,
			func() error {
				i := 0
				for {
					if i == frames {
						return nil
					} else if i > frames {
						return errFailedToReadAnimValueForBadFrames
						// return nil
					}
					vt := new(struct{ Valid, Total byte })
					if err := d.read(vt); err != nil {
						return fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
					}
					if vt.Total == 0 {
						return errFailedToReadAnimValueForBadFrames
					}
					for range int(vt.Valid) {
						var v int16
						if err := d.read(&v); err != nil {
							return fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
						}
						*store = append(*store, &v)
						(*store)[i] = &v
						i++
					}
					i += int(vt.Total - vt.Valid)
				}
			},
		); err != nil {
			return nil, fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
		}

	}
	return ptr, nil
}
