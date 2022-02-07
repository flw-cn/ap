ap -- auto-pager

ap 是一个 shell 工具，可以让其它 shell 命令的输出能够自动进入交互翻页模式。

ap 由两部分组成，一个 Go 语言编写的二进制程序，负责捕获命令的输出并支持翻页，
和一组 shell 脚本，负责为用户指定的命令清单创建与之同名的 wrapper。

经过 wrap 之后的命令用法与原来相同，不应当改变用户操作习惯，不会给用户造成困扰。

## 安装

```
go install github.com/flw-cn/ap
```

## 配置

在你的 `~/.zshrc` 里加入下面内容：

```
eval "$(ap --zsh)"
```

ap 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --zsh)"
```

或者也可只在默认清单之上增加新的命令：

```
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --zsh)"
```

另外，ap 也可以和 [grc](https://github.com/garabik/grc) 一起工作：

```
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --zsh)"
```

## 使用

经过 ap wrap 过的命令只需要像往常一样使用就可以了，
如果输出内容过多，会自动调用环境变量 `$PAGER` 所指定的分页器进行分页。

如果你的 `$PAGER` 变量比较特殊，无法与 ap 适配，你也可以通过以下两种方式定制：
* 通过环境变量 `$AP_PAGER`（优先级比 `$PAGER` 高）
* 通过命令行参数 `--pager <pager>`（优先级比 `$AP_PAGER` 和 `$PAGER` 要高）

如果 `--pager`、`$AP_PAGER` 和 `$PAGER` 都没有指定，那么将使用 `less -FR`。

## 常见问题

* 会影响命令的彩色输出吗？
    - 不会。
* 有的命令会检测终端，并为终端模式和非终端模式提供不同的输出，会改变它的输出吗？
    - 不会。
* 如果我怀疑 ap 影响了命令的输出，如何诊断？
    - 在 zsh 下你可以用 `command foo ...` 来执行 `foo`，这样就不会调用 ap。
* 如果我习惯性地在 ap 过的命令后面加了 `| less`，会出问题吗？
    - 不会。
* ap 过的命令还可以重定向它的输出吗？
    - 可以。
* ap 过的命令的自动补全会被破坏吗？
    - 不会。
* ap 支持像 `htop`、`vim` 这样的全荧幕应用吗？
    - 不支持。如果你希望支持它，请联系我并说明你想看到的效果。
