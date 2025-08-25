package main

import (
	"fmt"
	"github.com/simonkosina/bootdev-blog-aggregator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: &v", err)
	}

	err = cfg.SetUser("simon")
	if err != nil {
		log.Fatalf("Error set current user: &v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: &v", err)
	}

	fmt.Println(cfg)
}
