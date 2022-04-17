# ap --`auto-pager`

ap est un utilitaire shell qui permet à la sortie d'autres commandes shell d'entrer automatiquement en mode de retournement de page interactif.

ap se compose de deux parties, un programme binaire écrit en Go qui capture la sortie des commandes et prend en charge le retournement de page, et un ensemble de scripts shell qui créent un wrapper portant le même nom pour une liste de commandes spécifiée par l'utilisateur.

L'utilisation des commandes après le bouclage est la même qu'auparavant et ne devrait pas changer les habitudes de l'utilisateur ni causer de problèmes.

\*Lire ceci dans d'autres langues :[Anglais](README_en.md),[Chinois simplifié](README.md)

ap est un outil shell qui permet à la sortie d'autres commandes shell d'entrer automatiquement en mode de pagination interactive.

ap se compose de deux parties, un programme binaire écrit en langage Go, chargé de capturer la sortie de la commande et de prendre en charge le changement de page,
et un ensemble de scripts shell responsables de la création d'un wrapper du même nom pour une liste de commandes spécifiée par l'utilisateur.

L'utilisation de la commande après le bouclage est la même que celle d'origine, et les habitudes de fonctionnement de l'utilisateur ne doivent pas être modifiées, et cela ne causera pas de confusion à l'utilisateur.

## Installer

    go install github.com/flw-cn/ap@master

## configurer

-   `bash`(#frapper)
-   `fish`(#`fish`)
-   `zsh`(#`zsh`)

### frapper

À votre`~/.bashrc`Ajoutez-y ce qui suit :

```sh
eval "$(ap --bash)"
```

ap encapsule un lot de commandes par défaut. Si vous n'êtes pas satisfait, vous pouvez le re-personnaliser via des variables d'environnement :

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --bash)"
```

Ou ajoutez simplement de nouvelles commandes en haut de la liste par défaut :

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --bash)"
```

De plus, ap peut également être utilisé avec[grc](https://github.com/garabik/grc)En travaillant ensemble, grc peut être installé en utilisant Homebrew sous macOS :

```sh
brew install grc
```

ap + grc encapsule un lot de commandes par défaut. Si vous n'êtes pas satisfait, vous pouvez le re-personnaliser via des variables d'environnement :

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --bash)"
```

Ou ajoutez simplement de nouvelles commandes en haut de la liste par défaut :

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --bash)"
```

De plus, les variables d'environnement`$AUTO_PAGER_MIN_HEIGHT`Vous pouvez contrôler le nombre minimum de lignes pour démarrer la pagination :

```sh
AUTO_PAGER_MIN_HEIGHT=30        # 输出超过 30 行时才开始分页
eval "$(ap --bash)"
```

S'il est configuré comme un nombre négatif, il représente un pourcentage de la hauteur de la fenêtre du terminal :

```sh
AUTO_PAGER_MIN_HEIGHT='-50'     # 输出超过终端窗口高度的 50% 时才开始分页
eval "$(ap --bash)"
```

Si non spécifié`AUTO_PAGER_MIN_HEIGHT`,La valeur par défaut est`-80`,Tout de suite`80%`。

### `fish`

À votre`~/.config/fish/config.fish`Ajoutez-y ce qui suit :

```sh
ap --fish | source
```

ap encapsule un lot de commandes par défaut. Si vous n'êtes pas satisfait, vous pouvez le re-personnaliser via des variables d'environnement :

```sh
set AUTO_PAGER_CMDS go cargo make
ap --fish | source
```

或者也可只在默认清单之上增加新的命令：

```sh
set AUTO_PAGER_CMDS_EXTRA ps last
ap --fish | source
```

De plus, ap peut également être utilisé avec[grc](https://github.com/garabik/grc)En travaillant ensemble, grc peut être installé en utilisant Homebrew sous macOS :

```sh
brew install grc
```

ap + grc encapsule un lot de commandes par défaut. Si vous n'êtes pas satisfait, vous pouvez le re-personnaliser via des variables d'environnement :

```sh
set AUTO_PAGER_CMDS_WITH_GRC ps last dig diff
ap --fish | source
```

Ou ajoutez simplement de nouvelles commandes en haut de la liste par défaut :

```sh
set AUTO_PAGER_CMDS_WITH_GRC_EXTRA ps last
ap --fish | source
```

De plus, les variables d'environnement`$AUTO_PAGER_MIN_HEIGHT`Vous pouvez contrôler le nombre minimum de lignes pour démarrer la pagination :

```sh
set AUTO_PAGER_MIN_HEIGHT 30        # 输出超过 30 行时才开始分页
ap --fish | source
```

S'il est configuré comme un nombre négatif, il représente un pourcentage de la hauteur de la fenêtre du terminal :

```sh
set AUTO_PAGER_MIN_HEIGHT -50     # 输出超过终端窗口高度的 50% 时才开始分页
ap --fish | source
```

Si non spécifié`AUTO_PAGER_MIN_HEIGHT`,La valeur par défaut est`-80`,Tout de suite`80%`。

### `zsh`

À votre`~/.zshrc`Ajoutez-y ce qui suit :

```sh
eval "$(ap --zsh)"
```

ap encapsule un lot de commandes par défaut. Si vous n'êtes pas satisfait, vous pouvez le re-personnaliser via des variables d'environnement :

```sh
AUTO_PAGER_CMDS=(go cargo make)
eval "$(ap --zsh)"
```

Ou ajoutez simplement de nouvelles commandes en haut de la liste par défaut :

```sh
AUTO_PAGER_CMDS_EXTRA=(ps last)
eval "$(ap --zsh)"
```

De plus, ap peut également être utilisé avec[grc](https://github.com/garabik/grc)En travaillant ensemble, grc peut être installé en utilisant Homebrew sous macOS :

```sh
brew install grc
```

ap + grc encapsule un lot de commandes par défaut. Si vous n'êtes pas satisfait, vous pouvez le re-personnaliser via des variables d'environnement :

```sh
AUTO_PAGER_CMDS_WITH_GRC=(ps last dig diff)
eval "$(ap --zsh)"
```

Ou ajoutez simplement de nouvelles commandes en haut de la liste par défaut :

```sh
AUTO_PAGER_CMDS_WITH_GRC_EXTRA=(ps last)
eval "$(ap --zsh)"
```

De plus, les variables d'environnement`$AUTO_PAGER_MIN_HEIGHT`Vous pouvez contrôler le nombre minimum de lignes pour démarrer la pagination :

```sh
AUTO_PAGER_MIN_HEIGHT=30        # 输出超过 30 行时才开始分页
eval "$(ap --zsh)"
```

S'il est configuré comme un nombre négatif, il représente un pourcentage de la hauteur de la fenêtre du terminal :

```sh
AUTO_PAGER_MIN_HEIGHT='-50'     # 输出超过终端窗口高度的 50% 时才开始分页
eval "$(ap --zsh)"
```

Si non spécifié`AUTO_PAGER_MIN_HEIGHT`,La valeur par défaut est`-80`,Tout de suite`80%`。

## utiliser

Les commandes qui ont subi un wrap ap doivent simplement être utilisées comme d'habitude.
S'il y a trop de sortie, la variable d'environnement sera appelée automatiquement`$PAGER`Le téléavertisseur spécifié effectue la pagination.

si votre`$PAGER`Les variables sont spéciales et ne peuvent pas être adaptées à ap. Vous pouvez également passer des variables d'environnement`$AP_PAGER`Réglez le téléavertisseur séparément pour ap.

si`$AP_PAGER`et`$PAGER`ne sont pas spécifiés, alors utilisera`less -Fr`。

Les conditions suivantes ne démarreront pas le téléavertisseur :

-   Lorsque le contenu de sortie est trop petit, voir`$AUTO_PAGER_MIN_HEIGHT`。
-   Lorsque ap détecte que la sortie de la commande contient`ESC [?1049h`séquence, la commande est déterminée comme étant une application plein écran.
-   Lorsque la commande n'a pas encore été exécutée. Le téléavertisseur ne démarre pas tant que l'exécution de la commande n'est pas terminée.
    -   `ping` 和 `tcpdump`et d'autres commandes de ce type, vous devez appuyer sur`Ctrl-C`La pagination ne démarre pas avant la commande terminate.
    -   `python`et`gdb`Ces commandes doivent attendre la fin de la commande avant de démarrer la pagination.

## Problème commun

-   Cela affectera-t-il la sortie colorée de la commande ?
    -   Ne le fera pas.
-   Il existe des commandes qui détectent le terminal et fournissent une sortie différente pour le mode terminal et le mode non terminal, changera-t-il sa sortie ?
    -   Ne le fera pas.
-   Comment puis-je diagnostiquer si je soupçonne que ap affecte la sortie de la commande ?
    -   tu peux l'utiliser`command foo`éxécuter`foo`, de sorte que ap ne soit pas appelé.
-   Si j'ajoute habituellement la commande à la commande ap`| less`, y aura-t-il un problème ?
    -   Ne le fera pas.
-   La commande ap peut-elle toujours rediriger sa sortie ?
    -   Pouvez.
-   La saisie semi-automatique des commandes AP sera-t-elle interrompue ?
    -   Ne le fera pas.
-   ap prend en charge des choses comme`python`、`gdb`Une telle application interactive ?
    -   Support. Mais comme ces applications produisent des caractères de contrôle, le contenu que vous voyez après la pagination peut être un peu désordonné.
-   ap prend en charge des choses comme`htop`、`vim`Une application plein écran comme celle-ci ?
    -   Comment dire, je ne peux pas me tromper de toute façon, mais je ne vois pas l'intérêt pratique d'associer ap avec eux.
