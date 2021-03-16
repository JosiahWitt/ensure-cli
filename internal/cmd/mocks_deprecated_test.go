package cmd_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure-cli/internal/cmd"
	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_context"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_exitcleanup"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_mockgen"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestGenerateMocks(t *testing.T) {
	ensure := ensure.New(t)

	type ContextKey struct{}

	type Mocks struct {
		Context          *mock_context.MockContext `ensure:"ignoreunused"`
		EnsureFileLoader *mock_ensurefile.MockLoaderIface
		MockGen          *mock_mockgen.MockMockGenerator
		Cleanup          *mock_exitcleanup.MockExitCleaner
	}

	exampleError := errors.New("something went wrong")
	defaultWd := func() (string, error) {
		return "/test", nil
	}

	table := []struct {
		Name          string
		ExpectedError error
		Flags         []string

		Getwd      func() (string, error)
		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *cmd.App
	}{
		{
			Name:  "with valid execution",
			Getwd: defaultWd,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().
					LoadConfig("/test").
					Return(&ensurefile.Config{
						RootPath: "/some/root/path",
					}, nil)

				ctx := context.WithValue(m.Context, ContextKey{}, "123")
				m.Cleanup.EXPECT().ToContext(gomock.Any()).Return(ctx)

				m.MockGen.EXPECT().
					GenerateMocks(ctx, &ensurefile.Config{
						RootPath: "/some/root/path",
					}).
					Return(nil)
			},
		},

		{
			Name:  "with valid execution: disabled parallel generation",
			Flags: []string{"--disable-parallel"},
			Getwd: defaultWd,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().
					LoadConfig("/test").
					Return(&ensurefile.Config{
						RootPath: "/some/root/path",
					}, nil)

				ctx := context.WithValue(m.Context, ContextKey{}, "123")
				m.Cleanup.EXPECT().ToContext(gomock.Any()).Return(ctx)

				m.MockGen.EXPECT().
					GenerateMocks(ctx, &ensurefile.Config{
						RootPath:                  "/some/root/path",
						DisableParallelGeneration: true,
					}).
					Return(nil)
			},
		},

		{
			Name:          "when error loading working directory",
			Getwd:         func() (string, error) { return "", exampleError },
			ExpectedError: exampleError,
		},

		{
			Name:          "when cannot load config",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(nil, exampleError)
			},
		},

		{
			Name:          "when cannot generate mocks",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().
					LoadConfig("/test").
					Return(&ensurefile.Config{
						RootPath: "/some/root/path",
					}, nil)

				ctx := context.WithValue(m.Context, ContextKey{}, "123")
				m.Cleanup.EXPECT().ToContext(gomock.Any()).Return(ctx)

				m.MockGen.EXPECT().
					GenerateMocks(ctx, &ensurefile.Config{
						RootPath: "/some/root/path",
					}).
					Return(exampleError)
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]
		entry.Subject.Logger = log.New(ioutil.Discard, "", 0)
		entry.Subject.Getwd = entry.Getwd

		err := entry.Subject.Run(append([]string{"ensure", "generate", "mocks"}, entry.Flags...))
		ensure(err).IsError(entry.ExpectedError)
	})
}
