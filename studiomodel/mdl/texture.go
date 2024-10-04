package mdl

type Texture struct {
	Header *TextureHeader
	Name   string
}

type TextureHeader struct {
	NameIndex int32
	Flags     int32
	Used      int32
	_         [3]int32
	_         [10]int32
}

func (d *Decoder) decodeTexture(mdl *MDL) error {
	mdl.Textures = make([]*Texture, mdl.Header.TextureCount)
	for i := range mdl.Textures {
		header := new(TextureHeader)
		if err := d.read(header); err != nil {
			return err
		}
		texture := new(Texture)
		texture.Header = header
		texture.Name = d.readName(header.NameIndex - 64)
		mdl.Textures[i] = texture
	}
	return nil
}
