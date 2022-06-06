package dps_test

import (
	"fmt"
	dps "github.com/deliveryhero/pd-groceries-vendor"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ConsumerSuite struct {
	suite.Suite
}

var (
	pact *dsl.Pact
)

func TestBookClientSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}

func (s *ConsumerSuite) SetupSuite() {
	pact = &dsl.Pact{
		Consumer: "GVS",
		Provider: "DPS",
		PactDir:  "./pacts",
	}
}

func (s *ConsumerSuite) TearDownSuite() {
	pact.Teardown()
}

func (s *ConsumerSuite) TestGetCustomerThatDoesExist() {
	pact.AddInteraction().
		Given("DPS Response for Singapore country").
		UponReceiving("A post request for DPS").
		WithRequest(
			dsl.Request{
				Method: http.MethodPost,
				Path:   dsl.String("/api/v2/fees/FP_SG"),
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
			},
		).
		WillRespondWith(
			dsl.Response{
				Status: http.StatusOK,
				Headers: dsl.MapMatcher{
					"Content-Type": dsl.String("application/json"),
				},
				Body: dsl.Match(dps.Response{}),
			},
		)

	test := func() error {
		c := dps.NewClient(fmt.Sprintf("http://localhost:%d", pact.Server.Port))
		resp, err := c.PostCustomer("FP_SG")

		s.NoError(err)
		s.Equal("Original", resp.Customer.Variant)

		return nil
	}

	s.NoError(pact.Verify(test))
}
