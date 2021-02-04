{{- define "type_members" -}}
{{- $field := . -}}
<tr>
    <td>
        <code>{{ .Name }}</code><br/>
        <em>{{ htmlRenderType .Type }}</em>
    </td>
    <td>
        {{ if $field.Optional -}}
        <em>(Optional)</em>
        {{- end -}}
        {{- if eq $field.Name "metadata" }}
        <p>Refer to Kubernetes API documentation for fields of `metadata`.</p>
        {{ else -}}
        {{ $field.Doc | htmlRenderDoc}}
        {{- if eq $field.Name "spec" -}}
        </br>
        <table>
            {{ range $field.Type.Fields -}}
            {{ template "type_members" . }}
            {{- end -}}
        </table>
        {{- end -}}
    </td>
</tr>
{{- end -}}
{{- end -}}
