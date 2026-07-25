package main

import (
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetArgs(t *testing.T) {
	t.Run("valid arguments", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"name": "Bukhari",
				},
			},
		}
		args, errResult := getArgs(req)
		if errResult != nil {
			t.Fatalf("expected no error, got %v", errResult)
		}
		if args == nil {
			t.Fatal("expected non-nil args")
		}
	})

	t.Run("nil arguments", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: nil,
			},
		}
		_, errResult := getArgs(req)
		if errResult == nil {
			t.Fatal("expected error for nil arguments")
		}
	})

	t.Run("wrong type arguments", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: "not a map",
			},
		}
		_, errResult := getArgs(req)
		if errResult == nil {
			t.Fatal("expected error for wrong type arguments")
		}
	})
}

func TestGetString(t *testing.T) {
	args := handlerArgs{
		"name":   "Bukhari",
		"count":  float64(42),
		"nested": map[string]interface{}{"key": "val"},
	}

	t.Run("existing string key", func(t *testing.T) {
		v, err := args.getString("name")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if v != "Bukhari" {
			t.Fatalf("expected 'Bukhari', got '%s'", v)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		_, err := args.getString("missing")
		if err == nil {
			t.Fatal("expected error for missing key")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		_, err := args.getString("count")
		if err == nil {
			t.Fatal("expected error for non-string value")
		}
	})

	t.Run("nil value", func(t *testing.T) {
		args := handlerArgs{"x": nil}
		_, err := args.getString("x")
		if err == nil {
			t.Fatal("expected error for nil value")
		}
	})
}

func TestGetFloat(t *testing.T) {
	args := handlerArgs{
		"count": float64(42),
		"pi":    float64(3.14),
		"name":  "Bukhari",
	}

	t.Run("existing float key", func(t *testing.T) {
		v, err := args.getFloat("count")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if v != float64(42) {
			t.Fatalf("expected 42, got %f", v)
		}
	})

	t.Run("float with decimal", func(t *testing.T) {
		v, err := args.getFloat("pi")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if v != float64(3.14) {
			t.Fatalf("expected 3.14, got %f", v)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		_, err := args.getFloat("missing")
		if err == nil {
			t.Fatal("expected error for missing key")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		_, err := args.getFloat("name")
		if err == nil {
			t.Fatal("expected error for non-float value")
		}
	})
}

func TestGetOptionalFloat(t *testing.T) {
	args := handlerArgs{
		"limit":  float64(50),
		"zero":   float64(0),
		"neg":    float64(-5),
		"name":   "Bukhari",
	}

	t.Run("existing positive value", func(t *testing.T) {
		v := args.getOptionalFloat("limit", 20)
		if v != float64(50) {
			t.Fatalf("expected 50, got %f", v)
		}
	})

	t.Run("missing key returns default", func(t *testing.T) {
		v := args.getOptionalFloat("missing", 10)
		if v != float64(10) {
			t.Fatalf("expected 10, got %f", v)
		}
	})

	t.Run("zero value returns default", func(t *testing.T) {
		v := args.getOptionalFloat("zero", 20)
		if v != float64(20) {
			t.Fatalf("expected 20 (default), got %f", v)
		}
	})

	t.Run("negative value returns default", func(t *testing.T) {
		v := args.getOptionalFloat("neg", 20)
		if v != float64(20) {
			t.Fatalf("expected 20 (default), got %f", v)
		}
	})

	t.Run("wrong type returns default", func(t *testing.T) {
		v := args.getOptionalFloat("name", 7)
		if v != float64(7) {
			t.Fatalf("expected 7 (default), got %f", v)
		}
	})
}

func TestMarshalResult(t *testing.T) {
	t.Run("marshal struct", func(t *testing.T) {
		result := marshalResult(map[string]string{"key": "value"})
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		// Verify it's a text result with valid JSON
		tc, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatal("expected TextContent")
		}
		var data map[string]string
		if err := json.Unmarshal([]byte(tc.Text), &data); err != nil {
			t.Fatalf("expected valid JSON, got error: %v", err)
		}
		if data["key"] != "value" {
			t.Fatalf("expected 'value', got '%s'", data["key"])
		}
	})

	t.Run("marshal nil", func(t *testing.T) {
		result := marshalResult(nil)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
	})

	t.Run("marshal unencodable", func(t *testing.T) {
		// channel can't be marshalled
		result := marshalResult(make(chan int))
		if result == nil {
			t.Fatal("expected non-nil result even on error")
		}
		if result.IsError != true {
			t.Fatal("expected error result for unencodable value")
		}
	})
}
