package vortigaunt

import (
	"flag"
	"log"
)

var flagScale float64
var flagWithLODs bool
var flagForceStatic bool
var flagForceSkeletal bool

func init() {
	flag.Float64Var(&flagScale, "scale", 0.02, "scale factor")
	flag.BoolVar(&flagWithLODs, "with-lods", false, "export all LODs")
	flag.BoolVar(&flagForceStatic, "force-static", false, "force static prop")
	flag.BoolVar(&flagForceSkeletal, "force-skeletal", false, "force skeletal prop")
	flag.Parse()
}

func Run() {
	for _, n := range flag.Args() {
		if err := convert(n); err != nil {
			log.Println(err)
		}
	}
}
