package main

import (
	"log"
	"os"

	"github.com/mvo5/uboot-env/uboot"
)

func main() {
	// FIXME: argsparse ftw!
	envFile := os.Args[1]
	env, err := uboot.ReadUbootEnv(envFile)
	if err != nil {
		log.Fatalf("readUbootEnv failed for %s: %s", envFile, err)
	}

	switch os.Args[2] {
	case "print":
		uboot.PrintEnv(env)
	case "set":
		name := os.Args[3]
		value := os.Args[4]
		if err := uboot.SetEnv(env, name, value); err != nil {
			log.Fatalf("setEnv failed: %s", err)
		}
		if err := uboot.WriteUbootEnv(envFile, env); err != nil {
			log.Fatalf("writeUbootEnv failed for %s: %s", envFile, err)
		}
	}

}
