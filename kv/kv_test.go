package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	pairs := mapToPairs(map[string]string{
		"gonfig/fielda":        "bar",
		"gonfig/fieldb":        "1",
		"gonfig/fieldc":        "true",
		"gonfig/fieldd/0":      "one",
		"gonfig/fieldd/1":      "two",
		"gonfig/fielde":        "",
		"gonfig/fieldf/Test1":  "A",
		"gonfig/fieldf/Test2":  "B",
		"gonfig/fieldg/0/name": "A",
		"gonfig/fieldg/1/name": "B",
	})

	element := &sample{}

	err := Decode(pairs, element, "gonfig")
	require.NoError(t, err)

	expected := &sample{
		FieldA: "bar",
		FieldB: 1,
		FieldC: true,
		FieldD: []string{"one", "two"},
		FieldE: &struct {
			Name string
		}{},
		FieldF: map[string]string{
			"Test1": "A",
			"Test2": "B",
		},
		FieldG: []sub{
			{Name: "A"},
			{Name: "B"},
		},
	}
	assert.Equal(t, expected, element)
}

type sample struct {
	FieldA string
	FieldB int
	FieldC bool
	FieldD []string
	FieldE *struct {
		Name string
	} `label:"allowEmpty"`
	FieldF map[string]string
	FieldG []sub
}

type sub struct {
	Name string
}
