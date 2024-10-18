package vortigaunt

import (
	"flag"
	"log"
)

var flagScale float64
var flagWithLODs bool

func init() {
	flag.Float64Var(&flagScale, "scale", 0.02, "scale factor")
	flag.BoolVar(&flagWithLODs, "with-lods", false, "export all LODs")
	flag.Parse()
}

func Run() {
	for _, n := range flag.Args() {
		if err := convert(n); err != nil {
			log.Println(err)
		}
	}
}
