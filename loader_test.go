package gonfig

import (
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	example "github.com/crazy-max/gonfig/contrib/example/config"
	"github.com/crazy-max/gonfig/env"
	"github.com/crazy-max/gonfig/file"
	"github.com/crazy-max/gonfig/flag"
	"github.com/crazy-max/gonfig/parser"
	"github.com/crazy-max/gonfig/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type textDuration time.Duration

var errTextDuration = errors.New("invalid custom duration")

func (d *textDuration) UnmarshalText(text []byte) error {
	switch string(text) {
	case "30d":
		*d = textDuration(30 * 24 * time.Hour)
	case "2w":
		*d = textDuration(14 * 24 * time.Hour)
	case "":
		*d = 0
	default:
		return errTextDuration
	}
	return nil
}

type textName string

func (n *textName) UnmarshalText(text []byte) error {
	*n = textName(strings.ToUpper(string(text)))
	return nil
}

func TestTextUnmarshaler(t *testing.T) {
	type config struct {
		Duration    textDuration
		Pointer     *textDuration
		Absent      *textDuration
		Durations   []textDuration
		Pointers    []*textDuration
		Names       map[string]textName
		Times       []time.Time
		Time        time.Time
		TimePointer *time.Time
		IP          net.IP
		IPs         []net.IP
		Native      time.Duration
		Legacy      types.Duration
		Ignored     textDuration `label:"-" file:"-"`
	}
	for _, tc := range []struct {
		name   string
		decode func(*config) error
	}{
		{
			name: "environment",
			decode: func(c *config) error {
				return env.Decode([]string{
					"GONFIG_DURATION=30d", "GONFIG_POINTER=2w",
					"GONFIG_DURATIONS=30d,2w", "GONFIG_POINTERS=30d,2w",
					"GONFIG_NAMES_FIRST=hello", "GONFIG_TIMES=2026-01-01T00:00:00Z",
					"GONFIG_TIME=2026-01-01T00:00:00Z", "GONFIG_TIMEPOINTER=2026-01-01T00:00:00Z",
					"GONFIG_IP=192.0.2.1", "GONFIG_NATIVE=42", "GONFIG_LEGACY=42",
					"GONFIG_IGNORED=invalid",
					"GONFIG_IPS=192.0.2.1,192.0.2.2",
				}, env.DefaultNamePrefix, c)
			},
		},
		{
			name: "yaml",
			decode: func(c *config) error {
				return file.DecodeContent(`duration: 30d
pointer: 2w
durations: [30d, 2w]
pointers: [30d, 2w]
names:
  first: hello
times: ["2026-01-01T00:00:00Z"]
time: "2026-01-01T00:00:00Z"
timePointer: "2026-01-01T00:00:00Z"
ip: 192.0.2.1
ips: [192.0.2.1, 192.0.2.2]
native: 42
legacy: 42
ignored: invalid
`, ".yaml", c)
			},
		},
		{
			name: "toml",
			decode: func(c *config) error {
				return file.DecodeContent(`duration = "30d"
pointer = "2w"
durations = ["30d", "2w"]
pointers = ["30d", "2w"]
times = ["2026-01-01T00:00:00Z"]
time = "2026-01-01T00:00:00Z"
timePointer = "2026-01-01T00:00:00Z"
ip = "192.0.2.1"
ips = ["192.0.2.1", "192.0.2.2"]
native = 42
legacy = 42
ignored = "invalid"
[names]
first = "hello"
`, ".toml", c)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var c config
			require.NoError(t, tc.decode(&c))
			assert.Equal(t, textDuration(720*time.Hour), c.Duration)
			require.NotNil(t, c.Pointer)
			assert.Equal(t, textDuration(336*time.Hour), *c.Pointer)
			assert.Nil(t, c.Absent)
			assert.Equal(t, []textDuration{c.Duration, *c.Pointer}, c.Durations)
			require.Len(t, c.Pointers, 2)
			assert.Equal(t, &c.Duration, c.Pointers[0])
			assert.Equal(t, c.Pointer, c.Pointers[1])
			assert.Equal(t, map[string]textName{"first": "HELLO"}, c.Names)
			wantTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			assert.Equal(t, wantTime, c.Time)
			assert.Equal(t, &wantTime, c.TimePointer)
			assert.Equal(t, []time.Time{wantTime}, c.Times)
			assert.Equal(t, net.ParseIP("192.0.2.1"), c.IP)
			assert.Equal(t, []net.IP{net.ParseIP("192.0.2.1"), net.ParseIP("192.0.2.2")}, c.IPs)
			assert.Equal(t, time.Duration(42), c.Native)
			assert.Equal(t, types.Duration(42*time.Second), c.Legacy)
			assert.Zero(t, c.Ignored)
		})
	}
}

type textAddress struct {
	Host string
}

func (a *textAddress) UnmarshalText(text []byte) error {
	a.Host = strings.ToUpper(string(text))
	return nil
}

func TestTextUnmarshalerStructuredInput(t *testing.T) {
	type config struct {
		Address   textAddress
		Pointer   *textAddress
		Addresses []textAddress
	}
	for _, tc := range []struct {
		name   string
		decode func(*config) error
	}{
		{
			name: "environment",
			decode: func(c *config) error {
				return env.Decode([]string{"GONFIG_ADDRESS_HOST=example.org", "GONFIG_POINTER_HOST=example.org"}, env.DefaultNamePrefix, c)
			},
		},
		{
			name: "yaml",
			decode: func(c *config) error {
				return file.DecodeContent("address:\n  host: example.org\npointer:\n  host: example.org\naddresses:\n  - host: example.org\n", ".yaml", c)
			},
		},
		{
			name: "toml",
			decode: func(c *config) error {
				return file.DecodeContent("[address]\nhost = 'example.org'\n[pointer]\nhost = 'example.org'\n[[addresses]]\nhost = 'example.org'\n", ".toml", c)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var c config
			require.NoError(t, tc.decode(&c))
			assert.Equal(t, textAddress{Host: "example.org"}, c.Address)
			assert.Equal(t, &c.Address, c.Pointer)
			if tc.name != "environment" {
				assert.Equal(t, []textAddress{c.Address}, c.Addresses)
			}
		})
	}

	c := struct{ Address textAddress }{Address: textAddress{Host: "example.org"}}
	environment, err := env.Encode(env.DefaultNamePrefix, &c)
	require.NoError(t, err)
	assert.Contains(t, environment, parser.Flat{Name: "GONFIG_ADDRESS_HOST", Default: "example.org"})
	flags, err := flag.Encode(&c)
	require.NoError(t, err)
	assert.Contains(t, flags, parser.Flat{Name: "address.host", Default: "example.org"})

	require.NoError(t, env.Decode([]string{"GONFIG_ADDRESS=example.org"}, env.DefaultNamePrefix, &c))
	assert.Equal(t, "EXAMPLE.ORG", c.Address.Host)
}

func TestTextUnmarshalerErrors(t *testing.T) {
	for _, field := range []string{"duration", "pointer", "durations", "values.entry"} {
		t.Run(field, func(t *testing.T) {
			var c struct {
				Duration  textDuration
				Pointer   *textDuration
				Durations []textDuration
				Values    map[string]textDuration
			}
			key := "GONFIG_" + strings.ToUpper(strings.ReplaceAll(field, ".", "_"))
			err := env.Decode([]string{key + "=invalid"}, env.DefaultNamePrefix, &c)
			require.ErrorIs(t, err, errTextDuration)
		})
	}
	var c struct{ Duration textDuration }
	c.Duration = 1
	require.NoError(t, env.Decode([]string{"GONFIG_DURATION="}, env.DefaultNamePrefix, &c))
	assert.Zero(t, c.Duration)
}

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
