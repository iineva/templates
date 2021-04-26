{{/*
*/}}
{{- define "sub" }}
  {{- if typeIs "string" . }}
    {{- import . | b64dec | splitList "\n" | include "sub.list" "default" }}
  {{- else }}
    {{- $groupName := (index . 0) }}
    {{- $content := (index . 1) }}
    {{- import $content | b64dec | splitList "\n" | include "sub.list" $groupName }}
  {{- end }}
{{- end }}

{{/*
*/}}
{{- define "sub.list" }}
  {{- $groupName := (index . 0) }}
  {{- $list := (index . 1) }}
  {{- range $v := $list }}
    {{- include "sub.item" $groupName $v }}
  {{- end }}
{{- end }}

{{/*
*/}}
{{- define "sub.item" }}
  {{- $groupName := (index . 0) }}
  {{- $value := index . 1 | trim }}
  {{- if $value | lower | hasPrefix "ss" }}
    {{- include "sub.ss2surge" $groupName $value }}
  {{- else if $value | lower | hasPrefix "vmess" }}
    {{- include "sub.vmess2surge" $groupName $value }}
  {{- else if $value | lower | hasPrefix "trojan" }}
    {{- include "sub.trojan2surge" $groupName $value }}
  {{- else if len . | eq 0 }}
    {{- print "\n" }}
  {{- else if hasPrefix "REMARKS" $value }}
    {{- setValue (printf "sub.title:%s" $groupName) (query "REMARKS" $value) }}
  {{- else if hasPrefix "STATUS" $value }}
    {{- setValue (printf "sub.status:%s" $groupName) (query "STATUS" $value) }}
  {{- else if hasPrefix "#" $value }}
    {{- printf "$s\n" $value }}
  {{- else if contains "://" $value }}
    {{- printf "# know protocol: %s\n" $value }}
  {{- end }}
{{- end }}

{{/*
*/}}
{{- define "sub.ss2surge" }}
  {{- $groupName := (index . 0) }}
  {{- $raw := (index . 1) }}
  {{- with $u := fromUrl $raw }}
  {{- with $remark := $u.GoURL.Fragment }}
  {{- with $uri := b64dec $u.Host | printf "ss://%s" | fromUrl }}
    {{- addList (print "sub.remarks:" $groupName) $remark }}
    {{- printf "%s = ss, %v, %v, encrypt-method=%s, password=%s\n" $remark $uri.Hostname $uri.Port $uri.Username $uri.Password }}
  {{- end }}
  {{- end }}
  {{- end }}
{{- end }}

{{/*
*/}}
{{- define "sub.vmess2surge" }}
  {{- $groupName := (index . 0) }}
  {{- $raw := (index . 1) }}
  {{- with $obj := replace "vmess://" "" $raw | b64dec | fromJson }}
    {{- addList (print "sub.remarks:" $groupName) $obj.ps }}
    {{- if get $obj "net"  | eq "ws" }}{{ addList "sub.v2ray" "ws=true" }}{{ end }}
    {{- if get $obj "path" | ne "" }}{{ addList "sub.v2ray" (printf "ws-path=%s" $obj.path) }}{{ end }}
    {{- if get $obj "host" | ne "" }}{{ addList "sub.v2ray" (printf "ws-headers=HOST:%s" $obj.host) }}{{ end }}
    {{- if get $obj "tls" | eq "tls" }}{{ addList "sub.v2ray" "tls=true" }}{{ end }}
    {{- if get $obj "type" | ne "" | and (queryContext "type" | eq "shadowrocket") }}{{ addList "sub.v2ray" (printf "encrypt-method=%s" $obj.type) }}{{ end }}
    {{- printf "%s = vmess, %s, %s, username=%s, %s\n" $obj.ps $obj.add $obj.port $obj.id (getList "sub.v2ray" | join ", ") }}
    {{- delList "sub.v2ray"}}
  {{- end }}
{{- end }}

{{/*
*/}}
{{- define "sub.trojan2surge" }}
  {{- $groupName := (index . 0) }}
  {{- $raw := (index . 1) }}
  {{- with $uri := fromUrl $raw }}
  {{- addList (print "sub.remarks:" $groupName) $uri.GoURL.Fragment }}
  {{- if eq $uri.Password "" }}
    {{- printf "%s = trojan, %s, %v, password=%s\n" $uri.GoURL.Fragment $uri.Hostname $uri.Port $uri.Username }}
  {{- else }}
    {{- printf "%s = trojan, %s, %v, username=%s, password=%s\n" $uri.GoURL.Fragment $uri.Hostname $uri.Port $uri.Username $uri.Password }}
  {{- end }}
  {{- end }}
{{- end }}