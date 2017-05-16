package uenv

import (
	"bytes"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up check.v1 into the "go test" runner
func Test(t *testing.T) { TestingT(t) }

type uenvTestSuite struct {
	envFile string
}

var _ = Suite(&uenvTestSuite{})

func (u *uenvTestSuite) SetUpTest(c *C) {
	u.envFile = filepath.Join(c.MkDir(), "uboot.env")
}

func (u *uenvTestSuite) TestSetNoDuplicate(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)
	env.Set("foo", "bar")
	env.Set("foo", "bar")
	c.Assert(env.String(), Equals, "foo=bar\n")
}

func (u *uenvTestSuite) TestOpenEnv(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)
	env.Set("foo", "bar")
	c.Assert(env.String(), Equals, "foo=bar\n")
	err = env.Save()
	c.Assert(err, IsNil)

	env2, err := Open(u.envFile)
	c.Assert(err, IsNil)
	c.Assert(env2.String(), Equals, "foo=bar\n")
}

func (u *uenvTestSuite) TestGetSimple(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)
	env.Set("foo", "bar")
	c.Assert(env.Get("foo"), Equals, "bar")
}

func (u *uenvTestSuite) TestGetNoSuchEntry(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)
	c.Assert(env.Get("no-such-entry"), Equals, "")
}

func (u *uenvTestSuite) TestImport(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)

	r := strings.NewReader("foo=bar\n#comment\n\nbaz=baz")
	err = env.Import(r)
	c.Assert(err, IsNil)
	// order is alphabetic
	c.Assert(env.String(), Equals, "baz=baz\nfoo=bar\n")
}

func (u *uenvTestSuite) TestImportHasError(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)

	r := strings.NewReader("foxy")
	err = env.Import(r)
	c.Assert(err, ErrorMatches, "Invalid line: \"foxy\"")
}

func (u *uenvTestSuite) TestSetEmptyUnsets(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)

	env.Set("foo", "bar")
	c.Assert(env.String(), Equals, "foo=bar\n")
	env.Set("foo", "")
	c.Assert(env.String(), Equals, "")
}

func (u *uenvTestSuite) makeUbootEnvFromData(c *C, mockData []byte) {
	w := bytes.NewBuffer(nil)
	crc := crc32.ChecksumIEEE(mockData)
	w.Write(writeUint32(crc))
	w.Write([]byte{0})
	w.Write(mockData)

	f, err := os.Create(u.envFile)
	c.Assert(err, IsNil)
	defer f.Close()
	_, err = f.Write(w.Bytes())
	c.Assert(err, IsNil)
}

// ensure that the data after \0\0 is discarded (except for crc)
func (u *uenvTestSuite) TestReadStopsAfterDoubleNull(c *C) {
	mockData := []byte{
		// foo=bar
		0x66, 0x6f, 0x6f, 0x3d, 0x62, 0x61, 0x72,
		// eof
		0x00, 0x00,
		// junk after eof as written by fw_setenv sometimes
		// =b
		0x3d, 62,
		// empty
		0xff, 0xff,
	}
	u.makeUbootEnvFromData(c, mockData)

	env, err := Open(u.envFile)
	c.Assert(err, IsNil)
	c.Assert(env.String(), Equals, "foo=bar\n")
}

// ensure that the malformed data is not causing us to panic.
func (u *uenvTestSuite) TestErrorOnMalformedData(c *C) {
	mockData := []byte{
		// foo
		0x66, 0x6f, 0x6f,
		// eof
		0x00, 0x00,
	}
	u.makeUbootEnvFromData(c, mockData)

	env, err := Open(u.envFile)
	c.Assert(err, ErrorMatches, `cannot parse line "foo" as key=value pair`)
	c.Assert(env, IsNil)
}

// ensure that the malformed data is not causing us to panic.
func (u *uenvTestSuite) TestOpenBestEffort(c *C) {
	mockData := []byte{
		// key1=value1
		0x6b, 0x65, 0x79, 0x31, 0x3d, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, 0x00,
		// foo
		0x66, 0x6f, 0x6f, 0x00,
		// key2=value2
		0x6b, 0x65, 0x79, 0x32, 0x3d, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x00,
		// eof
		0x00, 0x00,
	}
	u.makeUbootEnvFromData(c, mockData)

	env, err := OpenWithFlags(u.envFile, OpenBestEffort)
	c.Assert(err, IsNil)
	c.Assert(env.String(), Equals, "key1=value1\nkey2=value2\n")
}

func (u *uenvTestSuite) TestReadEmptyFile(c *C) {
	mockData := []byte{
		// eof
		0x00, 0x00,
		// empty
		0xff, 0xff,
	}
	u.makeUbootEnvFromData(c, mockData)

	env, err := Open(u.envFile)
	c.Assert(err, IsNil)
	c.Assert(env.String(), Equals, "")
}

func (u *uenvTestSuite) TestWritesEmptyFileWithDoubleNewline(c *C) {
	env, err := Create(u.envFile, 12)
	c.Assert(err, IsNil)
	err = env.Save()
	c.Assert(err, IsNil)

	r, err := os.Open(u.envFile)
	c.Assert(err, IsNil)
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	c.Assert(err, IsNil)
	c.Assert(content, DeepEquals, []byte{
		// crc
		0x11, 0x38, 0xb3, 0x89,
		// redundant
		0x0,
		// eof
		0x0, 0x0,
		// footer
		0xff, 0xff, 0xff, 0xff, 0xff,
	})

	env, err = Open(u.envFile)
	c.Assert(err, IsNil)
	c.Assert(env.String(), Equals, "")
}

func (u *uenvTestSuite) TestWritesContentCorrectly(c *C) {
	totalSize := 16

	env, err := Create(u.envFile, totalSize)
	c.Assert(err, IsNil)
	env.Set("a", "b")
	env.Set("c", "d")
	err = env.Save()
	c.Assert(err, IsNil)

	r, err := os.Open(u.envFile)
	c.Assert(err, IsNil)
	defer r.Close()
	content, err := ioutil.ReadAll(r)
	c.Assert(err, IsNil)
	c.Assert(content, DeepEquals, []byte{
		// crc
		0xc7, 0xd9, 0x6b, 0xc5,
		// redundant
		0x0,
		// a=b
		0x61, 0x3d, 0x62,
		// eol
		0x0,
		// c=d
		0x63, 0x3d, 0x64,
		// eof
		0x0, 0x0,
		// footer
		0xff, 0xff,
	})

	env, err = Open(u.envFile)
	c.Assert(err, IsNil)
	c.Assert(env.String(), Equals, "a=b\nc=d\n")
	c.Assert(env.size, Equals, totalSize)
}
