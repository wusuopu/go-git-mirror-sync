package services_test

import (
	"app/di"
	"app/initialize"
	"app/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	initialize.InitServices()
	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	assert.Equal(t, 3, di.Container.RepositoryService.Add(1, 2))
}

func TestMinus(t *testing.T) {
	assert.Equal(t, 3, di.Container.RepositoryService.Minus(5, 2))
}

type MyMockObject struct {
	mock.Mock
	services.RepositoryService
}

func (m *MyMockObject) Minus(a, b int) int {
	args := m.Called(a, b)
	return args.Int(0)
}

func TestMinus2(t *testing.T) {
	var oldService = di.Container.RepositoryService

	testObj := new(MyMockObject)
	testObj.On("Minus", 2, 1).Return(10)
	di.Container.RepositoryService = testObj

	assert.Equal(t, 10, di.Container.RepositoryService.Minus(2, 1))
	testObj.AssertCalled(t, "Minus", 2, 1)

	di.Container.RepositoryService = oldService

	assert.Equal(t, 3, di.Container.RepositoryService.Add(1, 2))
}
