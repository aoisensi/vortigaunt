package mdl

import (
	"fmt"
)

type SeqDesc struct {
	Header       *SeqDescHeader
	Label        string
	ActivityName string
}

type SeqDescHeader struct {
	_                     int32
	LabelIndex            int32
	ActivityNameIndex     int32
	Flags                 int32
	Activity              int32
	ActWeight             int32
	NumEvents             int32
	EventIndex            int32
	BBMin                 [3]float32
	BBMax                 [3]float32
	NumBlends             int32
	AnimIndexIndex        int32
	MovementIndex         int32
	GroupSize             [2]int32
	ParamIndex            [2]int32
	ParamStart            [2]float32
	ParamEnd              [2]float32
	ParamParent           int32
	FadeInTime            float32
	FadeOutTime           float32
	LocalEntryNode        int32
	LocalExitNode         int32
	NodeFlags             int32
	EntryPhase            float32
	ExitPhase             float32
	LastFrame             float32
	NextSeq               int32
	Pose                  int32
	NumIKRules            int32
	NumAutoLayers         int32
	AutoLayerIndex        int32
	WeightListIndex       int32
	PoseKeyIndex          int32
	NumIKLocks            int32
	IKLockIndex           int32
	KeyValueIndex         int32
	KeyValueSize          int32
	CyclePoseIndex        int32
	ActivityModifierIndex int32
	NumActivityModifiers  int32
	_                     [5]int32
} // 212 bytes

func (d *Decoder) decodeSeqDescs(mdl *MDL) error {
	mdl.SecDescss = make([]*SeqDesc, 0, mdl.Header.LocalSeqCount)

	for range mdl.Header.LocalSeqCount {
		header := new(SeqDescHeader)
		if err := d.read(header); err != nil {
			return fmt.Errorf("vortigaunt: decodeSeqDescs: %w", err)
		}
		sd := new(SeqDesc)
		sd.Header = header

		sd.Label = d.readName(header.LabelIndex - 212)
		mdl.SecDescss = append(mdl.SecDescss, sd)
	}
	return nil
}
