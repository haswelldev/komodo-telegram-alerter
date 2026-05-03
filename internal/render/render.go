package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

// Alert mirrors Komodo's alert webhook payload.
type Alert struct {
	Level    string `json:"level"`
	Resolved bool   `json:"resolved"`
	Target   struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"target"`
	Data struct {
		Type string         `json:"type"`
		Data map[string]any `json:"data"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

const maxLen = 4096

var (
	registry map[string]*template.Template
	fallback *template.Template
	funcMap  template.FuncMap
)

func init() {
	funcMap = template.FuncMap{
		"esc":          htmlEsc,
		"bold":         func(s string) string { return "<b>" + htmlEsc(s) + "</b>" },
		"italic":       func(s string) string { return "<i>" + htmlEsc(s) + "</i>" },
		"code":         func(s string) string { return "<code>" + htmlEsc(s) + "</code>" },
		"emoji":        levelEmoji,
		"resolvedIcon": resolvedIcon,
		"json":         jsonBlock,
		"ts":           tsFormat,
		"get":          safeGet,
		"str":          anyToStr,
	}

	registry = make(map[string]*template.Template)
	for typ, src := range defaultTemplates {
		t, err := template.New(typ).Funcs(funcMap).Parse(src)
		if err != nil {
			log.Printf("[WARN] failed to parse built-in template %q: %v", typ, err)
			continue
		}
		registry[typ] = t
	}

	// env overrides: TEMPLATE_<TYPE>
	for typ := range defaultTemplates {
		envKey := "TEMPLATE_" + typ
		if src := os.Getenv(envKey); src != "" {
			t, err := template.New(typ).Funcs(funcMap).Parse(src)
			if err != nil {
				log.Printf("[WARN] env %s has invalid template, using built-in default: %v", envKey, err)
			} else {
				registry[typ] = t
				log.Printf("[INFO] loaded custom template for %s from env", typ)
			}
		}
	}

	// build fallback (TEMPLATE_DEFAULT or built-in genericTemplate)
	fallbackSrc := genericTemplate
	if src := os.Getenv("TEMPLATE_DEFAULT"); src != "" {
		t, err := template.New("default").Funcs(funcMap).Parse(src)
		if err != nil {
			log.Printf("[WARN] TEMPLATE_DEFAULT has invalid template, using built-in fallback: %v", err)
		} else {
			fallback = t
			log.Printf("[INFO] loaded custom default template from env")
		}
	}
	if fallback == nil {
		t, _ := template.New("default").Funcs(funcMap).Parse(fallbackSrc)
		fallback = t
	}
}

// Render converts an Alert to a Telegram HTML string.
func Render(a Alert) string {
	typ := strings.ToUpper(strings.ReplaceAll(a.Data.Type, " ", ""))
	tmpl, ok := registry[typ]
	if !ok {
		tmpl = fallback
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, a); err != nil {
		log.Printf("[WARN] template execute error for type %q: %v, using fallback", typ, err)
		buf.Reset()
		if err2 := fallback.Execute(&buf, a); err2 != nil {
			return fmt.Sprintf("%s <b>%s</b> · %s %s", levelEmoji(a.Level), htmlEsc(a.Level), htmlEsc(a.Data.Type), resolvedIcon(a.Resolved))
		}
	}

	msg := strings.TrimSpace(buf.String())
	if len(msg) > maxLen {
		// "…" is 3 bytes in UTF-8
		msg = msg[:maxLen-3] + "…"
	}
	return msg
}

// --- helper functions ---

func htmlEsc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func levelEmoji(level string) string {
	switch strings.ToUpper(level) {
	case "CRITICAL":
		return "🔴"
	case "ERROR":
		return "🚨"
	case "WARNING":
		return "⚠️"
	case "INFO":
		return "ℹ️"
	case "OK":
		return "✅"
	default:
		return "ℹ️"
	}
}

func resolvedIcon(resolved bool) string {
	if resolved {
		return "✅"
	}
	return "❌"
}

func jsonBlock(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return "<pre>" + htmlEsc(string(b)) + "</pre>"
}

func tsFormat(ts int64) string {
	if ts == 0 {
		return ""
	}
	t := time.Unix(ts/1000, 0).UTC()
	return t.Format("2006-01-02 15:04:05 UTC")
}

func safeGet(m map[string]any, key string) any {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	return v
}

func anyToStr(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
