# GoTools

[![GoDoc](https://camo.githubusercontent.com/ba58c24fb3ac922ec74e491d3ff57ebac895cf2deada3bf1c9eebda4b25d93da/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f67616d6d617a65726f2f776f726b6572706f6f6c3f7374617475732e737667)](https://pkg.go.dev/github.com/520MianXiangDuiXiang520/GoTools)

<a title="GPL" target="_blank" href="https://github.com/520MianXiangDuiXiang520/GoTools/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-red.svg?style=flat-square"></a>
<a title="Last Commit" target="_blank" href="https://github.com/520MianXiangDuiXiang520/JuneGoBlog/commits/master"><img src="https://img.shields.io/github/last-commit/520MianXiangDuiXiang520/GoTools.svg?style=flat-square&color=FF9900"></a>
<a href="https://gitmoji.carloscuesta.me">
<img src="https://img.shields.io/badge/gitmoji-%20😜%20😍-FFDD67.svg?style=flat-square" alt="Gitmoji"></a>
<a href="https://goreportcard.com/badge/github.com/520MianXiangDuiXiang520/GoTools"> <img src="https://goreportcard.com/badge/github.com/520MianXiangDuiXiang520/GoTools" /></a>
<a href="https://codeclimate.com/github/520MianXiangDuiXiang520/GoTools/maintainability"><img src="https://api.codeclimate.com/v1/badges/ed575aea812a025dfcc9/maintainability" /></a>

包含一些平时 Go 开发过程中积累的小工具。
使用：

```go
go get github.com/520MianXiangDuiXiang520/GoTools
```

## CheckTools

这是一个通过结构体标签快速检查值是否合法的工具，可以为结构体字段添加 `check` 标签，并使用 `Check()` 函数检查，目前支持以下标签：

`int, int8, int16, int32, int64`:

| 标签示例                 | 作用                            |
| ------------------------ | ------------------------------- |
| `not null` 或 `not zero` | 非零判断                        |
| `size: [0, 10]`          | 判断范围在 0 到 10 之间，开区间 |
| `more: 10`               | 判断值大于10                    |
| `less: 10`               | 判断值小于 10                   |
| `equal: 10`              | 判断值等于 10                   |

`string`:

| 标签示例       | 作用                                |
| -------------- | ----------------------------------- |
| `not null`     | 不为空                              |
| `len: [2, 10]` | 字符串长度在 2 到 10 之间（闭区间） |
| `email`        | 判断是否是一个电子邮件              |

* **注意**：len 判断的是字符串底层字符数组的长度，对于中文或其他语言可能产生意外

`slice`:

| 标签示例       | 作用                              |
| -------------- | --------------------------------- |
| `len: [2, 10]` | 元素长度在 2 到 10 之间（开区间） |

`ptr`:

| 标签示例                | 作用             |
| ----------------------- | ---------------- |
| `not null` 或 `not nil` | 判断是否为空指针 |

`struct`:

如果一个 struct 包含另一个 struct, 则允许递归判断, 具体用法请参考 [godoc](https://pkg.go.dev/github.com/520MianXiangDuiXiang520/GoTools/check_tools)

## daoTools

与数据库相关的工具函数，目前包含：

* `conn`: 一个数据库连接工具
* `Transaction`: 一个数据库事务工具
* `redis`: 一个 redis 连接工具

具体用法请参考 [godoc](https://pkg.go.dev/github.com/520MianXiangDuiXiang520/GoTools/gin_tools/dao_tools) 文档

## emailTools

对 goemail 的简单封装，可以更加简单的实现群发，抄送，密送，附件等功能, 具体使用请参考 [godoc](https://pkg.go.dev/github.com/520MianXiangDuiXiang520/GoTools/email_tools)

