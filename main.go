package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mvo5/uboot-env/uboot"
)

func main() {
	// FIXME: argsparse ftw!
	envFile := os.Args[1]
	cmd := os.Args[2]

	switch cmd {
	case "print":
		env, err := uboot.NewEnv(envFile)
		if err != nil {
			log.Fatalf("readUbootEnv failed for %s: %s", envFile, err)
		}
		fmt.Print(env)
	case "create":
		size, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("Atoi failed for %s: %s", envFile, err)
		}
		env, err := uboot.CreateEnv(envFile, size)
		if err != nil {
			log.Fatalf("env.Create failed for %s: %s", envFile, err)
		}
		if err := env.Write(); err != nil {
			log.Fatalf("env.Write failed: %s", err)
		}

	case "set":
		env, err := uboot.NewEnv(envFile)
		if err != nil {
			log.Fatalf("readUbootEnv failed for %s: %s", envFile, err)
		}
		name := os.Args[3]
		value := os.Args[4]
		if err := env.Set(name, value); err != nil {
			log.Fatalf("env.Set failed: %s", err)
		}
		if err := env.Write(); err != nil {
			log.Fatalf("env.Write failed for %s: %s", envFile, err)
		}
	default:
		log.Fatalf("unknown command %s", cmd)
	}

}
