package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// parse command line args
	port := flag.Int("p", 8080, "-p <port>")
	flag.Parse()

	g := gin.Default()
	initRoute(g)
	err := g.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
