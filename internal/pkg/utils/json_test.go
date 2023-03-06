package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonString(t *testing.T) {
	testData := struct {
		FieldName string
	}{
		FieldName: "jack",
	}
	resp := ToJSON(testData)
	assert.Equal(t, `{"FieldName":"jack"}`, resp)
}
