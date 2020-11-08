package env

import (
	"testing"

	"github.com/crazy-max/gonfig/generator"
	"github.com/crazy-max/gonfig/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		desc     string
		environ  []string
		element  interface{}
		expected interface{}
	}{
		{
			desc:     "no env vars",
			environ:  nil,
			expected: nil,
		},
		{
			desc:    "bool value",
			environ: []string{"GONFIG_FOO=true"},
			element: &struct {
				Foo bool
			}{},
			expected: &struct {
				Foo bool
			}{
				Foo: true,
			},
		},
		{
			desc:    "equal",
			environ: []string{"GONFIG_FOO=bar"},
			element: &struct {
				Foo string
			}{},
			expected: &struct {
				Foo string
			}{
				Foo: "bar",
			},
		},
		{
			desc:    "multiple bool flags without value",
			environ: []string{"GONFIG_FOO=true", "GONFIG_BAR=true"},
			element: &struct {
				Foo bool
				Bar bool
			}{},
			expected: &struct {
				Foo bool
				Bar bool
			}{
				Foo: true,
				Bar: true,
			},
		},
		{
			desc:    "map string",
			environ: []string{"GONFIG_FOO_NAME=bar"},
			element: &struct {
				Foo map[string]string
			}{},
			expected: &struct {
				Foo map[string]string
			}{
				Foo: map[string]string{
					"name": "bar",
				},
			},
		},
		{
			desc:    "map struct",
			environ: []string{"GONFIG_FOO_NAME_VALUE=bar"},
			element: &struct {
				Foo map[string]struct{ Value string }
			}{},
			expected: &struct {
				Foo map[string]struct{ Value string }
			}{
				Foo: map[string]struct{ Value string }{
					"name": {
						Value: "bar",
					},
				},
			},
		},
		{
			desc:    "map struct with sub-struct",
			environ: []string{"GONFIG_FOO_NAME_BAR_VALUE=bar"},
			element: &struct {
				Foo map[string]struct {
					Bar *struct{ Value string }
				}
			}{},
			expected: &struct {
				Foo map[string]struct {
					Bar *struct{ Value string }
				}
			}{
				Foo: map[string]struct {
					Bar *struct{ Value string }
				}{
					"name": {
						Bar: &struct {
							Value string
						}{
							Value: "bar",
						},
					},
				},
			},
		},
		{
			desc:    "map struct with sub-map",
			environ: []string{"GONFIG_FOO_NAME1_BAR_NAME2_VALUE=bar"},
			element: &struct {
				Foo map[string]struct {
					Bar map[string]struct{ Value string }
				}
			}{},
			expected: &struct {
				Foo map[string]struct {
					Bar map[string]struct{ Value string }
				}
			}{
				Foo: map[string]struct {
					Bar map[string]struct{ Value string }
				}{
					"name1": {
						Bar: map[string]struct{ Value string }{
							"name2": {
								Value: "bar",
							},
						},
					},
				},
			},
		},
		{
			desc:    "slice",
			environ: []string{"GONFIG_FOO=bar,baz"},
			element: &struct {
				Foo []string
			}{},
			expected: &struct {
				Foo []string
			}{
				Foo: []string{"bar", "baz"},
			},
		},
		{
			desc:    "struct pointer value",
			environ: []string{"GONFIG_FOO=true"},
			element: &struct {
				Foo *struct{ Field string } `label:"allowEmpty"`
			}{},
			expected: &struct {
				Foo *struct{ Field string } `label:"allowEmpty"`
			}{
				Foo: &struct{ Field string }{},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := Decode(test.environ, DefaultNamePrefix, test.element)
			require.NoError(t, err)

			assert.Equal(t, test.expected, test.element)
		})
	}
}

func TestEncode(t *testing.T) {
	element := &Ya{
		Foo: &Yaa{
			FieldIn1: "bar",
			FieldIn2: false,
			FieldIn3: 1,
			FieldIn4: map[string]string{
				parser.MapNamePlaceholder: "",
			},
			FieldIn5: map[string]int{
				parser.MapNamePlaceholder: 0,
			},
			FieldIn6: map[string]struct{ Field string }{
				parser.MapNamePlaceholder: {},
			},
			FieldIn7: map[string]struct{ Field map[string]string }{
				parser.MapNamePlaceholder: {
					Field: map[string]string{
						parser.MapNamePlaceholder: "",
					},
				},
			},
			FieldIn8: map[string]*struct{ Field string }{
				parser.MapNamePlaceholder: {},
			},
			FieldIn9: map[string]*struct{ Field map[string]string }{
				parser.MapNamePlaceholder: {
					Field: map[string]string{
						parser.MapNamePlaceholder: "",
					},
				},
			},
			FieldIn10: struct{ Field string }{},
			FieldIn11: &struct{ Field string }{},
			FieldIn12: func(v string) *string { return &v }(""),
			FieldIn13: func(v bool) *bool { return &v }(false),
			FieldIn14: func(v int) *int { return &v }(0),
		},
		Field1: "bir",
		Field2: true,
		Field3: 0,
		Field4: map[string]string{
			parser.MapNamePlaceholder: "",
		},
		Field5: map[string]int{
			parser.MapNamePlaceholder: 0,
		},
		Field6: map[string]struct{ Field string }{
			parser.MapNamePlaceholder: {},
		},
		Field7: map[string]struct{ Field map[string]string }{
			parser.MapNamePlaceholder: {
				Field: map[string]string{
					parser.MapNamePlaceholder: "",
				},
			},
		},
		Field8: map[string]*struct{ Field string }{
			parser.MapNamePlaceholder: {},
		},
		Field9: map[string]*struct{ Field map[string]string }{
			parser.MapNamePlaceholder: {
				Field: map[string]string{
					parser.MapNamePlaceholder: "",
				},
			},
		},
		Field10: struct{ Field string }{},
		Field11: &struct{ Field string }{},
		Field12: func(v string) *string { return &v }(""),
		Field13: func(v bool) *bool { return &v }(false),
		Field14: func(v int) *int { return &v }(0),
		Field15: []int{7},
	}
	generator.Generate(element)

	flats, err := Encode(DefaultNamePrefix, element)
	require.NoError(t, err)

	expected := []parser.Flat{
		{
			Name:        "GONFIG_FIELD1",
			Description: "",
			Default:     "bir",
		},
		{
			Name:        "GONFIG_FIELD10",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD10_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD11_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD12",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD13",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FIELD14",
			Description: "",
			Default:     "0",
		},
		{
			Name:        "GONFIG_FIELD15",
			Description: "",
			Default:     "7",
		},
		{
			Name:        "GONFIG_FIELD2",
			Description: "",
			Default:     "true",
		},
		{
			Name:        "GONFIG_FIELD3",
			Description: "",
			Default:     "0",
		},
		{
			Name:        "GONFIG_FIELD4_\u003cNAME\u003e",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD5_\u003cNAME\u003e",
			Description: "",
			Default:     "0",
		},
		{
			Name:        "GONFIG_FIELD6_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FIELD6_\u003cNAME\u003e_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD7_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FIELD7_\u003cNAME\u003e_FIELD_\u003cNAME\u003e",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD8_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FIELD8_\u003cNAME\u003e_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FIELD9_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FIELD9_\u003cNAME\u003e_FIELD_\u003cNAME\u003e",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN1",
			Description: "",
			Default:     "bar",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN10",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN10_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN11_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN12",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN13",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN14",
			Description: "",
			Default:     "0",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN2",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN3",
			Description: "",
			Default:     "1",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN4_\u003cNAME\u003e",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN5_\u003cNAME\u003e",
			Description: "",
			Default:     "0",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN6_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN6_\u003cNAME\u003e_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN7_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN7_\u003cNAME\u003e_FIELD_\u003cNAME\u003e",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN8_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN8_\u003cNAME\u003e_FIELD",
			Description: "",
			Default:     "",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN9_\u003cNAME\u003e",
			Description: "",
			Default:     "false",
		},
		{
			Name:        "GONFIG_FOO_FIELDIN9_\u003cNAME\u003e_FIELD_\u003cNAME\u003e",
			Description: "",
			Default:     "",
		},
	}

	assert.Equal(t, expected, flats)
}
