AUTO_PAGER_CMDS=(
    ls tree
    go cargo rustc
)

AUTO_PAGER_CMDS_WITH_GRC=(
    df du env id last lsof mount ps sysctl
    diff tar
    cc gcc g++ make
    curl dig ifconfig iostat ip iptables iptables-save netstat
    ping ping6 tcpdump traceroute traceroute6 whois
    docker docker-compose docker-machine kubectl
)

for _cmd in "${AUTO_PAGER_CMDS[@]}" "${AUTO_PAGER_CMDS_EXTRA[@]}"; do
    [ "$_cmd" ] && alias $_cmd="ap $_cmd"
done

if type -p grc >/dev/null; then
    for _cmd in "${AUTO_PAGER_CMDS_WITH_GRC[@]}" "${AUTO_PAGER_CMDS_WITH_GRC_EXTRA[@]}"; do
        [ "$_cmd" ] && alias $_cmd="ap grc $_cmd"
    done
fi

unset _cmd _prog
