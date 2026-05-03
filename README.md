# Komodo Telegram Alerter

[![Build](https://github.com/haswelldev/komodo-telegram-alerter/actions/workflows/docker.yml/badge.svg)](https://github.com/haswelldev/komodo-telegram-alerter/actions/workflows/docker.yml)

A lightweight Telegram alerter for [Komodo](https://komo.do), written in Go.
Forked from [SashaBusinaro/komodo-telegram-alerter](https://github.com/SashaBusinaro/komodo-telegram-alerter).

- ~10 MB Docker image (`scratch`-based), ~5–10 MB idle RSS
- Beautiful, alert-type-aware Telegram HTML messages
- Fully customizable via environment variables — no config required by default

## Example Notifications

**🔴 CRITICAL · ServerUnreachable**
```
🔴 CRITICAL · ServerUnreachable ❌
Server: M14
Error: Timed out waiting for Ping
       deadline has elapsed
```

**⚠️ WARNING · StackStateChange**
```
⚠️ WARNING · StackStateChange ✅
Stack: redis
Server: Hetzner
State: stopped → running
```

**✅ OK · ServerUnreachable**
```
✅ OK · ServerUnreachable ✅
Server: M14
Error: Failed to connect to websocket
       IO error: No route to host (os error 113)
```

## Quick Start (Docker Compose)

Minimal setup — works with no environment variables:

```yaml
services:
  komodo-telegram-alerter:
    image: ghcr.io/haswelldev/komodo-telegram-alerter:latest
    container_name: komodo-telegram-alerter
    restart: unless-stopped
    ports:
      - '3000:3000'
```

<details>
<summary>With template overrides</summary>

```yaml
services:
  komodo-telegram-alerter:
    image: ghcr.io/haswelldev/komodo-telegram-alerter:latest
    container_name: komodo-telegram-alerter
    restart: unless-stopped
    ports:
      - '3000:3000'
    environment:
      PORT: 3000

      # Override the ServerUnreachable template with a more concise format
      TEMPLATE_SERVERUNREACHABLE: >-
        {{emoji .Level}} <b>{{esc .Level}}</b> · ServerUnreachable {{resolvedIcon .Resolved}}
        <b>Server:</b> {{esc (str (get .Data.Data "name"))}}
        {{- with get .Data.Data "err"}}{{$e := .}}{{with get $e "error"}}
        <b>Error:</b> <code>{{esc (str .)}}</code>{{end}}{{end}}

      # Override StackStateChange to include a custom header line
      TEMPLATE_STACKSTATECHANGE: >-
        {{emoji .Level}} <b>Stack state changed</b> {{resolvedIcon .Resolved}}
        <b>Stack:</b> {{esc (str (get .Data.Data "name"))}}
        <b>Server:</b> {{esc (str (get .Data.Data "server_name"))}}
        <b>Transition:</b> <code>{{esc (str (get .Data.Data "from"))}}</code> → <code>{{esc (str (get .Data.Data "to"))}}</code>

      # Generic fallback for any alert type without a specific template
      TEMPLATE_DEFAULT: >-
        {{emoji .Level}} <b>{{esc .Level}}</b> · {{esc .Data.Type}} {{resolvedIcon .Resolved}}
        {{- with (get .Data.Data "name") | str}}{{if .}}
        <b>Name:</b> {{esc .}}{{end}}{{end}}
        {{json .Data.Data}}
```

</details>

<details>
<summary>Running alongside Komodo (same compose stack)</summary>

When the alerter is in the same Docker Compose stack as Komodo, no port forwarding is needed — use the container name as the hostname.

```yaml
services:
  komodo-core:
    image: ghcr.io/moghtech/komodo-core:latest
    container_name: komodo-core
    restart: unless-stopped
    # ... your existing Komodo config

  komodo-telegram-alerter:
    image: ghcr.io/haswelldev/komodo-telegram-alerter:latest
    container_name: komodo-telegram-alerter
    restart: unless-stopped
    # No ports: block needed — Komodo reaches it over the internal network
```

In Komodo, set your Custom Alerter URL to:

`http://komodo-telegram-alerter:3000/alert?token=[[TELEGRAM_TOKEN]]&chat_id=[[TELEGRAM_CHAT_ID]]`

</details>

### Configure Komodo

In Komodo, add a Custom Alerter with the following URL:

`http://<komodo-telegram-alerter-ip>:3000/alert?token=<TELEGRAM_TOKEN>&chat_id=<TELEGRAM_CHAT_ID>`

Or, leverage Komodo's interpolation:

`http://<komodo-telegram-alerter-ip>:3000/alert?token=[[TELEGRAM_TOKEN]]&chat_id=[[TELEGRAM_CHAT_ID]]`

**Recommended**: Use [Komodo Secrets & Variables](https://komo.do/docs/variables) to store your Telegram credentials.

---

## Configuration

All configuration is via environment variables. None are required.

| Variable | Default | Description |
|---|---|---|
| `PORT` | `3000` | HTTP port to listen on |
| `TEMPLATE_DEFAULT` | built-in | Fallback template used for unknown alert types |
| `TEMPLATE_<TYPE>` | built-in | Per-type template override (see below) |

### Per-type template overrides

Set `TEMPLATE_<ALERTTYPE>` (uppercased, no spaces) to override a specific alert type's message format.

**Supported types** (with built-in templates):

| Env var | Komodo alert type |
|---|---|
| `TEMPLATE_SERVERUNREACHABLE` | ServerUnreachable |
| `TEMPLATE_SERVERCPU` | ServerCpu |
| `TEMPLATE_SERVERMEM` | ServerMem |
| `TEMPLATE_SERVERDISK` | ServerDisk |
| `TEMPLATE_SERVERTEMP` | ServerTemp |
| `TEMPLATE_CONTAINERSTATECHANGE` | ContainerStateChange |
| `TEMPLATE_STACKSTATECHANGE` | StackStateChange |
| `TEMPLATE_STACKAUTOUPDATED` | StackAutoUpdated |
| `TEMPLATE_DEPLOYMENTSTATECHANGE` | DeploymentStateChange |
| `TEMPLATE_BUILDFAILED` | BuildFailed |
| `TEMPLATE_RESOURCESYNCPENDINGUPDATES` | ResourceSyncPendingUpdates |
| `TEMPLATE_AWSBUILDERTERMATIONFAILED` | AwsBuilderTerminationFailed |

Unknown types fall back to `TEMPLATE_DEFAULT` (or the built-in generic template).

### Template syntax

Templates use [Go `text/template`](https://pkg.go.dev/text/template) syntax. The root object (`.`) is the full Komodo alert:

| Field | Type | Description |
|---|---|---|
| `.Level` | string | `CRITICAL`, `ERROR`, `WARNING`, `INFO`, `OK` |
| `.Resolved` | bool | Whether the alert is resolved |
| `.Target.Type` | string | e.g. `Server`, `Stack`, `Container` |
| `.Target.ID` | string | Target resource ID |
| `.Data.Type` | string | Alert type name |
| `.Data.Data` | map | Type-specific payload fields |
| `.Ts` | int64 | Unix timestamp in milliseconds |

**Helper functions:**

| Function | Description |
|---|---|
| `esc .Value` | HTML-escape `&`, `<`, `>` (always use for user data) |
| `bold .Value` | Wrap in `<b>…</b>` with escaping |
| `italic .Value` | Wrap in `<i>…</i>` with escaping |
| `code .Value` | Wrap in `<code>…</code>` with escaping |
| `emoji .Level` | Level → emoji (🔴 ⚠️ 🚨 ℹ️ ✅) |
| `resolvedIcon .Resolved` | `true`→✅ / `false`→❌ |
| `json .Value` | Pretty-print as `<pre>…</pre>` |
| `ts .Ts` | Format unix-ms timestamp as `2006-01-02 15:04:05 UTC` |
| `get .Data.Data "key"` | Safe map lookup (returns `""` if missing) |
| `str .Value` | Convert any value to string |

**Example custom template:**

```
TEMPLATE_SERVERUNREACHABLE='{{emoji .Level}} <b>{{esc .Level}}</b> — {{esc (str (get .Data.Data "name"))}} is down {{resolvedIcon .Resolved}}'
```

If a template fails to parse at startup, a warning is logged and the built-in default is used — the service always starts.

---

## Getting Your Telegram Credentials

### Bot Token
1. Message [@BotFather](https://t.me/botfather) on Telegram
2. Create a new bot with `/newbot`
3. Follow the instructions to get your bot token

### Chat ID
1. Add your bot to the desired chat/channel
2. Send a message to the chat
3. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Look for the `chat.id` field in the response
