package categories_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type CategoriesServiceTestSuite struct {
	suite.Suite
}

func (ts *CategoriesServiceTestSuite) SetupTest() {

}

func TestCategoriesService(t *testing.T) {
	suite.Run(t, new(CategoriesServiceTestSuite))
}
