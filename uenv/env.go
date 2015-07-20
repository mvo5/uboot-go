package uenv

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

// Env contains the data of the uboot environment
type Env struct {
	fname string
	size  int
	data  map[string]string
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

// Create a new empty uboot env file with the given size
func Create(fname string, size int) (*Env, error) {
	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	env := &Env{
		fname: fname,
		size:  size,
		data:  make(map[string]string),
	}

	return env, nil
}

// Open opens a existing uboot env file
func Open(fname string) (*Env, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	crc := readUint32(content)
	actualCRC := crc32.ChecksumIEEE(content[headerSize:])
	if crc != actualCRC {
		return nil, fmt.Errorf("bad CRC: %v != %v", crc, actualCRC)
	}

	env := &Env{
		fname: fname,
		size:  len(content),
		data:  parseData(content[headerSize:]),
	}

	return env, nil
}

func parseData(data []byte) map[string]string {
	out := make(map[string]string)

	for _, envStr := range bytes.Split(data, []byte{0}) {
		if len(envStr) == 0 || envStr[0] == 0 || envStr[0] == 255 {
			continue
		}
		l := strings.SplitN(string(envStr), "=", 2)
		out[l[0]] = l[1]
	}

	return out
}

func (env *Env) String() string {
	out := ""

	for k, v := range env.data {
		out += fmt.Sprintf("%s=%s\n", k, v)
	}

	return out
}

// Get returns the value of the environment variable of the given name
func (env *Env) Get(name string) string {
	return env.data[name]
}

// Set sets an environment name to the given value, if the value is empty
// the variable will be removed from the environment
func (env *Env) Set(name, value string) {
	env.data[name] = value
}

// Save will write out the environment data
func (env *Env) Save() error {
	w := bytes.NewBuffer(nil)
	w.Grow(env.size)
	for k, v := range env.data {
		w.Write([]byte(fmt.Sprintf("%s=%s", k, v)))
		w.Write([]byte{0})
	}
	crc := crc32.ChecksumIEEE(w.Bytes())

	// the size of the env file never changes so we not truncate
	// we also do not O_TRUNC to avoid reallocations on the FS
	// to minimize risk of fs corruption
	f, err := os.OpenFile(env.fname, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(writeUint32(crc))
	// padding bytes (e.g. for redundant header)
	pad := make([]byte, headerSize-binary.Size(crc))
	f.Write(pad)
	f.Write(w.Bytes())

	return nil
}
