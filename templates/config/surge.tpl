{{- loadDefine "/define/sub.tpl" -}}
{{- loadDefine "/define/surge.rule.tpl" -}}
#!MANAGED-CONFIG {{ requestURL }} interval=86400 strict=false

[General]
loglevel = info
skip-proxy = 127.0.0.1,192.168.0.0/16,10.0.0.0/8,172.16.0.0/12,100.64.0.0/10,localhost,*.local
internet-test-url = http://baidu.com
test-timeout = 5
proxy-test-url = http://www.google.com/generate_204
ipv6 = false

# https://manual.nssurge.com/others/doh.html
doh-server = https://dns.alidns.com/dns-query,https://9.9.9.9/dns-query
# doh-follow-outbound-mode = true

# iOS
bypass-system = true
bypass-tun = 192.168.0.0/16, 172.16.0.0/12, 100.64.0.0/10, 10.0.0.0/8, 17.0.0.0/8

# macOS
interface = 127.0.0.1
socks-interface = 127.0.0.1

[Proxy]

[Proxy Group]
DIRECT = select, direct

[Rule]
RULE-SET,SYSTEM,DIRECT
RULE-SET,LAN,DIRECT
GEOIP,CN,DIRECT
FINAL,DIRECT

[URL Rewrite]
^http://www.google.cn http://www.google.com header
