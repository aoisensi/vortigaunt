package vortigaunt

import (
	"flag"
	"log"
)

func init() {
	flag.Parse()

}

func Run() {
	for _, n := range flag.Args() {
		if err := convert(n); err != nil {
			log.Println(err)
		}
	}
}
