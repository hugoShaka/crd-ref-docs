{{- define "gvList" -}}
{{- $groupVersions := . -}}

<!--Generated documentation. Please do not edit.-->

<p>Packages:</p>
<ul>
    {{- range $groupVersions }}
    <li>
        {{ htmlRenderGVLink . }}
    </li>
    {{- end }}
</ul>

{{ range $groupVersions }}
{{ template "gvDetails" . }}
{{ end }}

{{- end -}}