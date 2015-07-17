package uboot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strings"
)

// FIXME: add config option for that so that the user can select if
//        he/she wants env with or without flags
var headerSize = 5

type Env struct {
	fname string
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

func writeUint32(u uint32) []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, &u)
	return buf.Bytes()
}

func CreateEnv(fname string, size int) (*Env, error) {
	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content := make([]byte, size)
	env := &Env{
		fname: fname,
		crc:   crc32.ChecksumIEEE(content[headerSize:]),
		flags: content[4],
		data:  content[headerSize:],
	}

	return env, nil

}

func NewEnv(fname string) (*Env, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	env := &Env{
		fname: fname,
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

func (env *Env) String() string {
	out := ""
	for _, envStr := range bytes.Split(env.data, []byte{0}) {
		if len(envStr) == 0 || envStr[0] == 0 || envStr[0] == 255 {
			continue
		}
		out += string(envStr) + "\n"
	}

	return out
}

func (env *Env) Set(name, value string) error {
	envStrs := []string{}

	for _, envStr := range bytes.Split(env.data, []byte{0}) {
		if len(envStr) == 0 || envStr[0] == 0 || envStr[0] == 255 {
			continue
		}
		// remove str we want to rewrite
		if strings.HasPrefix(string(envStr), fmt.Sprintf("%s=", name)) {
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

func (env *Env) Write() error {
	// FIXME: need to open so that we keep the file intact,
	//        i.e. no truncate, we always override with the full
	//        size anyway
	f, err := os.Create(env.fname)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(writeUint32(env.crc))
	f.Write([]byte{env.flags})
	f.Write(env.data)

	return nil
}
