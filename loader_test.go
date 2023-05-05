package gonfig

import (
	"strings"
	"testing"
	"time"

	example "github.com/crazy-max/gonfig/_example/config"
	"github.com/crazy-max/gonfig/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvLoader(t *testing.T) {
	testCases := []struct {
		desc     string
		cfgfile  string
		environ  []string
		found    bool
		expected interface{}
		wantErr  bool
	}{
		{
			desc:     "no env vars",
			environ:  nil,
			found:    false,
			expected: example.Config{},
			wantErr:  false,
		},
		{
			desc: "ftp server",
			environ: []string{
				env.DefaultNamePrefix + "SERVER_FTP_HOST=test.rebex.net",
				env.DefaultNamePrefix + "SERVER_FTP_USERNAME=demo",
				env.DefaultNamePrefix + "SERVER_FTP_PASSWORD=password",
				env.DefaultNamePrefix + "SERVER_FTP_SOURCES=/",
			},
			found: true,
			expected: example.Config{
				Server: &example.Server{
					FTP: &example.ServerFTP{
						Host:     "test.rebex.net",
						Port:     21,
						Username: "demo",
						Password: "password",
						Sources: []string{
							"/",
						},
						Timeout:            example.NewDuration(5 * time.Second),
						DisableUTF8:        example.NewFalse(),
						DisableEPSV:        example.NewFalse(),
						DisableMLSD:        example.NewFalse(),
						EscapeRegexpMeta:   example.NewFalse(),
						TLS:                example.NewFalse(),
						InsecureSkipVerify: example.NewFalse(),
						LogTrace:           example.NewFalse(),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.environ != nil {
				for _, environ := range tt.environ {
					n := strings.SplitN(environ, "=", 2)
					t.Setenv(n[0], n[1])
				}
			}

			var cfg example.Config
			envLoader := NewEnvLoader(EnvLoaderConfig{
				Prefix: env.DefaultNamePrefix,
			})

			found, err := envLoader.Load(&cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestFileLoader(t *testing.T) {
	cases := []struct {
		name     string
		cfgfile  string
		found    bool
		expected *example.Config
		wantErr  bool
	}{
		{
			name:     "Non-existing file",
			cfgfile:  "",
			found:    false,
			expected: &example.Config{},
		},
		{
			name:    "Fail on wrong file format",
			cfgfile: "./fixtures/config.invalid.yml",
			found:   true,
			wantErr: true,
		},
		{
			name:    "Success",
			cfgfile: "./fixtures/config.test.yml",
			found:   true,
			expected: &example.Config{
				Server: &example.Server{
					FTP: &example.ServerFTP{
						Host:     "test.rebex.net",
						Port:     21,
						Username: "demo",
						Password: "password",
						Sources: []string{
							"/",
						},
						Timeout:            example.NewDuration(5 * time.Second),
						DisableUTF8:        example.NewTrue(),
						DisableEPSV:        example.NewFalse(),
						DisableMLSD:        example.NewFalse(),
						EscapeRegexpMeta:   example.NewFalse(),
						TLS:                example.NewFalse(),
						InsecureSkipVerify: example.NewFalse(),
						LogTrace:           example.NewFalse(),
					},
				},
				Download: &example.Download{
					Output:        "./fixtures/downloads",
					UID:           1000,
					GID:           1000,
					ChmodFile:     0o644,
					ChmodDir:      0o755,
					Include:       []string{`^Foo\.Bar\.S01.+(VOSTFR|SUBFRENCH).+(720p).+(HDTV|WEB-DL|WEBRip).+`},
					Exclude:       []string{`\.nfo$`},
					Since:         "2019-02-01T18:50:05Z",
					Retry:         3,
					HideSkipped:   example.NewFalse(),
					TempFirst:     example.NewFalse(),
					CreateBaseDir: example.NewFalse(),
				},
				Notif: &example.Notif{
					Mail: &example.NotifMail{
						Host:               "smtp.example.com",
						Port:               587,
						SSL:                example.NewFalse(),
						InsecureSkipVerify: example.NewFalse(),
						From:               "from@example.com",
						To:                 "webmaster@example.com",
					},
					Webhook: &example.NotifWebhook{
						Endpoint: "http://webhook.foo.com/sd54qad89azd5a",
						Method:   "GET",
						Headers: map[string]string{
							"content-type":  "application/json",
							"authorization": "Token123456",
						},
						Timeout: example.NewDuration(10 * time.Second),
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			fileLoader := NewFileLoader(FileLoaderConfig{
				Filename: tt.cfgfile,
			})

			cfg := &example.Config{}
			found, err := fileLoader.Load(cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestFlagLoader(t *testing.T) {
	testCases := []struct {
		desc     string
		cfgfile  string
		args     []string
		found    bool
		expected interface{}
		wantErr  bool
	}{
		{
			desc:     "no flag arguments",
			args:     nil,
			found:    false,
			expected: example.Config{},
			wantErr:  false,
		},
		{
			desc: "ftp server",
			args: []string{
				"--server.ftp.host=test.rebex.net",
				"--server.ftp.username=demo",
				"--server.ftp.password=password",
				"--server.ftp.sources=/src1,/src2",
				"--server.ftp.disableEPSV=true",
			},
			found: true,
			expected: example.Config{
				Server: &example.Server{
					FTP: &example.ServerFTP{
						Host:     "test.rebex.net",
						Port:     21,
						Username: "demo",
						Password: "password",
						Sources: []string{
							"/src1",
							"/src2",
						},
						Timeout:            example.NewDuration(5 * time.Second),
						DisableUTF8:        example.NewFalse(),
						DisableEPSV:        example.NewTrue(),
						DisableMLSD:        example.NewFalse(),
						EscapeRegexpMeta:   example.NewFalse(),
						TLS:                example.NewFalse(),
						InsecureSkipVerify: example.NewFalse(),
						LogTrace:           example.NewFalse(),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			var cfg example.Config
			flagLoader := NewFlagLoader(FlagLoaderConfig{
				Args: tt.args,
			})

			found, err := flagLoader.Load(&cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
