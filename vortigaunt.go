package vortigaunt

import (
	"flag"
	"log"
)

var flagScale float64

func init() {
	flag.Float64Var(&flagScale, "scale", 0.02, "scale factor")
	flag.Parse()
}

func Run() {
	for _, n := range flag.Args() {
		if err := convert(n); err != nil {
			log.Println(err)
		}
	}
}
