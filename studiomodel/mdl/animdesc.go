package mdl

import (
	"errors"
	"fmt"
	"log"
)

// mstudioanimdesc_t
type AnimDesc struct {
	Header       *AnimDescHeader
	Name         string
	Anims        []*Anim
	AnimSections []*AnimSection
}

type AnimDescFlags uint32

var (
	AnimDescFlagFrameAnim = AnimDescFlags(0x40)
)

type AnimDescHeader struct {
	Baseptr               int32
	NameOffset            int32
	FPS                   float32
	Flags                 AnimDescFlags
	NumFrames             int32 // total animation frames count
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
} // 100 bytes

func (d *Decoder) decodeAnimDesc(mdl *MDL) error {
	mdl.AnimDescs = make([]*AnimDesc, 0, mdl.Header.LocalAnimCount)

	for range mdl.Header.LocalAnimCount {
		header := new(AnimDescHeader)

		if err := d.read(header); err != nil {
			return fmt.Errorf("vortigaunt: decodeAnimDesc: %w", err)
		}
		ad := new(AnimDesc)
		ad.Header = header

		ad.Name = d.readName(ad.Header.NameOffset - 100)
		if header.SectionFrameCount == 0 {
			err := d.ppush(
				header.AnimOffset-100,
				func() error {
					anim, err := d.decodeAnim(int(header.NumFrames))
					ad.Anims = []*Anim{anim}
					return err
				},
			)
			if errors.Is(err, errFailedToReadAnimValueForBadFrames) {
				log.Println("decodeAnimDesc: skipping anim", ad.Name, "due to bad frames")
				continue
			}
			if err != nil {
				return fmt.Errorf("vortigaunt: decodeAnimDesc: %w", err)
			}
		} else {
			sections := int(header.NumFrames/header.SectionFrameCount) + 2 // https://github.com/ValveSoftware/source-sdk-2013/blob/master/mp/src/public/studio.cpp#L111
			if err := d.ppush(
				header.SectionOffset-100,
				func() error {
					ad.AnimSections = make([]*AnimSection, sections)
					for i := range ad.AnimSections {
						section := new(AnimSection)
						if err := d.read(section); err != nil {
							return fmt.Errorf("vortigaunt: decodeAnimDesc: %w", err)
						}
						ad.AnimSections[i] = section
					}
					return nil
				},
			); err != nil {
				return fmt.Errorf("vortigaunt: decodeAnimDesc: %w", err)
			}

			ad.Anims = make([]*Anim, 0, sections)
			for si, as := range ad.AnimSections {
				framesCount := int(header.SectionFrameCount) + 1
				if si >= sections-2 {
					framesCount = int(header.NumFrames) - (sections-2)*int(header.SectionFrameCount)
				}
				err := d.ppush(
					as.AnimOffset-100, // +header.SectionOffset-ad.AnimSections[0].AnimOffset,
					func() error {
						anim, err := d.decodeAnim(framesCount)
						ad.Anims = append(ad.Anims, anim)
						return err
					},
				)
				if errors.Is(err, errFailedToReadAnimValueForBadFrames) {
					log.Println("decodeAnimDesc: skipping anim", ad.Name, "due to bad frames")
					continue
				}
				if err != nil {
					return fmt.Errorf("vortigaunt: decodeAnimDesc: %w", err)
				}
			}
		}
		mdl.AnimDescs = append(mdl.AnimDescs, ad)
	}
	return nil
}
