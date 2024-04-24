package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type IDoSomething interface {
	DoSomething(int) (bool, error)
}

func targetFuncThatDoesSomethingWithObj(d IDoSomething) int {
	d.DoSomething(123)
	return 123
}

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) DoSomething(number int) (bool, error) {
	args := m.Called(number)
	return args.Bool(0), args.Error(1)
}

func TestMockedObject(t *testing.T) {
	testObject := new(MyMockedObject)

	testObject.On("DoSomething", 123).Return(true, nil)

	assert.Equal(t, 123, targetFuncThatDoesSomethingWithObj(testObject))

	testObject.AssertExpectations(t)
}

func TestSomethingWithPlaceholder(t *testing.T) {

	// create an instance of our test object
	testObj := new(MyMockedObject)

	// set up expectations with a placeholder in the argument list
	testObj.On("DoSomething", mock.Anything).Return(true, nil)

	// call the code we are testing
	targetFuncThatDoesSomethingWithObj(testObj)

	// assert that the expectations were met
	testObj.AssertExpectations(t)

}

func TestSomethingElse2(t *testing.T) {

	// create an instance of our test object
	testObj := new(MyMockedObject)

	// set up expectations with a placeholder in the argument list
	mockCall := testObj.On("DoSomething", mock.Anything).Return(true, nil)

	// call the code we are testing
	targetFuncThatDoesSomethingWithObj(testObj)

	// assert that the expectations were met
	testObj.AssertExpectations(t)

	// remove the handler now so we can add another one that takes precedence
	mockCall.Unset()

	// return false now instead of true
	testObj.On("DoSomething", mock.Anything).Return(false, nil)

	testObj.AssertExpectations(t)
}
