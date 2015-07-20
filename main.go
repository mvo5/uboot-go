package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mvo5/uboot-go/uenv"
)

func main() {
	// FIXME: argsparse ftw!
	envFile := os.Args[1]
	cmd := os.Args[2]

	switch cmd {
	case "print":
		env, err := uenv.Open(envFile)
		if err != nil {
			log.Fatalf("uenv.Open failed for %s: %s", envFile, err)
		}
		fmt.Print(env)
	case "create":
		size, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("Atoi failed for %s: %s", envFile, err)
		}
		env, err := uenv.Create(envFile, size)
		if err != nil {
			log.Fatalf("uenv.Create failed for %s: %s", envFile, err)
		}
		if err := env.Save(); err != nil {
			log.Fatalf("env.Save failed: %s", err)
		}

	case "set":
		env, err := uenv.Open(envFile)
		if err != nil {
			log.Fatalf("uenv.Open failed for %s: %s", envFile, err)
		}
		name := os.Args[3]
		value := os.Args[4]
		env.Set(name, value)
		if err := env.Save(); err != nil {
			log.Fatalf("env.Save failed for %s: %s", envFile, err)
		}
	case "import":
		env, err := uenv.Open(envFile)
		if err != nil {
			log.Fatalf("uenv.Open failed for %s: %s", envFile, err)
		}
		fname := os.Args[3]
		r, err := os.Open(fname)
		if err != nil {
			log.Fatalf("Open failed for %s: %s", fname, err)
		}
		if err := env.Import(r); err != nil {
			log.Fatalf("env.Import failed for %s: %s", envFile, err)
		}
		if err := env.Save(); err != nil {
			log.Fatalf("env.Save failed for %s: %s", envFile, err)
		}
	default:
		log.Fatalf("unknown command %s", cmd)
	}

}
