declare -p AUTO_PAGER_CMDS >/dev/null 2>&1 || AUTO_PAGER_CMDS=(
    ls tree
    go cargo rustc
)

declare -p AUTO_PAGER_CMDS_WITH_GRC >/dev/null 2>&1 || AUTO_PAGER_CMDS_WITH_GRC=(
    df du env id last lsof mount ps sysctl
    diff tar
    cc gcc g++ make
    curl dig ifconfig iostat ip iptables iptables-save netstat
    ping ping6 tcpdump traceroute traceroute6 whois
    docker docker-compose docker-machine kubectl
)

# can't use alias like ap.bash because it breaks zsh auto-completion for these commands.
for _cmd in "${AUTO_PAGER_CMDS[@]}" "${AUTO_PAGER_CMDS_EXTRA[@]}"; do
    [ "$_cmd" ] || continue
    _alias=`alias $_cmd`
    _alias=${_alias#*\'}
    _alias=${_alias%\'}
    [ "$_alias" ] && unalias $_cmd || _alias=$_cmd
    eval "$_cmd() { ap $_alias \"\$@\" }"
done

type -p grc >/dev/null && _grc="grc "

for _cmd in "${AUTO_PAGER_CMDS_WITH_GRC[@]}" "${AUTO_PAGER_CMDS_WITH_GRC_EXTRA[@]}"; do
    [ "$_cmd" ] || continue
    _alias=`alias $_cmd`
    _alias=${_alias#*\'}
    _alias=${_alias%\'}
    [ "$_alias" ] && unalias $_cmd || _alias=$_cmd
    eval "$_cmd() { ap $_grc$_alias \"\$@\" }"
done

unset _alias _cmd _grc
