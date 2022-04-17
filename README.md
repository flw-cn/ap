# ap -- `auto-pager`

ap is a shell utility that allows the output of other shell commands to automatically enter interactive page-flipping mode.

ap consists of two parts, a binary program written in Go that captures the output of commands and supports page-flipping, and a set of shell scripts that create a wrapper with the same name for a user-specified list of commands.

The usage of the commands after wrap is the same as before, and should not change the user's habits or cause any problems.

Read this in other languages: [English](README.en.md), [简体中文](README.md),[繁体中文](README.zh-TW.md),[Arabic](README.ar.md),[French](README.fr.md),[Hindi](README.hi.md)

ap 是一个 shell 工具，可以让其它 shell 命令的输出能够自动进入交互翻页模式。

ap 由两部分组成，一个 Go 语言编写的二进制程序，负责捕获命令的输出并支持翻页，
和一组 shell 脚本，负责为用户指定的命令清单创建与之同名的 wrapper。

经过 wrap 之后的命令用法与原来相同，不应当改变用户操作习惯，不会给用户造成困扰。





## 安装

```
go install github.com/flw-cn/ap@master
```

## 配置

* `bash`(#bash)
* `fish`(#`fish`)
* `zsh` (#`zsh`)

### bash

在你的 `~/.bashrc` 里加入下面内容：

```sh
eval "$(ap --bash)"
```

ap 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --bash)"
```

或者也可只在默认清单之上增加新的命令：

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --bash)"
```

另外，ap 也可以和 [grc](https://github.com/garabik/grc) 一起工作，在 macOS  下可以使用 Homebrew 安装 grc:

```sh
brew install grc
```

ap + grc 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --bash)"
```

或者也可只在默认清单之上增加新的命令：

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --bash)"
```



另外，环境变量 `$AUTO_PAGER_MIN_HEIGHT` 可以控制开始分页的最小行数：

```sh
AUTO_PAGER_MIN_HEIGHT=30        # 输出超过 30 行时才开始分页
eval "$(ap --bash)"
```

如果配置为负数，则代表终端窗口高度的百分比：

```sh
AUTO_PAGER_MIN_HEIGHT='-50'     # 输出超过终端窗口高度的 50% 时才开始分页
eval "$(ap --bash)"
```

如果不指定 `AUTO_PAGER_MIN_HEIGHT`，默认为 `-80`，即 `80%`。

### `fish`

在你的 `~/.config/fish/config.fish` 里加入下面内容：

```sh
ap --fish | source
```

ap 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```sh
set AUTO_PAGER_CMDS go cargo make
ap --fish | source
```

或者也可只在默认清单之上增加新的命令：

```sh
set AUTO_PAGER_CMDS_EXTRA ps last
ap --fish | source
```

另外，ap 也可以和 [grc](https://github.com/garabik/grc) 一起工作，在 macOS  下可以使用 Homebrew 安装 grc:

```sh
brew install grc
```

ap + grc 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```sh
set AUTO_PAGER_CMDS_WITH_GRC ps last dig diff
ap --fish | source
```

或者也可只在默认清单之上增加新的命令：

```sh
set AUTO_PAGER_CMDS_WITH_GRC_EXTRA ps last
ap --fish | source
```



另外，环境变量 `$AUTO_PAGER_MIN_HEIGHT` 可以控制开始分页的最小行数：

```sh
set AUTO_PAGER_MIN_HEIGHT 30        # 输出超过 30 行时才开始分页
ap --fish | source
```

如果配置为负数，则代表终端窗口高度的百分比：

```sh
set AUTO_PAGER_MIN_HEIGHT -50     # 输出超过终端窗口高度的 50% 时才开始分页
ap --fish | source
```

如果不指定 `AUTO_PAGER_MIN_HEIGHT`，默认为 `-80`，即 `80%`。

###  `zsh`

在你的 `~/.zshrc` 里加入下面内容：

```sh
eval "$(ap --zsh)"
```

ap 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --zsh)"
```

或者也可只在默认清单之上增加新的命令：

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --zsh)"
```

另外，ap 也可以和 [grc](https://github.com/garabik/grc) 一起工作，在 macOS  下可以使用 Homebrew 安装 grc:

```sh
brew install grc
```

ap + grc 默认 wrap 了一批命令。如果你不满意，可以通过环境变量重新定制：

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --zsh)"
```

或者也可只在默认清单之上增加新的命令：

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --zsh)"
```



另外，环境变量 `$AUTO_PAGER_MIN_HEIGHT` 可以控制开始分页的最小行数：

```sh
AUTO_PAGER_MIN_HEIGHT=30        # 输出超过 30 行时才开始分页
eval "$(ap --zsh)"
```

如果配置为负数，则代表终端窗口高度的百分比：

```sh
AUTO_PAGER_MIN_HEIGHT='-50'     # 输出超过终端窗口高度的 50% 时才开始分页
eval "$(ap --zsh)"
```

如果不指定 `AUTO_PAGER_MIN_HEIGHT`，默认为 `-80`，即 `80%`。

## 使用

经过 ap wrap 过的命令只需要像往常一样使用就可以了。
如果输出内容过多，会自动调用环境变量 `$PAGER` 所指定的分页器进行分页。

如果你的 `$PAGER` 变量比较特殊，无法与 ap 适配，你也可以通过环境变量 `$AP_PAGER`
为 ap 单独设置分页器。

如果 `$AP_PAGER` 和 `$PAGER` 都没有指定，那么将使用 `less -Fr`。

以下情况并不会启动分页器：
* 输出内容过少时，参见 `$AUTO_PAGER_MIN_HEIGHT`。
* 当 ap 检测到命令输出中包含 `ESC [?1049h` 序列时，此时命令被判定为全荧幕应用。
* 当命令尚未执行完成时。分页器只有命令执行完成后才会启动。
    - `ping` 和 `tcpdump` 等此类命令需要先按 `Ctrl-C` 终止命令后才会启动分页。
    - `python` 和 `gdb` 等此类命令需要先等待命令退出后才会启动分页。

## 常见问题

* 会影响命令的彩色输出吗？
    - 不会。
* 有的命令会检测终端，并为终端模式和非终端模式提供不同的输出，会改变它的输出吗？
    - 不会。
* 如果我怀疑 ap 影响了命令的输出，如何诊断？
    - 你可以用 `command foo` 来执行 `foo`，这样就不会调用 ap。
* 如果我习惯性地在 ap 过的命令后面加了 `| less`，会出问题吗？
    - 不会。
* ap 过的命令还可以重定向它的输出吗？
    - 可以。
* ap 过的命令的自动补全会被破坏吗？
    - 不会。
* ap 支持像 `python`、`gdb` 这样的交互式应用吗？
    - 支持。但是由于这些应用会输出一些控制字符，所以分页后看到的内容可能会有点乱。
* ap 支持像 `htop`、`vim` 这样的全荧幕应用吗？
    - 怎么说呢，反正不会出错，但我想不明白把 ap 和它们搭配在一起有什么实际意义。
