auto-pager-wrapper() {
    local grc=$1; shift
    local name=$funcstack[2]
    if (( ! $+commands[$name] )) ; then
        return
    fi
    if [[ "x$grc" == "xgrc" ]]; then
        ap grc $name $@
    else
        ap $name $@
    fi
}

AUTO_PAGER_CMDS_DEFAULT=(
    go make
    ps lsof netstat
)

if [[ ${#AUTO_PAGER_CMDS[@]} -eq 0 ]]; then
    AUTO_PAGER_CMDS=($AUTO_PAGER_CMDS_DEFAULT)
fi

for cmd in $AUTO_PAGER_CMDS; do
    $cmd() { auto-pager-wrapper - $@ }
done

for cmd in $AUTO_PAGER_CMDS_EXTRA; do
    $cmd() { auto-pager-wrapper - $@ }
done

for cmd in $AUTO_PAGER_CMDS_WITH_GRC; do
    $cmd() { auto-pager-wrapper grc $@ }
done

unset cmd
