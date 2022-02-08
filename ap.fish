set AUTO_PAGER_CMDS \
    ls tree \
    go cargo rustc

set AUTO_PAGER_CMDS_WITH_GRC \
    df du env id last lsof mount ps sysctl \
    diff tar \
    cc gcc g++ make \
    curl dig ifconfig iostat ip iptables iptables-save netstat \
    ping ping6 tcpdump traceroute traceroute6 whois \
    docker docker-compose docker-machine kubectl

for _cmd in $AUTO_PAGER_CMDS $AUTO_PAGER_CMDS_EXTRA;
    [ "$_cmd" ]; and alias $_cmd "ap $_cmd"
end

if type -p grc >/dev/null;
    for _cmd in $AUTO_PAGER_CMDS_WITH_GRC $AUTO_PAGER_CMDS_WITH_GRC_EXTRA;
        [ "$_cmd" ]; and alias $_cmd "ap grc $_cmd"
    end
end

set -e _cmd _prog
