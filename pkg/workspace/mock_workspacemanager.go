// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package workspace

import (
	"sync"
)

var (
	lockManagerMockCheckWorkingDir sync.RWMutex
)

// Ensure, that ManagerMock does implement Manager.
// If this is not the case, regenerate this file with moq.
var _ Manager = &ManagerMock{}

// ManagerMock is a mock implementation of Manager.
//
//     func TestSomethingThatUsesManager(t *testing.T) {
//
//         // make and configure a mocked Manager
//         mockedManager := &ManagerMock{
//             CheckWorkingDirFunc: func() error {
// 	               panic("mock out the CheckWorkingDir method")
//             },
//         }
//
//         // use mockedManager in code that requires Manager
//         // and then make assertions.
//
//     }
type ManagerMock struct {
	// CheckWorkingDirFunc mocks the CheckWorkingDir method.
	CheckWorkingDirFunc func() error

	// calls tracks calls to the methods.
	calls struct {
		// CheckWorkingDir holds details about calls to the CheckWorkingDir method.
		CheckWorkingDir []struct {
		}
	}
}

// CheckWorkingDir calls CheckWorkingDirFunc.
func (mock *ManagerMock) CheckWorkingDir() error {
	if mock.CheckWorkingDirFunc == nil {
		panic("ManagerMock.CheckWorkingDirFunc: method is nil but Manager.CheckWorkingDir was just called")
	}
	callInfo := struct {
	}{}
	lockManagerMockCheckWorkingDir.Lock()
	mock.calls.CheckWorkingDir = append(mock.calls.CheckWorkingDir, callInfo)
	lockManagerMockCheckWorkingDir.Unlock()
	return mock.CheckWorkingDirFunc()
}

// CheckWorkingDirCalls gets all the calls that were made to CheckWorkingDir.
// Check the length with:
//     len(mockedManager.CheckWorkingDirCalls())
func (mock *ManagerMock) CheckWorkingDirCalls() []struct {
} {
	var calls []struct {
	}
	lockManagerMockCheckWorkingDir.RLock()
	calls = mock.calls.CheckWorkingDir
	lockManagerMockCheckWorkingDir.RUnlock()
	return calls
}
