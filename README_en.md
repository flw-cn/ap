# ap -- auto-pager

ap is a shell utility that enables the output of other shell commands to automatically enter interactive page-flipping mode.

ap consists of two parts, a binary program written in Go that captures the output of commands and supports pagination, and a set of shell scripts that create a wrapper with the same name for a user-specified list of commands.
and a set of shell scripts that create a wrapper with the same name for a user-specified list of commands.

The usage of the commands after wrap is the same as before, and should not change the user's habits or cause any problems.

## Installation

```
go install github.com/flw-cn/ap@master
```

## Configure

* [bash](#bash)
* [fish](#fish)
* [zsh](#zsh)

### bash

Add the following to your `~/.bashrc`.

```sh
eval "$(ap --bash)"
```

ap wrap a bunch of commands by default. If you are not satisfied, you can re-customize it via environment variables.

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --bash)"
```

Or you can just add new commands on top of the default manifest: ```sh

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --bash)"
```

Also, ap works with [grc](https://github.com/garabik/grc), which can be installed on macOS using Homebrew:

```sh
brew install grc
```

ap + grc wrap a bunch of commands by default. If you are not satisfied, you can re-customize it via environment variables: ``sh

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --bash)"
```

Alternatively, just add new commands on top of the default list.

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --bash)"
```



In addition, the environment variable `$AUTO_PAGER_MIN_HEIGHT` can control the minimum number of lines to start paging: ``sh

```sh
AUTO_PAGER_MIN_HEIGHT=30 # Paging starts only when the output exceeds 30 lines
eval "$(ap --bash)"
```

If configured as a negative number, the percentage representing the height of the terminal window.

```sh
AUTO_PAGER_MIN_HEIGHT='-50' # Paging starts only when the output exceeds 50% of the terminal window height
eval "$(ap --bash)"
```

If you don't specify `AUTO_PAGER_MIN_HEIGHT`, the default is `-80`, i.e. `80%`.

### fish

Add the following to your `~/.config/fish/config.fish`.

```sh
ap --fish | source
```

ap wrap a bunch of commands by default. If you are not satisfied, you can re-customize it via environment variables.

```sh
set AUTO_PAGER_CMDS go cargo make
ap --fish | source
```

Or you can just add new commands on top of the default list: ```sh

```sh
set AUTO_PAGER_CMDS_EXTRA ps last
ap --fish | source
```

Alternatively, ap can also work with [grc](https://github.com/garabik/grc), and on macOS you can install grc using Homebrew:

```sh
brew install grc
```

ap + grc wrap a bunch of commands by default. If you are not satisfied, you can re-customize it via environment variables: ``sh

```sh
set AUTO_PAGER_CMDS_WITH_GRC ps last dig diff
ap --fish | source
```

Or you can just add new commands on top of the default list: ```sh

```sh
set AUTO_PAGER_CMDS_WITH_GRC_EXTRA ps last
ap --fish | source
```



In addition, the environment variable `$AUTO_PAGER_MIN_HEIGHT` controls the minimum number of lines to start paging: ```sh

```sh
set AUTO_PAGER_MIN_HEIGHT 30 # Paging starts only when the output exceeds 30 lines
ap --fish | source
```

If configured as a negative number, it represents the percentage of terminal window height: ```sh

```sh
set AUTO_PAGER_MIN_HEIGHT -50 # Paging does not start until the output exceeds 50% of the terminal window height
ap --fish | source
```

If ``AUTO_PAGER_MIN_HEIGHT`` is not specified, the default is ``-80``, i.e. ``80%``.

### zsh

Add the following to your `~/.zshrc`.

```sh
eval "$(ap --zsh)"
```

ap wrap a bunch of commands by default. If you are not satisfied, you can re-customize it via environment variables.

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --zsh)"
```

Or you can just add new commands on top of the default manifest: ```sh

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --zsh)"
```

Also, ap works with [grc](https://github.com/garabik/grc), which can be installed on macOS using Homebrew:

```sh
brew install grc
```

ap + grc wrap a bunch of commands by default. If you are not satisfied, you can re-customize it via environment variables: ``sh

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --zsh)"
```

Or you can just add new commands on top of the default list.

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --zsh)"
```



In addition, the environment variable `$AUTO_PAGER_MIN_HEIGHT` can control the minimum number of lines to start paging.

```sh
AUTO_PAGER_MIN_HEIGHT=30 # Paging starts only when the output exceeds 30 lines
eval "$(ap --zsh)"
```

If configured as a negative number, the percentage representing the height of the terminal window.

```sh
AUTO_PAGER_MIN_HEIGHT='-50' # Paging starts only when the output exceeds 50% of the terminal window height
eval "$(ap --zsh)"
```

If `AUTO_PAGER_MIN_HEIGHT` is not specified, the default is `-80`, i.e. `80%`.

## Use

After ap wrap the command just use it as usual.
If there is too much output, the pager specified by the environment variable `$PAGER` will be called automatically to paginate.

If you have a special `$PAGER` variable that does not work with ap, you can also set a separate pager for ap using the environment variable `$AP_PAGER`
to set a separate pager for the ap.

If neither `$AP_PAGER` nor `$PAGER` is specified, then `less -Fr` will be used.

The paginator is not started in the following cases.
* When the output is too small, see `$AUTO_PAGER_MIN_HEIGHT`.
* When ap detects that the command output contains the `ESC [?1049h` sequence, the command is determined to be a full-screen application.
* When the command has not finished executing. The pager will start only after the command execution is completed.
    - Such commands as `ping` and `tcpdump` require `Ctrl-C` to terminate the command before paging is started.
    - Commands such as `python` and `gdb` need to wait for the command to exit before paging is started.

## Frequently Asked Questions

* Does it affect the color output of the command?
    - No.
* Some commands detect the terminal and provide different output for terminal and non-terminal modes, will it change its output?
    - No.
* How do I diagnose if I suspect ap is affecting the output of a command?
    - You can use `command foo` to execute `foo`, which will not invoke ap.
* If I habitually put `| less` after a command that ap has passed, will that cause problems?
    - No.
* Can an ap-passed command still redirect its output?
    - Yes.
* Will the auto-completion of ap-executed commands be broken?
    - No.
* Does ap support interactive applications like `python`, `gdb`?
    - Yes. However, since these applications output some control characters, the paginated content may be a bit messy.
* Does ap support full-screen applications like `htop`, `vim`?
    - What can I say, it won't be wrong anyway, but I can't see the practical point of putting ap with them.
