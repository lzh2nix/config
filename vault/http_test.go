// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package vault

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

type testHandler struct{}

type testObject struct {
	Field string `json:"field"`
}

func TestUnmarshalJSON(t *testing.T) {
	input := `{"test":"data","validation":"process"}`
	expected := map[string]interface{}{
		"test":       "data",
		"validation": "process",
	}

	var actual map[string]interface{}
	err := unmarshalJSON(bytes.NewReader([]byte(input)), &actual)
	if err != nil {
		fmt.Printf("decoding err: %v\n", err)
	}

	assert.EqualValues(t, expected, actual)
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(&testObject{
		Field: "response",
	})
	w.Write(b)
}

func TestPostGet(t *testing.T) {
	s := httptest.NewServer(new(testHandler))
	defer s.Close()
	body := testObject{Field: "hello"}
	expect := &testObject{Field: "response"}

	{
		output := new(testObject)
		err := httpPost(s.URL, body, output)
		assert.NoError(t, err)
		assert.EqualValues(t, expect, output)
	}

	{
		output := new(testObject)
		err := httpGet(s.URL, output)
		assert.NoError(t, err)
		assert.EqualValues(t, expect, output)
	}
}
