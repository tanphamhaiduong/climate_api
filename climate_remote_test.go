package climate

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type ClimateRemoteTestSuite struct {
	suite.Suite
	client ClientImpl
}

func TestClimateRemoteTestSuite(t *testing.T) {
	suite.Run(t, new(ClimateRemoteTestSuite))
}

func (s *ClimateRemoteTestSuite) SetupTest() {
	validate := validator.New()
	client := NewClient(http.DefaultClient, validate, "http://climatedataapi.worldbank.org/climateweb/rest/v1")
	s.client = *client
}

func (s *ClimateRemoteTestSuite) TestNewClient_Success() {
	s.NotNil(s.client)
}

func (s *ClimateRemoteTestSuite) TestNewGetRequestWithRelativeURL_Success() {
	var (
		input    = "/country/annualavg/pr/1980/1999/GBR.xml"
		expected = fmt.Sprintf("%s/country/annualavg/pr/1980/1999/GBR.xml", "http://climatedataapi.worldbank.org/climateweb/rest/v1")
		ctx      = context.Background()
	)
	r, err := s.client.NewGetRequest(ctx, input)
	s.Equal(expected, r.URL.String())
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestNewGetRequestWithAbsoluteURL_Success() {
	var (
		input    = fmt.Sprintf("%s/country/annualavg/pr/1980/1999/GBR.xml", "http://climatedataapi.worldbank.org/climateweb/rest/v1")
		expected = fmt.Sprintf("%s/country/annualavg/pr/1980/1999/GBR.xml", "http://climatedataapi.worldbank.org/climateweb/rest/v1")
		ctx      = context.Background()
	)
	r, err := s.client.NewGetRequest(ctx, input)
	s.Equal(expected, r.URL.String())
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestGetAnnualRainfall_Success() {
	var (
		input = GetAnnualRainfallArgs{
			FromCCYY:   "1980",
			ToCCYY:     "1999",
			CountryISO: "GBR",
		}
		ctx = context.Background()
	)
	result, err := s.client.GetAnnualRainfall(ctx, input)
	s.NotNil(result)
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestGetAnnualRainfall_Failed() {
	var (
		input = GetAnnualRainfallArgs{
			FromCCYY:   "1980",
			ToCCYY:     "1999",
			CountryISO: "GB",
		}
		ctx      = context.Background()
		expected = List{}
	)
	result, err := s.client.GetAnnualRainfall(ctx, input)
	s.Equal(expected, result)
	s.NotNil(err)
}

func (s *ClimateRemoteTestSuite) TestCalculateAveAnual_Success() {
	var (
		list = List{
			DomainWebAnnualGcmDatum: []DomainWebAnnualGcmDatum{
				{
					AnnualData: AnnualData{
						Double: "10",
					},
				},
				{
					AnnualData: AnnualData{
						Double: "11",
					},
				},
			},
		}
		fromCCYY = int64(1980)
		toCCYY   = int64(1990)
		expected = decimal.NewFromFloat32(10.5)
	)
	result, err := s.client.calculateAveAnual(list, fromCCYY, toCCYY)
	s.Equal(expected.String(), result.String())
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestCalculateAveAnual_Failed() {
	var (
		list = List{
			DomainWebAnnualGcmDatum: []DomainWebAnnualGcmDatum{},
		}
		fromCCYY = int64(1980)
		toCCYY   = int64(1990)
		expected = decimal.NewFromInt(0)
	)
	result, err := s.client.calculateAveAnual(list, fromCCYY, toCCYY)
	s.Equal(result, expected)
	s.NotNil(err)
}

func (s *ClimateRemoteTestSuite) TestAverageRainfallForGreatBritainFrom1980to1999Exists() {
	var (
		ctx      = context.Background()
		expected = float64(988.8454972331014)
	)
	result, err := s.client.GetAveAnnualRainfall(ctx, 1980, 1999, "gbr")
	s.Equal(expected, result)
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestAverageRainfallForFranceFrom1980to1999Exists() {
	var (
		ctx      = context.Background()
		expected = 913.7986955122727
	)
	result, err := s.client.GetAveAnnualRainfall(ctx, 1980, 1999, "fra")
	s.Equal(expected, result)
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestAverageRainfallForEgyptFrom1980to1999Exists() {
	var (
		ctx      = context.Background()
		expected = float64(54.58587712129825)
	)
	result, err := s.client.GetAveAnnualRainfall(ctx, 1980, 1999, "egy")
	s.Equal(expected, result)
	s.Nil(err)
}

func (s *ClimateRemoteTestSuite) TestAverageRainfallForGreatBritainFrom1985to1995DoesNotExist() {
	var (
		ctx      = context.Background()
		expected = float64(0)
	)
	result, err := s.client.GetAveAnnualRainfall(ctx, 1985, 1995, "gbr")
	s.Equal(expected, result)
	s.Error(err)
}

func (s *ClimateRemoteTestSuite) TestAverageRainfallForMiddleEarthFrom1980to1999DoesNotExist() {
	var (
		ctx      = context.Background()
		expected = float64(0)
	)
	result, err := s.client.GetAveAnnualRainfall(ctx, 1980, 1999, "mde")
	s.Equal(expected, result)
	s.Error(err)
}
