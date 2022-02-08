auto-pager-wrapper() {
    local grc=$1; shift
    local name=$funcstack[2]
    if ! type -p $name >/dev/null; then
        return
    fi

    local height="${AUTO_PAGER_MIN_HEIGHT:--80}"
    if [[ "x$grc" == "xgrc" ]] && type -p grc >/dev/null; then
        ap --height $height grc $name $@
    else
        ap --height $height $name $@
    fi
}

AUTO_PAGER_CMDS_DEFAULT_GRC=(
    df du env id last lsof mount ps sysctl
    diff tar
    cc gcc g++ make
    curl dig ifconfig iostat ip iptables iptables-save netstat
    ping ping6 tcpdump traceroute traceroute6 whois
    docker docker-compose docker-machine kubectl
)

AUTO_PAGER_CMDS_DEFAULT=(
    ls tree
    go cargo rustc
)

if [[ ${#AUTO_PAGER_CMDS[@]} -gt 0 ]]; then
    for cmd in $AUTO_PAGER_CMDS; do
        $cmd() { auto-pager-wrapper - $@ }
    done
else
    for cmd in $AUTO_PAGER_CMDS_DEFAULT; do
        $cmd() { auto-pager-wrapper - $@ }
    done
    for cmd in $AUTO_PAGER_CMDS_DEFAULT_GRC; do
        $cmd() { auto-pager-wrapper grc $@ }
    done
fi

for cmd in $AUTO_PAGER_CMDS_EXTRA; do
    $cmd() { auto-pager-wrapper - $@ }
done

for cmd in $AUTO_PAGER_CMDS_WITH_GRC; do
    $cmd() { auto-pager-wrapper grc $@ }
done

unset cmd
