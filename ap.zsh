auto-pager-wrapper() {
    local name=$funcstack[2]
    if (( ! $+commands[$name] )) ; then
        return
    fi
    ap $name $@
}

for cmd in $AUTO_PAGER_CMDS; do
    $cmd() { auto-pager-wrapper $@ }
done

unset cmd
