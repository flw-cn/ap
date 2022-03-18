# ap - 自動尋呼機

ap 是一個 shell 實用程序，它允許其他 shell 命令的輸出自動進入交互式翻頁模式。

ap 由兩部分組成，一個用 Go 編寫的二進製程序，它捕獲命令的輸出並支持頁面翻轉，以及一組 shell 腳本，它為用戶指定的命令列表創建一個同名的包裝器。

wrap 之後的命令用法和之前一樣，不應該改變用戶的習慣或造成任何問題。

\*用其他語言閱讀：[英語](README_en.md),[簡體中文](README.md)

ap 是一個 shell 工具，可以讓其它 shell 命令的輸出能夠自動進入交互翻頁模式。

ap 由兩部分組成，一個 Go 語言編寫的二進製程序，負責捕獲命令的輸出並支持翻頁，
和一組 shell 腳本，負責為用戶指定的命令清單創建與之同名的 wrapper。

經過 wrap 之後的命令用法與原來相同，不應當改變用戶操作習慣，不會給用戶造成困擾。

## 安裝

    go install github.com/flw-cn/ap@master

## 配置

-   [重擊](#bash)
-   [魚](#fish)
-   [zsh](#zsh)

### 重擊

在你的`~/.bashrc`裡加入下面內容：

```sh
eval "$(ap --bash)"
```

ap 默認 wrap 了一批命令。如果你不滿意，可以通過環境變量重新定制：

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --bash)"
```

或者也可只在默認清單之上增加新的命令：

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --bash)"
```

另外，ap 也可以和[grc](https://github.com/garabik/grc)一起工作，在 macOS  下可以使用 Homebrew 安裝 grc:

```sh
brew install grc
```

ap + grc 默認 wrap 了一批命令。如果你不滿意，可以通過環境變量重新定制：

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --bash)"
```

或者也可只在默认清单之上增加新的命令：

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --bash)"
```

另外，環境變量`$AUTO_PAGER_MIN_HEIGHT` 可以控制开始分页的最小行数：

```sh
AUTO_PAGER_MIN_HEIGHT=30        # 输出超过 30 行时才开始分页
eval "$(ap --bash)"
```

如果配置為負數，則代表終端窗口高度的百分比：

```sh
AUTO_PAGER_MIN_HEIGHT='-50'     # 输出超过终端窗口高度的 50% 时才开始分页
eval "$(ap --bash)"
```

如果不指定`AUTO_PAGER_MIN_HEIGHT`，默認為`-80`，即`80%`。

### 魚

在你的`~/.config/fish/config.fish`裡加入下面內容：

```sh
ap --fish | source
```

ap 默認 wrap 了一批命令。如果你不滿意，可以通過環境變量重新定制：

```sh
set AUTO_PAGER_CMDS go cargo make
ap --fish | source
```

或者也可只在默認清單之上增加新的命令：

```sh
set AUTO_PAGER_CMDS_EXTRA ps last
ap --fish | source
```

另外，ap 也可以和[grc](https://github.com/garabik/grc)一起工作，在 macOS  下可以使用 Homebrew 安裝 grc:

```sh
brew install grc
```

ap + grc 默認 wrap 了一批命令。如果你不滿意，可以通過環境變量重新定制：

```sh
set AUTO_PAGER_CMDS_WITH_GRC ps last dig diff
ap --fish | source
```

或者也可只在默認清單之上增加新的命令：

```sh
set AUTO_PAGER_CMDS_WITH_GRC_EXTRA ps last
ap --fish | source
```

另外，環境變量`$AUTO_PAGER_MIN_HEIGHT`可以控制開始分頁的最小行數：

```sh
set AUTO_PAGER_MIN_HEIGHT 30        # 输出超过 30 行时才开始分页
ap --fish | source
```

如果配置為負數，則代表終端窗口高度的百分比：

```sh
set AUTO_PAGER_MIN_HEIGHT -50     # 输出超过终端窗口高度的 50% 时才开始分页
ap --fish | source
```

如果不指定`AUTO_PAGER_MIN_HEIGHT`，默認為`-80`，即`80%`。

### zsh

在你的`~/.zshrc`裡加入下面內容：

```sh
eval "$(ap --zsh)"
```

ap 默認 wrap 了一批命令。如果你不滿意，可以通過環境變量重新定制：

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --zsh)"
```

或者也可只在默認清單之上增加新的命令：

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --zsh)"
```

另外，ap 也可以和[grc](https://github.com/garabik/grc)一起工作，在 macOS  下可以使用 Homebrew 安裝 grc:

```sh
brew install grc
```

ap + grc 默認 wrap 了一批命令。如果你不滿意，可以通過環境變量重新定制：

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --zsh)"
```

或者也可只在默認清單之上增加新的命令：

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --zsh)"
```

另外，環境變量`$AUTO_PAGER_MIN_HEIGHT`可以控制開始分頁的最小行數：

```sh
AUTO_PAGER_MIN_HEIGHT=30        # 输出超过 30 行时才开始分页
eval "$(ap --zsh)"
```

如果配置為負數，則代表終端窗口高度的百分比：

```sh
AUTO_PAGER_MIN_HEIGHT='-50'     # 输出超过终端窗口高度的 50% 时才开始分页
eval "$(ap --zsh)"
```

如果不指定`AUTO_PAGER_MIN_HEIGHT`，默認為`-80`，即`80%`。

## 使用

經過 ap wrap 過的命令只需要像往常一樣使用就可以了。
如果輸出內容過多，會自動調用環境變量`$PAGER`所指定的分頁器進行分頁。

如果你的`$PAGER`變量比較特殊，無法與 ap 適配，你也可以通過環境變量`$AP_PAGER`為 ap 單獨設置分頁器。

如果`$AP_PAGER`和`$PAGER`都沒有指定，那麼將使用`less -Fr`。

以下情況並不會啟動分頁器：

-   輸出內容過少時，參見`$AUTO_PAGER_MIN_HEIGHT`。
-   當 ap 檢測到命令輸出中包含`ESC [?1049h`序列時，此時命令被判定為全熒幕應用。
-   當命令尚未執行完成時。分頁器只有命令執行完成後才會啟動。
    -   `ping`和`tcpdump`等此類命令需要先按`Ctrl-C`終止命令後才會啟動分頁。
    -   `python`和`gdb`等此類命令需要先等待命令退出後才會啟動分頁。

## 常見問題

-   會影響命令的彩色輸出嗎？
    -   不會。
-   有的命令會檢測終端，並為終端模式和非終端模式提供不同的輸出，會改變它的輸出嗎？
    -   不會。
-   如果我懷疑 ap 影響了命令的輸出，如何診斷？
    -   你可以用`command foo`來執行`foo`，這樣就不會調用 ap。
-   如果我習慣性地在 ap 過的命令後面加了`| less`，會出問題嗎？
    -   不會。
-   ap 過的命令還可以重定向它的輸出嗎？
    -   可以。
-   ap 過的命令的自動補全會被破壞嗎？
    -   不會。
-   ap 支持像`python`、`gdb`這樣的交互式應用嗎？
    -   支持。但是由於這些應用會輸出一些控製字符，所以分頁後看到的內容可能會有點亂。
-   ap 支持像`htop`、`vim`這樣的全熒幕應用嗎？
    -   怎麼說呢，反正不會出錯，但我想不明白把 ap 和它們搭配在一起有什麼實際意義。
