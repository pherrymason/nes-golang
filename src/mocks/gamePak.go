package mocks

import (
	"github.com/raulferras/nes-golang/src/nes/gamePak"
	"github.com/stretchr/testify/mock"
)

type MockableGamePak struct {
	mock.Mock
}

func (m *MockableGamePak) Header() gamePak.Header {
	args := m.Called()

	return args.Get(0).(gamePak.Header)
}

type MockableHeader struct {
	mock.Mock
}

func (m *MockableHeader) ProgramSize() byte {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) CHRSize() byte {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) Mirroring() byte {
	args := m.Called()

	return args.Get(0).(byte)
}

func (m *MockableHeader) HasPersistentMemory() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) HasTrainer() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) IgnoreMirroringControl() bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) MapperNumber() byte {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) PRGRAM() byte {
	//TODO implement me
	panic("implement me")
}

func (m *MockableHeader) TvSystem() byte {
	//TODO implement me
	panic("implement me")
}
