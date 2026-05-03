package render

import (
	"encoding/json"
	"strings"
	"testing"
)

func makeAlert(level, alertType string, resolved bool, data string) Alert {
	var a Alert
	a.Level = level
	a.Resolved = resolved
	a.Target.Type = "Server"
	a.Data.Type = alertType
	if err := json.Unmarshal([]byte(data), &a.Data.Data); err != nil {
		panic(err)
	}
	return a
}

func TestRenderServerUnreachable(t *testing.T) {
	a := makeAlert("CRITICAL", "ServerUnreachable", false, `{
		"id": "69f484d03be021935cd04721",
		"name": "M14",
		"region": null,
		"err": {
			"error": "Timed out waiting for Ping",
			"trace": "deadline has elapsed"
		}
	}`)
	out := Render(a)
	if !strings.Contains(out, "🔴") {
		t.Error("expected critical emoji")
	}
	if !strings.Contains(out, "<b>CRITICAL</b>") {
		t.Error("expected bold level")
	}
	if !strings.Contains(out, "M14") {
		t.Error("expected server name")
	}
	if !strings.Contains(out, "Timed out waiting for Ping") {
		t.Error("expected error message")
	}
	if !strings.Contains(out, "❌") {
		t.Error("expected unresolved icon")
	}
}

func TestRenderStackStateChange(t *testing.T) {
	a := makeAlert("WARNING", "StackStateChange", true, `{
		"id": "69c910d0e7509802222b72a7",
		"name": "redis",
		"server_id": "69f49a30acb6f25705157c46",
		"server_name": "Hetzner",
		"from": "stopped",
		"to": "running"
	}`)
	out := Render(a)
	if !strings.Contains(out, "⚠️") {
		t.Error("expected warning emoji")
	}
	if !strings.Contains(out, "redis") {
		t.Error("expected stack name")
	}
	if !strings.Contains(out, "Hetzner") {
		t.Error("expected server name")
	}
	if !strings.Contains(out, "stopped") || !strings.Contains(out, "running") {
		t.Error("expected state transition")
	}
	if !strings.Contains(out, "✅") {
		t.Error("expected resolved icon")
	}
}

func TestRenderFallback(t *testing.T) {
	a := makeAlert("INFO", "SomeUnknownType", false, `{"name":"test","extra":"value"}`)
	out := Render(a)
	if !strings.Contains(out, "ℹ️") {
		t.Error("expected info emoji")
	}
	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRenderOKLevel(t *testing.T) {
	a := makeAlert("OK", "ServerUnreachable", true, `{"name":"myserver"}`)
	out := Render(a)
	if !strings.Contains(out, "✅") {
		t.Error("expected OK emoji")
	}
}

func TestHTMLEscape(t *testing.T) {
	a := makeAlert("INFO", "ServerUnreachable", false, `{
		"name": "<evil>&name",
		"err": {"error": "bad & broken", "trace": ""}
	}`)
	out := Render(a)
	if strings.Contains(out, "<evil>") {
		t.Error("HTML injection not escaped")
	}
	if !strings.Contains(out, "&lt;evil&gt;") {
		t.Error("expected escaped angle brackets")
	}
}

func TestMaxLength(t *testing.T) {
	// build a payload that will produce a very long message
	long := strings.Repeat("x", 5000)
	a := makeAlert("INFO", "SomeUnknownType", false, `{"name":"`+long+`"}`)
	out := Render(a)
	if len(out) > maxLen {
		t.Errorf("message length %d exceeds %d", len(out), maxLen)
	}
}
