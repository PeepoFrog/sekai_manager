package main

import (
	"log"

	"github.com/PeepoFrog/sekai_manager/src/cfg"
	"github.com/PeepoFrog/sekai_manager/src/cmd"
)

func main() {
	cfg, err := cfg.DefaultCfg()
	if err != nil {
		log.Fatal(err)
	}
	root := cmd.NewRootCmd(cfg)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
