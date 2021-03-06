package mockgen_test

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/mockgen"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_exitcleanup"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_fswrite"
	"github.com/JosiahWitt/ensure-cli/internal/mocks/mock_runcmd"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func TestTidyMocks(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		CmdRun  *mock_runcmd.MockRunnerIface
		FSWrite *mock_fswrite.MockFSWriteIface
		Cleanup *mock_exitcleanup.MockExitCleaner
	}

	table := []struct {
		Name          string
		Config        *ensurefile.Config
		ExpectedError error

		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *mockgen.MockGen
	}{
		{
			Name: "with files to delete",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/some/pkg/qwerty",
							Interfaces: []string{"Iface2"},
						},
						{
							Path:       "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
							Interfaces: []string{"Iface3"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const primaryMocksDir = "/root/path/primary_mocks"
				m.FSWrite.EXPECT().ListRecursive(primaryMocksDir).
					Return([]string{
						primaryMocksDir + "/github.com",
						primaryMocksDir + "/github.com/some",
						primaryMocksDir + "/github.com/some/pkg",
						primaryMocksDir + "/github.com/some/pkg/mock_abc",
						primaryMocksDir + "/github.com/some/pkg/mock_abc/mock_abc.go",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/mock_qwerty.go",

						// Extra files

						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_file.go",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir/with_file.go",

						primaryMocksDir + "/github.com/d3l3t3.m3",
						primaryMocksDir + "/somefile.txt",
						primaryMocksDir + "/some",
						primaryMocksDir + "/some/nesting",
						primaryMocksDir + "/some/nesting/file1.txt",
						primaryMocksDir + "/some/nesting/file2.txt",
						primaryMocksDir + "/some/hello.txt",
					}, nil)

				const internalMocksDir = "/root/path/layer1/layer2/internal/layer3/layer4/internal/internal_mocks"
				m.FSWrite.EXPECT().ListRecursive(internalMocksDir).
					Return([]string{
						internalMocksDir + "/layer5",
						internalMocksDir + "/layer5/layer6",
						internalMocksDir + "/layer5/layer6/mock_xyz",
						internalMocksDir + "/layer5/layer6/mock_xyz/mock_xyz.go",

						// Extra files
						internalMocksDir + "/layer5/layer6/mock_xyz/extra123.go",
						internalMocksDir + "/layer5/layer6/mock_xyz/nested",
						internalMocksDir + "/layer5/layer6/mock_xyz/nested/more.go",
						internalMocksDir + "/layer5/hello",
						internalMocksDir + "/layer5/hello/there.hi",
						internalMocksDir + "/garbage.txt",
					}, nil)

				expectedPathsToDelete := []string{
					// Primary mocks
					primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_file.go",
					primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir",
					primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir/with_file.go",

					primaryMocksDir + "/github.com/d3l3t3.m3",
					primaryMocksDir + "/somefile.txt",
					primaryMocksDir + "/some",
					primaryMocksDir + "/some/nesting",
					primaryMocksDir + "/some/nesting/file1.txt",
					primaryMocksDir + "/some/nesting/file2.txt",
					primaryMocksDir + "/some/hello.txt",

					// Internal mocks
					internalMocksDir + "/layer5/layer6/mock_xyz/extra123.go",
					internalMocksDir + "/layer5/layer6/mock_xyz/nested",
					internalMocksDir + "/layer5/layer6/mock_xyz/nested/more.go",
					internalMocksDir + "/layer5/hello",
					internalMocksDir + "/layer5/hello/there.hi",
					internalMocksDir + "/garbage.txt",
				}

				for _, expectedPathToDelete := range expectedPathsToDelete {
					m.FSWrite.EXPECT().RemoveAll(expectedPathToDelete).Return(nil)
				}
			},
		},

		{
			Name: "when already tidy",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/some/pkg/qwerty",
							Interfaces: []string{"Iface2"},
						},
						{
							Path:       "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
							Interfaces: []string{"Iface3"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const primaryMocksDir = "/root/path/primary_mocks"
				m.FSWrite.EXPECT().ListRecursive(primaryMocksDir).
					Return([]string{
						primaryMocksDir + "/github.com",
						primaryMocksDir + "/github.com/some",
						primaryMocksDir + "/github.com/some/pkg",
						primaryMocksDir + "/github.com/some/pkg/mock_abc",
						primaryMocksDir + "/github.com/some/pkg/mock_abc/mock_abc.go",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/mock_qwerty.go",
					}, nil)

				const internalMocksDir = "/root/path/layer1/layer2/internal/layer3/layer4/internal/internal_mocks"
				m.FSWrite.EXPECT().ListRecursive(internalMocksDir).
					Return([]string{
						internalMocksDir + "/layer5",
						internalMocksDir + "/layer5/layer6",
						internalMocksDir + "/layer5/layer6/mock_xyz",
						internalMocksDir + "/layer5/layer6/mock_xyz/mock_xyz.go",
					}, nil)
			},
		},

		{
			Name:          "with invalid config: missing mock config",
			ExpectedError: mockgen.ErrMissingMockConfig,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
			},
		},

		{
			Name:          "with invalid config: internal package outside module",
			ExpectedError: mockgen.ErrInternalPackageOutsideModule,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/not/my/mod/internal/xyz",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},
		},

		{
			Name:          "when unable to list files recursively",
			ExpectedError: mockgen.ErrTidyUnableToList,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				m.FSWrite.EXPECT().ListRecursive("/root/path/primary_mocks").
					Return(nil, errors.New("you can't do that"))
			},
		},

		{
			Name:          "when unable to delete files",
			ExpectedError: mockgen.ErrTidyUnableToCleanup,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const primaryMocksDir = "/root/path/primary_mocks"
				m.FSWrite.EXPECT().ListRecursive(primaryMocksDir).
					Return([]string{
						primaryMocksDir + "/github.com",
						primaryMocksDir + "/github.com/mock_abc",
						primaryMocksDir + "/github.com/mock_abc/mock_abc.go",
						primaryMocksDir + "/github.com/extra1.go",
						primaryMocksDir + "/github.com/extra2.go",
					}, nil)

				m.FSWrite.EXPECT().RemoveAll(primaryMocksDir + "/github.com/extra1.go").Return(nil)
				m.FSWrite.EXPECT().RemoveAll(primaryMocksDir + "/github.com/extra2.go").Return(errors.New("oops"))
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]
		entry.Subject.Logger = log.New(ioutil.Discard, "", 0)

		err := entry.Subject.TidyMocks(entry.Config)
		ensure(err).IsError(entry.ExpectedError)
	})
}
