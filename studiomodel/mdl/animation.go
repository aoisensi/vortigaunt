package mdl

import (
	"fmt"
)

type AnimDesc struct {
	Header *AnimDescHeader
	Name   string
	Anim   *Anim
}

type AnimType byte

const (
	AnimTypeRawPos AnimType = 1 << iota
	AnimTypeRawRot
	AnimTypeAnimPos
	AnimTypeAnimRot
	AnimTypeDelta
	AnimTypeRawRot2
)

type Anim struct {
	Bone       byte
	AnimType   AnimType
	NextOffset int16
	Next       *Anim
	PosPtr     *AnimValuePtr
	RotPtr     *AnimValuePtr
	RawRot     *[4]float32
	RawPos     *[3]float32
}

type AnimValuePtr struct {
	Offsets [3]int16
	X, Y, Z []*int16
}

type AnimValue struct {
	Valid, Total byte
	Values       []int16
	Next         *AnimValue
}

type AnimDescHeader struct {
	Baseptr               int32
	NameOffset            int32
	FPS                   float32
	Flags                 uint32
	NumFrames             int32
	NumMovements          int32
	MovementOffset        int32
	_                     [6]int32
	AnimBlock             int32
	AnimOffset            int32
	NumIKRules            int32
	IKRuleOffset          int32
	AnimBlockIKRuleOffset int32
	NumLocalHierarchy     int32
	LocalHierarchyOffset  int32
	SectionOffset         int32
	SectionFrameCount     int32
	ZeroFrameSpan         int16
	ZeroFrameCount        int16
	ZeroFrameIndex        int32
	ZeroFramesTallTime    float32
}

func (d *Decoder) decodeAnimDesc(mdl *MDL) error {
	mdl.AnimDescs = make([]*AnimDesc, mdl.Header.LocalAnimCount)

	for i := range mdl.AnimDescs {
		header := new(AnimDescHeader)

		if err := d.read(header); err != nil {
			return fmt.Errorf("vortigaunt: decodeAnimDesc: %w", err)
		}
		ad := new(AnimDesc)
		ad.Header = header

		ad.Name = d.readName(ad.Header.NameOffset - 100)
		if err := d.ppush(
			header.AnimOffset-100,
			func() error {
				var err error
				ad.Anim, err = d.decodeAnim(int(header.NumFrames))
				return err
			},
		); err != nil {
			return err
		}
		mdl.AnimDescs[i] = ad
	}
	return nil
}

func (d *Decoder) decodeAnim(frames int) (*Anim, error) {
	anim := new(Anim)
	if err := d.read(&anim.Bone); err != nil {
		return nil, fmt.Errorf("vortigaunt: decodeAnim: %w", err)
	}
	if err := d.read(&anim.AnimType); err != nil {
		return nil, fmt.Errorf("vortigaunt: decodeAnim: %w", err)
	}
	if err := d.read(&anim.NextOffset); err != nil {
		return nil, fmt.Errorf("vortigaunt: decodeAnim: %w", err)
	}
	if anim.NextOffset > 0 {
		if err := d.ppush(
			int32(anim.NextOffset)-4,
			func() error {
				var err error
				anim.Next, err = d.decodeAnim(frames)
				return err
			},
		); err != nil {
			return nil, err
		}
	}
	raw := false
	var err error
	if anim.AnimType&AnimTypeRawRot > 0 {
		anim.RawRot, err = d.quat48()
		if err != nil {
			return nil, err
		}
		raw = true
	}
	if anim.AnimType&AnimTypeRawRot2 > 0 {
		anim.RawRot, err = d.quat64()
		if err != nil {
			return nil, err
		}
		raw = true
	}
	if anim.AnimType&AnimTypeRawPos > 0 {
		anim.RawPos, err = d.vec32()
		if err != nil {
			return nil, err
		}
		raw = true
	}
	if raw {
		return anim, nil
	}

	if anim.AnimType&AnimTypeAnimPos > 0 {
		anim.PosPtr, err = d.decodeAnimValuePtr(frames)
		if err != nil {
			return nil, err
		}
	}
	if anim.AnimType&AnimTypeAnimRot > 0 {
		anim.RotPtr, err = d.decodeAnimValuePtr(frames)
		if err != nil {
			return nil, err
		}
	}

	return anim, nil
}

func (d *Decoder) decodeAnimValuePtr(frames int) (*AnimValuePtr, error) {
	ptr := new(AnimValuePtr)
	if err := d.read(&ptr.Offsets); err != nil {
		return nil, fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
	}
	for offset, store := range map[int16]*[]*int16{
		ptr.Offsets[0]: &ptr.X,
		ptr.Offsets[1]: &ptr.Y,
		ptr.Offsets[2]: &ptr.Z,
	} {
		if offset <= 0 {
			continue
		}
		*store = make([]*int16, frames)
		if err := d.ppush(
			int32(offset)-6,
			func() error {
				i := 0
				for {
					vt := new(struct {
						Valid, Total byte
					})
					if err := d.read(vt); err != nil {
						return fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
					}
					if i == frames {
						return nil
					} else if i > frames {
						panic("vortigaunt: decodeAnimValuePtr: invalid frame count")
					}
					for range int(vt.Valid) {
						var v int16
						if err := d.read(&v); err != nil {
							return fmt.Errorf("vortigaunt: decodeAnimValuePtr: %w", err)
						}
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
