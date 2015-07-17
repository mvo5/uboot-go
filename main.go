package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mvo5/uboot-env/uboot"
)

func main() {
	// FIXME: argsparse ftw!
	envFile := os.Args[1]
	env, err := uboot.NewEnv(envFile)
	if err != nil {
		log.Fatalf("readUbootEnv failed for %s: %s", envFile, err)
	}

	switch os.Args[2] {
	case "print":
		fmt.Print(env)
	case "set":
		name := os.Args[3]
		value := os.Args[4]
		if err := env.Set(name, value); err != nil {
			log.Fatalf("env.Set failed: %s", err)
		}
		if err := env.Write(); err != nil {
			log.Fatalf("env.Write failed for %s: %s", envFile, err)
		}
	}

}
