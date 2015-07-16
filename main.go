package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// FIXME: add config option for that so that the user can select if
//        he/she wants env with or without flags
var headerSize = 5

type EnvHeader struct {
	crc   uint32
	flags byte
	data  []byte
}

func readUint32(data []byte) uint32 {
	var ret uint32
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}

func readUbootEnv(fname string) (*EnvHeader, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	env := &EnvHeader{
		crc:   readUint32(content),
		flags: content[4],
		data:  content[headerSize:],
	}

	actualCRC := crc32.ChecksumIEEE(content[headerSize:])
	if env.crc != actualCRC {
		return nil, fmt.Errorf("bad CRC: %v != %v", env.crc, actualCRC)
	}

	return env, nil
}

func printEnv(env *EnvHeader) error {
	for _, envStr := range bytes.Split(env.data, []byte{0}) {
		if len(envStr) == 0 || envStr[0] == 0 || envStr[0] == 255 {
			continue
		}
		fmt.Println(string(envStr))
	}

	return nil
}

func setEnv(env *EnvHeader, name, value string) error {
	envStrs := []string{}

	for _, envStr := range bytes.Split(env.data, []byte{0}) {
		if len(envStr) == 0 || envStr[0] == 0 || envStr[0] == 255 {
			continue
		}
		// remove str we want to rewrite
		if strings.HasPrefix(string(envStr), fmt.Sprintf("%s=", value)) {
			continue
		}
		envStrs = append(envStrs, string(envStr))
	}

	// append new str
	if value != "" {
		envStrs = append(envStrs, fmt.Sprintf("%s=%s", name, value))
	}

	w := bytes.NewBuffer(nil)
	for _, envStr := range envStrs {
		//println(envStr)
		fmt.Fprintf(w, "%s", envStr)
		w.Write([]byte{0})
	}
	pad := make([]byte, len(env.data)-w.Len())
	w.Write(pad)
	env.data = w.Bytes()
	env.crc = crc32.ChecksumIEEE(env.data)

	return nil
}

func writeUint32(u uint32) []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, &u)
	return buf.Bytes()
}

func writeUbootEnv(envFile string, env *EnvHeader) error {
	f, err := os.Create(envFile)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(writeUint32(env.crc))
	f.Write([]byte{env.flags})
	f.Write(env.data)

	return nil
}

func main() {
	// FIXME: argsparse ftw!
	envFile := os.Args[1]
	env, err := readUbootEnv(envFile)
	if err != nil {
		log.Fatalf("readUbootEnv failed for %s: %s", envFile, err)
	}

	switch os.Args[2] {
	case "print":
		printEnv(env)
	case "set":
		name := os.Args[3]
		value := os.Args[4]
		if err := setEnv(env, name, value); err != nil {
			log.Fatalf("setEnv failed: %s", err)
		}
		if err := writeUbootEnv(envFile, env); err != nil {
			log.Fatalf("writeUbootEnv failed for %s: %s", envFile, err)
		}
	}

}
