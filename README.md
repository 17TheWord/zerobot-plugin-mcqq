[![zerobot-plugin-mcqq](https://socialify.git.ci/17TheWord/zerobot-plugin-mcqq/image?description=1&forks=1&issues=1&language=1&logo=https%3A%2F%2Fraw.githubusercontent.com%2F17TheWord%2Fnonebot-adapter-minecraft%2Fmain%2Fassets%2Flogo.png&name=1&owner=1&pulls=1&stargazers=1&theme=Auto)](https://github.com/17TheWord/zerobot-plugin-mcqq)

# ZeroBot-Plugin-MCQQ

`mcqq` 的 `ZeroBot` 插件实现

- 支持 QQ 群

## 使用

- 克隆本项目

- 修改 `main.go` 中的配置

- 启动 `go run main.go`

## 对接 Minecraft 服务器

配套 **插件/模组** 请前往 [`鹊桥`](https://github.com/17TheWord/QueQiao) 仓库查看详情

## 功能

- 推送消息列表

    - 服务器 -> QQ
        - [x] 加入 / 离开 服务器消息
        - [x] 玩家聊天信息
        - [x] 玩家死亡信息（死亡信息为英文，原版端不适用，用**正则**匹配死亡信息是大工程！）
    - QQ -> 服务器
        - [ ] 指令
        - [x] 群员聊天文本
        - [x] 图片、视频等内容转换为可点击在浏览器打开的 `[图片]`、`[视频]`
        - [x] (可选功能)借助 [`@kitUIN/ChatImage`](https://github.com/kitUIN/ChatImage) 直接在游戏内显示图片

> 个人能力问题，暂未实现 `Command`、`Title`、`ActionBar` 等内容，欢迎 `PR`

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
