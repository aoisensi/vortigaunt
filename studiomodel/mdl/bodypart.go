package mdl

// mstudiobodyparts_t
type BodyPart struct {
	Header *BodyPartHeader
	Name   string
	Models []*Model
}

type BodyPartHeader struct {
	NameIndex  int32
	NumModels  int32
	_          int32
	ModelIndex int32
} // 16 bytes

func (d *Decoder) decodeBodyParts(mdl *MDL) error {
	mdl.BodyParts = make([]*BodyPart, mdl.Header.BodyPartCount)
	for i := range mdl.BodyParts {
		header := new(BodyPartHeader)

		if err := d.read(header); err != nil {
			return err
		}
		bp := new(BodyPart)
		bp.Header = header
		bp.Name = d.readName(bp.Header.NameIndex - 16)

		err := d.ppush(
			header.ModelIndex-16,
			func() error { return d.decodeModel(bp) },
		)
		if err != nil {
			return err
		}

		mdl.BodyParts[i] = bp
	}
	return nil
}
