{{/*
  {{ include "surge.rule" <groupName: string> [loadRuleSet: bool] <content: string> }}
  example:
  {{ import "rule.conf" | include "surge.rule" "GroupName" }}
  {{ include "rule.rule" "GroupName" (import "rule.conf") }}
  {{ import "rule.conf" | include "surge.rule" "GroupName" true }}
*/}}
{{- define "surge.rule" }}
  {{- if len . | eq 2 }}
    {{- include "surge.rule.range" (index . 0) false (index . 1) }}
  {{- else if len . | eq 3 }}
    {{- include "surge.rule.range" (index . 0) (index . 1) (index . 2) }}
  {{- else }}
    {{- printf "error args: %v" . | fail }}
  {{- end }}
{{- end }}

{{/*
  {{ include "surge.rule.range" <groupName: string> <loadRuleSet: bool> <content: string> }}
  example:
  {{ import "rule.conf" | include "surge.rule.range" "DIRECT" true }}
*/}}
{{- define "surge.rule.range" }}
  {{- $groupName := index . 0 }}
  {{- $loadRuleSet := index . 1 }}
  {{- $list := index . 2 | splitList "\n" }}
  {{- range $value := $list }}
    {{- $v := trim $value }}
    {{- $args := $v | split "," }}
    {{- if $v | len | eq 0 }}
      {{- print "\n" }}
    {{- else if $v | hasPrefix "#" }}
      {{- printf "%s\n" $v }}
    {{- else if and ($args | len | ne 2) ($args | len | ne 3) }}
      {{- print "unkonw line: " $value | fail }}
    {{- else if $loadRuleSet | and ($v | upper | hasPrefix "RULE-SET") | and ($v | lower | contains "http") }}
      {{- import $args._1 | trim | include "surge.rule.range" $groupName true  }}
    {{- else if $args | len | eq 2 }}
      {{- printf "%s,%s,%s\n" $args._0 $args._1 $groupName }}
    {{- else if $args | len | eq 3 }}
      {{- printf "%s,%s,%s,%s\n" $args._0 $args._1 $groupName $args._2 }}
    {{- end }}
  {{- end }}
{{- end }}