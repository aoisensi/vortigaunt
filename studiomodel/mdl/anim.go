package mdl

import "fmt"

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

type AnimType byte

const (
	AnimTypeRawPos AnimType = 1 << iota
	AnimTypeRawRot
	AnimTypeAnimPos
	AnimTypeAnimRot
	AnimTypeDelta
	AnimTypeRawRot2
)

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
	var err error
	// fmt.Println("AnimType:", anim.AnimType)
	if anim.AnimType&AnimTypeRawRot2 > 0 {
		anim.RawRot, err = d.quat64()
		if err != nil {
			return nil, err
		}
	} else if anim.AnimType&AnimTypeRawRot > 0 {
		anim.RawRot, err = d.quat48()
		if err != nil {
			return nil, err
		}
	}
	if anim.AnimType&AnimTypeRawPos > 0 {
		anim.RawPos, err = d.vec32()
		if err != nil {
			return nil, err
		}
	}
	if anim.AnimType&(AnimTypeRawPos&AnimTypeRawRot&AnimTypeRawRot2) > 0 {
		return anim, nil
	}

	if anim.AnimType&AnimTypeAnimRot > 0 {
		anim.RotPtr, err = d.decodeAnimValuePtr(frames)
		if err != nil {
			return nil, err
		}
	}
	if anim.AnimType&AnimTypeAnimPos > 0 {
		anim.PosPtr, err = d.decodeAnimValuePtr(frames)
		if err != nil {
			return nil, err
		}
	}

	return anim, nil
}
