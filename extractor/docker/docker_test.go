package docker

import (
	"context"
	"sort"
	"testing"

	godeptypes "github.com/aquasecurity/go-dep-parser/pkg/types"
	digest "github.com/opencontainers/go-digest"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/fanal/extractor"
	"github.com/aquasecurity/fanal/extractor/image"
	"github.com/aquasecurity/fanal/types"
)

func TestApplyLayers(t *testing.T) {
	testCases := []struct {
		name                string
		inputLayers         []types.LayerInfo
		expectedImageDetail types.ImageDetail
	}{
		{
			name: "happy path",
			inputLayers: []types.LayerInfo{
				{
					SchemaVersion: 1,
					Digest:        "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					OS: &types.OS{
						Family: "alpine",
						Name:   "3.10",
					},
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "lib/apk/db/installed",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
							},
						},
					},
					Applications: []types.Application{
						{
							Type:     "gem",
							FilePath: "app/Gemfile.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "gemlibrary1",
										Version: "1.2.3",
									},
								},
							},
						},
						{
							Type:     "composer",
							FilePath: "app/composer.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
					DiffID:        "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "lib/apk/db/installed",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
								{ // added
									Name:    "musl",
									Version: "1.2.4",
									Release: "4.5.7",
								},
							},
						},
					},
					WhiteoutFiles: []string{"app/composer.lock"},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "lib/apk/db/installed",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
								{
									Name:    "musl",
									Version: "1.2.4",
									Release: "4.5.8", // updated
								},
							},
						},
					},
				},
			},
			expectedImageDetail: types.ImageDetail{
				OS: &types.OS{
					Family: "alpine",
					Name:   "3.10",
				},
				Packages: []types.Package{
					{
						Name:    "musl",
						Version: "1.2.4",
						Release: "4.5.8",
						Layer: types.Layer{
							Digest: "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
							DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
						},
					},
					{
						Name:    "openssl",
						Version: "1.2.3",
						Release: "4.5.6",
						Layer: types.Layer{
							Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
							DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
						},
					},
				},
				Applications: []types.Application{
					{
						Type:     "gem",
						FilePath: "app/Gemfile.lock",
						Libraries: []types.LibraryInfo{
							{
								Library: godeptypes.Library{
									Name:    "gemlibrary1",
									Version: "1.2.3",
								},
								Layer: types.Layer{
									Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
									DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "happy path with removed and updated lockfile",
			inputLayers: []types.LayerInfo{
				{
					SchemaVersion: 1,
					Digest:        "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					OS: &types.OS{
						Family: "alpine",
						Name:   "3.10",
					},
					Applications: []types.Application{
						{
							Type:     "gem",
							FilePath: "app/Gemfile.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "rails",
										Version: "5.0.0",
									},
								},
								{
									Library: godeptypes.Library{
										Name:    "rack",
										Version: "4.0.0",
									},
								},
							},
						},
						{
							Type:     "composer",
							FilePath: "app/composer.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
					DiffID:        "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
					Applications: []types.Application{
						{
							Type:     "gem",
							FilePath: "app/Gemfile.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "rails",
										Version: "6.0.0",
									},
								},
								{
									Library: godeptypes.Library{
										Name:    "rack",
										Version: "4.0.0",
									},
								},
							},
						},
						{
							Type:     "composer",
							FilePath: "app/composer2.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
					WhiteoutFiles: []string{"app/composer.lock"},
				},
			},
			expectedImageDetail: types.ImageDetail{
				OS: &types.OS{
					Family: "alpine",
					Name:   "3.10",
				},
				Applications: []types.Application{
					{
						Type:     "gem",
						FilePath: "app/Gemfile.lock",
						Libraries: []types.LibraryInfo{
							{
								Library: godeptypes.Library{
									Name:    "rack",
									Version: "4.0.0",
								},
								Layer: types.Layer{
									Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
									DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
								},
							},
							{
								Library: godeptypes.Library{
									Name:    "rails",
									Version: "6.0.0",
								},
								Layer: types.Layer{
									Digest: "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
									DiffID: "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
								},
							},
						},
					},
					{
						Type:     "composer",
						FilePath: "app/composer2.lock",
						Libraries: []types.LibraryInfo{
							{
								Library: godeptypes.Library{
									Name:    "phplibrary1",
									Version: "6.6.6",
								},
								Layer: types.Layer{
									Digest: "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
									DiffID: "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "happy path with status.d",
			inputLayers: []types.LayerInfo{
				{
					SchemaVersion: 1,
					Digest:        "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
					DiffID:        "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
					OS: &types.OS{
						Family: "debian",
						Name:   "8",
					},
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "var/lib/dpkg/status.d/openssl",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
							},
						},
					},
					Applications: []types.Application{
						{
							Type:     "composer",
							FilePath: "app/composer.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "var/lib/dpkg/status.d/libc",
							Packages: []types.Package{
								{
									Name:    "libc",
									Version: "1.2.4",
									Release: "4.5.7",
								},
							},
						},
					},
					OpaqueDirs: []string{"app"},
				},
			},
			expectedImageDetail: types.ImageDetail{
				OS: &types.OS{
					Family: "debian",
					Name:   "8",
				},
				Packages: []types.Package{
					{
						Name:    "libc",
						Version: "1.2.4",
						Release: "4.5.7",
						Layer: types.Layer{
							Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
							DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
						},
					},
					{
						Name:    "openssl",
						Version: "1.2.3",
						Release: "4.5.6",
						Layer: types.Layer{
							Digest: "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
							DiffID: "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotImageDetail := ApplyLayers(tc.inputLayers)
			sort.Slice(gotImageDetail.Packages, func(i, j int) bool {
				return gotImageDetail.Packages[i].Name < gotImageDetail.Packages[j].Name
			})
			sort.Slice(gotImageDetail.Applications, func(i, j int) bool {
				return gotImageDetail.Applications[i].FilePath < gotImageDetail.Applications[j].FilePath
			})
			for _, app := range gotImageDetail.Applications {
				sort.Slice(app.Libraries, func(i, j int) bool {
					return app.Libraries[i].Library.Name < app.Libraries[j].Library.Name
				})
			}
			assert.Equal(t, tc.expectedImageDetail, gotImageDetail, tc.name)
		})
	}
}

func TestExtractor_ExtractLayerFiles(t *testing.T) {
	type fields struct {
		option types.DockerOption
		image  image.RealImage
	}
	type args struct {
		ctx       context.Context
		dig       digest.Digest
		filenames []string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		imagePath       string
		expectedDigest  digest.Digest
		expectedFileMap extractor.FileMap
		expectedOpqDirs []string
		expectedWhFiles []string
		wantErr         string
	}{
		{
			name:      "happy path",
			imagePath: "testdata/image1.tar",
			args: args{
				ctx:       nil,
				dig:       "sha256:d9ff549177a94a413c425ffe14ae1cc0aa254bc9c7df781add08e7d2fba25d27",
				filenames: []string{"etc/hostname"},
			},
			expectedDigest: "sha256:d9ff549177a94a413c425ffe14ae1cc0aa254bc9c7df781add08e7d2fba25d27",
			expectedFileMap: extractor.FileMap{
				"etc/hostname": []byte("localhost\n"),
			},
		},
		{
			name:      "opq file path",
			imagePath: "testdata/image1.tar",
			args: args{
				ctx:       nil,
				dig:       "sha256:a8b87ccf2f2f94b9e23308560800afa3f272aa6db5cc7d9b0119b6843889cff2",
				filenames: []string{"etc/test/"},
			},
			expectedDigest: "sha256:a8b87ccf2f2f94b9e23308560800afa3f272aa6db5cc7d9b0119b6843889cff2",
			expectedFileMap: extractor.FileMap{
				"etc/test/bar": []byte("bar\n"),
			},
			expectedOpqDirs: []string{"etc/test/"},
			expectedWhFiles: []string{"var/foo"},
		},
		{
			name:      "sad path with GetLayer fails",
			imagePath: "testdata/image1.tar",
			args: args{
				ctx:       nil,
				dig:       "sha256:unknown",
				filenames: []string{"var/foo"},
			},
			expectedDigest: "sha256:f75441026d68038ca80e92f342fb8f3c0f1faeec67b5a80c98f033a65beaef5a",
			expectedFileMap: extractor.FileMap{
				"var/foo": []byte(""),
			},
			wantErr: "Unknown blob",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, cleanup, err := NewDockerArchiveExtractor(context.Background(), tt.imagePath, types.DockerOption{})
			require.NoError(t, err)
			defer cleanup()

			actualDigest, actualFileMap, actualOpqDirs, actualWhFiles, err := d.ExtractLayerFiles(tt.args.ctx, tt.args.dig, tt.args.filenames)
			if tt.wantErr != "" {
				require.NotNil(t, err, tt.name)
				assert.Contains(t, err.Error(), tt.wantErr, tt.name)
				return
			} else {
				require.NoError(t, err, tt.name)
			}

			assert.Equal(t, tt.expectedDigest, actualDigest)
			assert.Equal(t, tt.expectedFileMap, actualFileMap)
			assert.Equal(t, tt.expectedOpqDirs, actualOpqDirs)
			assert.Equal(t, tt.expectedWhFiles, actualWhFiles)
		})
	}
}
