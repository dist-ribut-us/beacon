package main

import (
	"github.com/dist-ribut-us/beacon"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/prog"
)

func main() {
	log.Contents = log.Truncate
	log.Panic(log.ToFile(prog.Root() + "beacon.log"))
	log.Go()
	log.SetDebug(true)

	proc, pool, _, err := prog.ReadArgs()
	log.Panic(err)

	beacon.New(proc, pool).Run()
}
