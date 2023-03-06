package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
