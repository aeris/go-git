package filesystem

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/cozy/go-git/plumbing"
	"github.com/cozy/go-git/utils/merkletrie"
	"github.com/cozy/go-git/utils/merkletrie/noder"
	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v2"
	"gopkg.in/src-d/go-billy.v2/memfs"
)

func Test(t *testing.T) { TestingT(t) }

type NoderSuite struct{}

var _ = Suite(&NoderSuite{})

func (s *NoderSuite) TestDiff(c *C) {
	fsA := memfs.New()
	WriteFile(fsA, "foo", []byte("foo"), 0644)
	WriteFile(fsA, "qux/bar", []byte("foo"), 0644)
	WriteFile(fsA, "qux/qux", []byte("foo"), 0644)

	fsB := memfs.New()
	WriteFile(fsB, "foo", []byte("foo"), 0644)
	WriteFile(fsB, "qux/bar", []byte("foo"), 0644)
	WriteFile(fsB, "qux/qux", []byte("foo"), 0644)

	ch, err := merkletrie.DiffTree(
		NewRootNode(fsA, nil),
		NewRootNode(fsB, nil),
		IsEquals,
	)

	c.Assert(err, IsNil)
	c.Assert(ch, HasLen, 0)
}

func (s *NoderSuite) TestDiffChangeContent(c *C) {
	fsA := memfs.New()
	WriteFile(fsA, "foo", []byte("foo"), 0644)
	WriteFile(fsA, "qux/bar", []byte("foo"), 0644)
	WriteFile(fsA, "qux/qux", []byte("foo"), 0644)

	fsB := memfs.New()
	WriteFile(fsB, "foo", []byte("foo"), 0644)
	WriteFile(fsB, "qux/bar", []byte("bar"), 0644)
	WriteFile(fsB, "qux/qux", []byte("foo"), 0644)

	ch, err := merkletrie.DiffTree(
		NewRootNode(fsA, nil),
		NewRootNode(fsB, nil),
		IsEquals,
	)

	c.Assert(err, IsNil)
	c.Assert(ch, HasLen, 1)
}

func (s *NoderSuite) TestDiffChangeMissing(c *C) {
	fsA := memfs.New()
	WriteFile(fsA, "foo", []byte("foo"), 0644)

	fsB := memfs.New()
	WriteFile(fsB, "bar", []byte("bar"), 0644)

	ch, err := merkletrie.DiffTree(
		NewRootNode(fsA, nil),
		NewRootNode(fsB, nil),
		IsEquals,
	)

	c.Assert(err, IsNil)
	c.Assert(ch, HasLen, 2)
}

func (s *NoderSuite) TestDiffChangeMode(c *C) {
	fsA := memfs.New()
	WriteFile(fsA, "foo", []byte("foo"), 0644)

	fsB := memfs.New()
	WriteFile(fsB, "foo", []byte("foo"), 0755)

	ch, err := merkletrie.DiffTree(
		NewRootNode(fsA, nil),
		NewRootNode(fsB, nil),
		IsEquals,
	)

	c.Assert(err, IsNil)
	c.Assert(ch, HasLen, 1)
}

func (s *NoderSuite) TestDiffChangeModeNotRelevant(c *C) {
	fsA := memfs.New()
	WriteFile(fsA, "foo", []byte("foo"), 0644)

	fsB := memfs.New()
	WriteFile(fsB, "foo", []byte("foo"), 0655)

	ch, err := merkletrie.DiffTree(
		NewRootNode(fsA, nil),
		NewRootNode(fsB, nil),
		IsEquals,
	)

	c.Assert(err, IsNil)
	c.Assert(ch, HasLen, 0)
}

func (s *NoderSuite) TestDiffDirectory(c *C) {
	fsA := memfs.New()
	fsA.MkdirAll("qux/bar", 0644)

	fsB := memfs.New()
	fsB.MkdirAll("qux/bar", 0644)

	ch, err := merkletrie.DiffTree(
		NewRootNode(fsA, map[string]plumbing.Hash{
			"qux/bar": plumbing.NewHash("aa102815663d23f8b75a47e7a01965dcdc96468c"),
		}),
		NewRootNode(fsB, map[string]plumbing.Hash{
			"qux/bar": plumbing.NewHash("19102815663d23f8b75a47e7a01965dcdc96468c"),
		}),
		IsEquals,
	)

	c.Assert(err, IsNil)
	c.Assert(ch, HasLen, 1)

	a, err := ch[0].Action()
	c.Assert(err, IsNil)
	c.Assert(a, Equals, merkletrie.Modify)
}

func WriteFile(fs billy.Filesystem, filename string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

var empty = make([]byte, 24)

func IsEquals(a, b noder.Hasher) bool {
	if bytes.Equal(a.Hash(), empty) || bytes.Equal(b.Hash(), empty) {
		return false
	}

	return bytes.Equal(a.Hash(), b.Hash())
}
