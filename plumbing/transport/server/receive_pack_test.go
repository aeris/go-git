package server_test

import (
	"github.com/cozy/go-git/plumbing/transport"
	"github.com/cozy/go-git/plumbing/transport/test"

	. "gopkg.in/check.v1"
)

type ReceivePackSuite struct {
	BaseSuite
	test.ReceivePackSuite
}

var _ = Suite(&ReceivePackSuite{})

func (s *ReceivePackSuite) SetUpSuite(c *C) {
	s.BaseSuite.SetUpSuite(c)
	s.ReceivePackSuite.Client = s.client
}

func (s *ReceivePackSuite) SetUpTest(c *C) {
	s.prepareRepositories(c, &s.Endpoint, &s.EmptyEndpoint, &s.NonExistentEndpoint)
}

func (s *ReceivePackSuite) TearDownTest(c *C) {
	s.Suite.TearDownSuite(c)
}

// TODO
func (s *ReceivePackSuite) TestSendPackAddDeleteReference(c *C) {
	c.Skip("delete reference not supported yet")
}

// Overwritten, server returns error earlier.
func (s *ReceivePackSuite) TestAdvertisedReferencesNotExists(c *C) {
	r, err := s.Client.NewReceivePackSession(s.NonExistentEndpoint, s.EmptyAuth)
	c.Assert(err, Equals, transport.ErrRepositoryNotFound)
	c.Assert(r, IsNil)
}
