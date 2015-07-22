package uenv

import (
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
	env.Save()

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
	// order is input order
	c.Assert(env.String(), Equals, "foo=bar\nbaz=baz\n")
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

func (u *uenvTestSuite) TestSetOrder(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)

	// set in order
	env.Set("foo", "bar")
	env.Set("baz", "baz")
	env.Set("hrm", "hrm")
	// original order of baz is kept
	env.Set("baz", "baz")

	// order is set order
	c.Assert(env.String(), Equals, "foo=bar\nbaz=baz\nhrm=hrm\n")
}

func (u *uenvTestSuite) TestSetOrderPathological(c *C) {
	env, err := Create(u.envFile, 4096)
	c.Assert(err, IsNil)

	// set in order
	env.Set("delme", "delme")
	env.Set("delme2", "delme2")
	env.Set("first", "first")

	// now delete some keys
	env.Set("delme", "")
	env.Set("delme2", "")

	// add a new one
	env.Set("second", "second")

	// ensure order is kept
	c.Assert(env.String(), Equals, "first=first\nsecond=second\n")
}
