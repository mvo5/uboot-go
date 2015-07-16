package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
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

func main() {
	envFile := os.Args[1]
	env, err := readUbootEnv(envFile)
	if err != nil {
		log.Fatalf("readUboot failed for %s: %s", envFile, err)
	}
	printEnv(env)
}
