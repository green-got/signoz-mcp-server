package types

import (
	"encoding/json"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestDashboardInputSchemaUsesObjectItemsForArrayFields(t *testing.T) {
	reflector := jsonschema.Reflector{
		DoNotReference:            true,
		Anonymous:                 true,
		AllowAdditionalProperties: true,
	}

	schema := reflector.Reflect(Dashboard{})
	schema.Version = ""

	rawSchema, err := json.Marshal(schema)
	require.NoError(t, err)

	var schemaMap map[string]any
	require.NoError(t, json.Unmarshal(rawSchema, &schemaMap))

	assertArrayItemsAreObjectSchema(t, schemaMap, []string{
		"properties", "widgets", "items", "properties", "query", "properties", "builder",
		"properties", "queryData", "items", "properties", "functions", "items",
		"properties", "args",
	})

	assertArrayItemsAreObjectSchema(t, schemaMap, []string{
		"properties", "widgets", "items", "properties", "contextLinks", "properties", "linksData",
	})

	assertArrayItemsAreObjectSchema(t, schemaMap, []string{
		"properties", "widgets", "items", "properties", "query", "properties", "builder",
		"properties", "queryTraceOperator",
	})
}

func assertArrayItemsAreObjectSchema(t *testing.T, schemaMap map[string]any, path []string) {
	t.Helper()

	target := getSchemaPath(t, schemaMap, path)
	arraySchema, ok := target.(map[string]any)
	require.True(t, ok, "schema at path %v must be an object", path)
	require.Equal(t, "array", arraySchema["type"], "schema at path %v must be array", path)

	_, ok = arraySchema["items"].(map[string]any)
	require.True(t, ok, "array items at path %v must be an object schema", path)
}

func getSchemaPath(t *testing.T, schemaMap map[string]any, path []string) any {
	t.Helper()

	var current any = schemaMap
	for _, key := range path {
		obj, ok := current.(map[string]any)
		require.True(t, ok, "schema node before key %q must be an object", key)
		next, ok := obj[key]
		require.True(t, ok, "missing schema key %q", key)
		current = next
	}

	return current
}
