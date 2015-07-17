package uboot

import (
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up check.v1 into the "go test" runner
func Test(t *testing.T) { TestingT(t) }

type ubootTestSuite struct {
	envFile string
}

var _ = Suite(&ubootTestSuite{})

func (u *ubootTestSuite) SetUpTest(c *C) {
	u.envFile = filepath.Join(c.MkDir(), "uboot.env")
}

func (u *ubootTestSuite) TestSetNoDuplicate(c *C) {
	env, err := CreateEnv(u.envFile, 4096)
	c.Assert(err, IsNil)
	env.Set("foo", "bar")
	env.Set("foo", "bar")
	c.Assert(env.String(), Equals, "foo=bar\n")
}

func (u *ubootTestSuite) TestOpenEnv(c *C) {
	env, err := CreateEnv(u.envFile, 4096)
	c.Assert(err, IsNil)
	env.Set("foo", "bar")
	c.Assert(env.String(), Equals, "foo=bar\n")
	env.Write()

	env2, err := OpenEnv(u.envFile)
	c.Assert(err, IsNil)
	c.Assert(env2.String(), Equals, "foo=bar\n")
}
