package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/e1m0re/grdn/internal/storage/store"
	"github.com/e1m0re/grdn/internal/storage/store/mocks"
)

func TestConnection(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		mockStore func() store.Store
		args      args
		want      want
		name      string
	}{
		{
			name: "something wrong",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: errors.New("something wrong"),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Ping", mock.Anything).
					Return(errors.New("something wrong"))

				return mockStore
			},
		},
		{
			name: "successfully case",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Ping", mock.Anything).
					Return(nil)

				return mockStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &service{
				Store: test.mockStore(),
			}
			err := s.TestConnection(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestNewService(t *testing.T) {
	mockStore := mocks.NewStore(t)
	got := NewService(mockStore)
	assert.Implements(t, (*Service)(nil), got)
}

func Test_service_Clear(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		mockStore func() store.Store
		args      args
		want      want
		name      string
	}{
		{
			name: "something wrong",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: errors.New("something wrong"),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Clear", mock.Anything).
					Return(errors.New("something wrong"))

				return mockStore
			},
		},
		{
			name: "successfully case",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Clear", mock.Anything).
					Return(nil)

				return mockStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &service{
				Store: test.mockStore(),
			}
			err := s.Clear(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_service_Close(t *testing.T) {
	type want struct {
		err error
	}
	tests := []struct {
		mockStore func() store.Store
		want      want
		name      string
	}{
		{
			name: "something wrong",
			want: want{
				err: errors.New("something wrong"),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Close", mock.Anything).
					Return(errors.New("something wrong"))

				return mockStore
			},
		},
		{
			name: "successfully case",
			want: want{
				err: nil,
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Close", mock.Anything).
					Return(nil)

				return mockStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &service{
				Store: test.mockStore(),
			}
			err := s.Close()
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_service_Restore(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		mockStore func() store.Store
		args      args
		want      want
		name      string
	}{
		{
			name: "something wrong",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: errors.New("something wrong"),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Restore", mock.Anything).
					Return(errors.New("something wrong"))

				return mockStore
			},
		},
		{
			name: "successfully case",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Restore", mock.Anything).
					Return(nil)

				return mockStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &service{
				Store: test.mockStore(),
			}
			err := s.Restore(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func Test_service_Save(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		mockStore func() store.Store
		args      args
		want      want
		name      string
	}{
		{
			name: "something wrong",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: errors.New("something wrong"),
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Save", mock.Anything).
					Return(errors.New("something wrong"))

				return mockStore
			},
		},
		{
			name: "successfully case",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				err: nil,
			},
			mockStore: func() store.Store {
				mockStore := mocks.NewStore(t)
				mockStore.
					On("Save", mock.Anything).
					Return(nil)

				return mockStore
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &service{
				Store: test.mockStore(),
			}
			err := s.Save(test.args.ctx)
			assert.Equal(t, test.want.err, err)
		})
	}
}
