package server_test

import (
	"fmt"
	"testing"

	"github.com/cozy/go-git/plumbing/transport"
	"github.com/cozy/go-git/plumbing/transport/client"
	"github.com/cozy/go-git/plumbing/transport/server"
	"github.com/cozy/go-git/storage/filesystem"
	"github.com/cozy/go-git/storage/memory"
	"github.com/src-d/go-git-fixtures"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

const inprocScheme = "inproc"

type BaseSuite struct {
	fixtures.Suite
	loader       server.MapLoader
	client       transport.Transport
	clientBackup transport.Transport
}

func (s *BaseSuite) SetUpSuite(c *C) {
	s.Suite.SetUpSuite(c)
	s.loader = server.MapLoader{}
	s.client = server.NewServer(s.loader)
	s.clientBackup = client.Protocols[inprocScheme]
	client.Protocols[inprocScheme] = s.client
}

func (s *BaseSuite) TearDownSuite(c *C) {
	if s.clientBackup == nil {
		delete(client.Protocols, inprocScheme)
	} else {
		client.Protocols[inprocScheme] = s.clientBackup
	}
}

func (s *BaseSuite) prepareRepositories(c *C, basic *transport.Endpoint,
	empty *transport.Endpoint, nonExistent *transport.Endpoint) {

	f := fixtures.Basic().One()
	fs := f.DotGit()
	path := fs.Base()
	url := fmt.Sprintf("%s://%s", inprocScheme, path)
	ep, err := transport.NewEndpoint(url)
	c.Assert(err, IsNil)
	*basic = ep
	sto, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)
	s.loader[ep.String()] = sto

	path = "/empty.git"
	url = fmt.Sprintf("%s://%s", inprocScheme, path)
	ep, err = transport.NewEndpoint(url)
	c.Assert(err, IsNil)
	*empty = ep
	s.loader[ep.String()] = memory.NewStorage()

	path = "/non-existent.git"
	url = fmt.Sprintf("%s://%s", inprocScheme, path)
	ep, err = transport.NewEndpoint(url)
	c.Assert(err, IsNil)
	*nonExistent = ep
}
