package main

import (
	"github.com/korchasa/kogdaeda/web"
	"log"
)

func main() {
	//sch, err := zk.Schedule(50.45767445379344,30.502764582939378)
	//if err != nil {
	//	panic(err)
	//}
	//spew.Dump(sch)
	addr := ":8090"
	s := web.New()
	log.Printf("Listen on %s\n", addr)
	if err := s.Start(addr); err != nil {
		panic(err)
	}
}
