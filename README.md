[![zerobot-plugin-mcqq](https://socialify.git.ci/17TheWord/zerobot-plugin-mcqq/image?description=1&forks=1&issues=1&language=1&logo=https%3A%2F%2Fraw.githubusercontent.com%2F17TheWord%2Fnonebot-adapter-minecraft%2Fmain%2Fassets%2Flogo.png&name=1&owner=1&pulls=1&stargazers=1&theme=Auto)](https://github.com/17TheWord/zerobot-plugin-mcqq)

# ZeroBot-Plugin-MCQQ

`mcqq` 的 `ZeroBot` 插件实现

- 支持 QQ 群
- 支持鹊桥 Protocol V2

## 使用

- 克隆本项目

- 复制 `config.example.yml` 为 `config.yml`，并按注释修改配置

- 启动 `go run . -config config.yml`

- 编译后启动

    ```bash
    go build -o mcqq
    ./mcqq -config config.yml
    ```

    Windows:

    ```powershell
    go build -o mcqq.exe
    .\mcqq.exe -config config.yml
    ```

## 对接 Minecraft 服务器

配套 **插件/模组** 请前往 [`鹊桥`](https://github.com/17TheWord/QueQiao) 仓库查看详情

当前版本按 [`鹊桥 Protocol V2`](https://queqiao-docs.pages.dev/) 对接：

- QQ 群消息通过 `broadcast` API 广播到 Minecraft 服务器
- Minecraft 事件通过 WebSocket JSON 推送到本插件
- WebSocket Header 默认使用 `x-client-origin: zerobot`
- WebSocket Client 支持断线自动重连，可在 `config.yml` 中配置 `reconnect_interval` 和 `reconnect_max_times`

## 配置说明

项目只支持 YAML 配置文件。默认读取 `config.yml`，也可以通过 `-config` 指定其他路径。

`zero.drivers` 用于配置 ZeroBot 和 OneBot 实现之间的连接：

- `websocket_server`：本程序作为 WebSocket 服务端，等待 OneBot 实现连接。
- `websocket_client`：本程序主动连接 OneBot 实现。

`mcqq.websocket_server` 和 `mcqq.websocket_client` 用于配置本程序和 Minecraft 鹊桥端之间的连接：

- `mcqq.websocket_server.enable: true`：本程序作为 WebSocket 服务端，等待 Minecraft 端连接。
- `mcqq.websocket_client`：本程序主动连接 Minecraft 端。
- 两种 Minecraft 连接方式至少需要启用一种。

`mcqq.server_map` 是服务器到 QQ 群的路由表：

- key 需要和鹊桥 Minecraft 端 `server_name` 一致。
- `group_list[].bot_id` 是发送群消息时使用的 Bot ID。
- `group_list[].group_id` 是 QQ 群号。
- 同一个 `group_id` 不能同时映射到多个服务器，避免 QQ 群消息被意外广播到多个 Minecraft 服务器。

`mcqq.websocket_client[].reconnect_max_times` 为 `0` 时表示无限重连。

## 常见问题

### 提示 `配置校验失败`

请优先对照 `config.example.yml` 检查字段名和缩进。YAML 对缩进敏感，建议使用两个空格缩进。

### 提示 `No active websocket connection`

表示 QQ 群消息找到了目标服务器配置，但该服务器没有可用的 Minecraft WebSocket 连接。请检查：

- 鹊桥端 `server_name` 是否和 `mcqq.server_map`、`mcqq.websocket_client[].server_name` 一致。
- `mcqq.websocket_client[].url` 是否正确。
- 如果使用 `mcqq.websocket_server`，Minecraft 端是否成功连接到本程序。
- `access_token` 是否两端一致。

### 提示 `Failed to get bot with id`

表示 `group_list[].bot_id` 没有匹配到当前 ZeroBot 已连接的 Bot。请检查 Bot ID 是否填写为实际机器人 QQ 号或对应适配器要求的 Bot ID。

### QQ 消息没有进入 Minecraft

请检查：

- QQ 群号是否写在 `mcqq.server_map.*.group_list[].group_id`。
- 当前 Bot 是否已经加入该群。
- Minecraft WebSocket 是否已连接。
- 鹊桥端是否使用 Protocol V2。

## 功能

- 推送消息列表

    - 服务器 -> QQ
        - [x] 加入 / 离开 服务器消息
        - [x] 玩家聊天信息
        - [x] 玩家命令信息
        - [x] 玩家死亡信息（死亡信息为英文，原版端不适用，用**正则**匹配死亡信息是大工程！）
        - [x] 玩家成就信息
    - QQ -> 服务器
        - [x] 鹊桥 V2 API 命令调用
        - [x] 群员聊天文本
        - [x] 图片、视频等内容转换为可点击在浏览器打开的 `[图片]`、`[视频]`
        - [x] (可选功能)借助 [`@kitUIN/ChatImage`](https://github.com/kitUIN/ChatImage) 直接在游戏内显示图片

## QQ 命令

以下命令需要 ZeroBot `super_users` 权限，并且只能在已配置的 QQ 群中使用。命令会通过鹊桥 V2 API 发送到 Minecraft 端，并根据响应包的 `echo` 将结果返回到原 QQ 群。

```text
/mc status <server>
/mc rcon <server> <command>
/mc title <server> <title> [| <subtitle>]
/mc actionbar <server> <message>
/mc private <server> <uuid|-> <nickname|-> <message>
```

示例：

```text
/mc status Server
/mc rcon Server list
/mc title Server 服务器公告 | 欢迎回来
/mc actionbar Server 当前正在维护
/mc private Server - Steve 你好
```

说明：

- `server` 需要和 `mcqq.server_map` 中的服务器名一致。
- `/mc private` 中 `uuid` 或 `nickname` 至少填写一个，不使用的字段写 `-`。
- API 响应默认等待 60 秒，超时会向原 QQ 群发送超时提示。

## 特别感谢

- [`ZeroBot`](https://github.com/wdvxdr1123/ZeroBot)：插件使用的开发框架。
- [`@kitUIN/ChatImage`](https://github.com/kitUIN/ChatImage)：用于在游戏内显示图片的插件。

## 其他项目

- [`@17TheWord/nonebot-plugin-mcqq`](https://github.com/17TheWord/nonebot-plugin-mcqq)；适用于 `NoneBot` 的互通消息插件
- [`@17TheWord/nonebot-adapter-minecraft`](https://github.com/17TheWord/nonebot-adapter-minecraft)：适用于 `NoneBot` 的
  `Minecraft Server` 适配器
- [`@CikeyQi/mc-plugin`](https://github.com/CikeyQi/mc-plugin)：适用于 `Yunzai` 的互通消息插件

## 贡献与支持

觉得好用可以给这个项目点个 `Star` 或者去 [爱发电](https://afdian.com/a/17TheWord) 投喂我。

有意见或者建议也欢迎提交 [Issues](https://github.com/17TheWord/zerobot-plugin-mcqq/issues)
和 [Pull requests](https://github.com/17TheWord/zerobot-plugin-mcqq/pulls)。

## 许可证

本项目使用 [`MIT`](./LICENSE) 作为开源许可证。
