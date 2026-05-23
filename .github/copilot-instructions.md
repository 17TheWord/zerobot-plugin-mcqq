# Project Instructions

## Project Overview

This repository is `github.com/17TheWord/zerobot-plugin-mcqq`, a Go 1.20 ZeroBot plugin that bridges QQ group messages and Minecraft server messages.

The plugin connects ZeroBot to Minecraft-side QueQiao-compatible plugins/mods over WebSocket. It uses QueQiao Protocol V2 only. It supports forwarding QQ group messages to Minecraft as Minecraft chat components through the V2 `broadcast` API, and forwarding Minecraft player events back to configured QQ groups.

Keep changes small and aligned with the existing simple package layout. Prefer direct, readable code over introducing abstractions unless they are clearly reused.

## Key Technologies

- Language: Go.
- Module path: `github.com/17TheWord/zerobot-plugin-mcqq`.
- Bot framework: `github.com/wdvxdr1123/ZeroBot`.
- WebSocket implementation: `github.com/RomiChan/websocket`.
- Logging: `github.com/sirupsen/logrus` with `logrus-easy-formatter`.
- Minecraft bridge protocol: QueQiao Protocol V2 JSON events and API payloads.

## Repository Layout

- `main.go`: executable entrypoint. It reads YAML config, initializes logging, builds ZeroBot drivers, initializes `mcqq`, and starts ZeroBot.
- `app_config.go`: application-level YAML config loading, defaults, validation, and ZeroBot driver construction.
- `config.example.yml`: documented YAML example. Users should copy it to `config.yml` and edit values before running the binary.
- `mcqq/config.go`: configuration structs for WebSocket server/client, Minecraft servers, QQ groups, command priority, ChatImage, and access token.
- `mcqq/mcqq.go`: plugin initialization, group filter setup, command priority bounds, WebSocket startup, and ZeroBot message handler registration.
- `mcqq/wsserver.go`: WebSocket server mode. Accepts Minecraft connections on `/minecraft`, `/minecraft/`, `/minecraft/ws`, and `/minecraft/ws/`.
- `mcqq/wsclient.go`: WebSocket client mode. Connects to configured Minecraft-side WebSocket endpoints.
- `mcqq/handle_message.go`: core message conversion and forwarding logic in both directions.
- `mcqq/proto.go`: Minecraft chat component structs, V2 API constants, and outbound WebSocket API envelope.
- `mcqq/event.go`: Minecraft V2 event payload structs for chat, command, join/quit, death, achievements, `Translate`, and API responses.
- `mcqq/rule.go`: ZeroBot group-message rule based on configured QQ groups.
- `mcqq/util.go`: target server lookup and WebSocket cleanup helpers.
- `mcqq/cache.go`: in-memory QQ group and group member display-name caches.
- `.github/workflows/release.yml`: release workflow. It runs `go mod tidy`, `go test -v`, `go vet`, `golangci-lint`, and creates a GitHub release on push to `main`.

## Runtime Model

`main.go` reads `config.yml` by default, assigns `mcqq.PluginConfig`, calls `mcqq.InitPlugin()`, then starts ZeroBot through `zero.RunAndBlock`.

Configuration is YAML-only. Do not add JSON, dotenv, or multi-format config support unless explicitly requested. Users should not edit `main.go` for deployment configuration.

`PluginConfig.ServerMap` is the main routing table. Each server name maps to a `ServerConfig` containing QQ group IDs and bot IDs. QQ messages from configured groups are forwarded to every configured Minecraft server that includes that group. Minecraft events are sent back to every QQ group configured for that server.

The plugin can use either or both WebSocket modes:

- Server mode: `WebsocketServer.Enable` starts an HTTP WebSocket server. Minecraft connects to this plugin and must provide `x-self-name`. If `AccessToken` is configured, `Authorization: Bearer <token>` is required.
- Client mode: each `WebsocketClientConfig` starts a goroutine that connects to a Minecraft-side WebSocket URL. The plugin sends `x-client-origin: zerobot`, `x-self-name: <ServerName>`, and `Authorization: Bearer <AccessToken>` headers.

Active Minecraft connections are managed by a mutex-protected connection store keyed by server name. Do not access connection maps directly; use the store/helper functions.

## Message Flow

QQ to Minecraft:

- `zero.OnMessage(GroupRule)` registers `handleQQMessage`.
- `GroupRule` accepts only group messages whose group ID appears in `PluginConfig.ServerMap`.
- `processQQMessage2MinecraftProtocol` creates a Minecraft chat component list containing QQ group name, sender name, and converted message elements.
- `processQQMessageList` maps QQ segments to Minecraft components:
  - `text`: white plain text.
  - `reply`: gray reply metadata and hover text for the original message.
  - `face`: `[表情]` with face ID hover text.
  - `file`: `[文件]` with file name hover text.
  - `image`: either ChatImage `[[CICode,url=...,name=图片]]` when `ChatImage` is enabled, or clickable `[图片]` with URL hover/click metadata.
  - `record`: `[语音]`.
  - `video`: clickable `[视频]`.
  - `at`: resolves display name from cache or ZeroBot group member info.
  - unknown segment types become `[type]`.
- `handleQQMessage` wraps the component list in `WebsocketData{API: "broadcast", Data: {"message": ...}, Echo: <millisecond timestamp>}` and writes it to each target Minecraft WebSocket.

Minecraft to QQ:

- WebSocket readers call `handleMinecraftMessage` for every incoming message.
- `post_type == "response"` is logged and ignored.
- Supported `sub_type` values are `player_chat`, `player_command`, `player_join`, `player_quit`, `player_death`, and `player_achievement`.
- The QQ message format is plain text prefixed with `[serverName]`.
- `sendMcMsg2QQGroup` looks up `PluginConfig.ServerMap[serverName]`, obtains the configured ZeroBot instance with `zero.GetBot(BotId)`, and sends group messages.

## Configuration Notes

- `CommandPriority` is clamped to the inclusive range `2..98` during initialization.
- `AccessToken` is optional. Empty token disables server-side authorization checks but client mode still sends `Authorization: Bearer `.
- `ChatImage` changes image forwarding behavior. When enabled, image message text is replaced with ChatImage CICode for in-game image rendering.
- `RconMsg` exists in `ServerConfig` but is not currently used by the implementation.
- QueQiao `0.4.1+` uses `Translate` for death and achievement text. Prefer `translate.text`, then legacy `text`, then key fallback.
- WebSocket client reconnect settings are configured per client with `reconnect_interval` seconds and `reconnect_max_times`; `0` max means reconnect indefinitely.
- `main.go` is the application entrypoint. Keep deployment configuration in YAML rather than hard-coding values in Go.

## Protocol Expectations

Outbound messages to Minecraft use:

```json
{
  "api": "broadcast",
  "data": {
    "message": [
      { "text": "...", "color": "..." }
    ]
  },
  "echo": "..."
}
```

Inbound Minecraft events are expected to include at least `post_type`, `server_name`, and `sub_type`. Event-specific fields are defined in `mcqq/event.go`.

API response packets use `post_type: "response"` and are logged by the adapter. QQ command wrappers for Rcon/status/title/actionbar are not implemented yet.

Minecraft chat components are modeled in `mcqq/proto.go`. Preserve JSON tags and omit-empty behavior when modifying protocol structs.

## Development Guidelines

- Follow existing Go style and run `gofmt` on modified Go files.
- Keep package code in `mcqq` unless adding a new executable entrypoint.
- Keep runtime configuration in YAML. The default path is `config.yml`, and the repository should track `config.example.yml` rather than real local config files.
- Preserve existing public struct field names unless a breaking change is intentional.
- Be careful with global maps (`PluginConfig`, caches). Minecraft WebSocket connections are protected by a connection store; keep connection access behind that abstraction.
- Prefer exact JSON field names already used by QueQiao-compatible plugins/mods.
- When adding support for new QQ message segments or Minecraft `sub_type` events, update conversion logic and relevant protocol/event structs together.
- Do not reintroduce V1 protocol behavior unless explicitly requested. This adapter targets QueQiao Protocol V2 only.
- Do not remove README-documented behavior such as image/video clickable links or optional ChatImage support.
- Keep Chinese user-facing/logging text consistent with the existing codebase unless there is a reason to change it.

## Validation

Use these commands before finalizing code changes when feasible:

```bash
go test ./...
go vet ./...
go mod tidy
```

The GitHub release workflow currently runs `go test -v` and `go vet` without `./...`, but local validation should prefer package-wide commands.

## Known Limitations

- QQ to Minecraft command forwarding is listed as not implemented in the README.
- QQ wrappers for Minecraft `Rcon`, `Title`, `ActionBar`, private message, and status query APIs are not implemented.
- Death and achievement text are read from QueQiao V2 `Translate` or legacy text fields, but no local translation files are loaded by this project.
- There are no dedicated tests in the repository at the time this instruction file was written.
