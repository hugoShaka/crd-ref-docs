{{- define "type" -}}
{{- $type := . -}}
{{- if htmlShouldRenderType $type -}}

<h3 id="{{ htmlTypeID $type | htmlRenderAnchorID }}">
    {{- $type.Name  }}
    {{ if $type.IsAlias }}(<code>{{ htmlRenderTypeLink $type.UnderlyingType  }}</code> alias) {{ end -}}
</h3>

{{ if $type.References -}}
<p>
    (<em>Appears In: </em>
    {{- $prev := "" -}}
    {{- range $type.SortedReferences }}
        {{- if $prev -}}, {{ end -}}
        {{- $prev = . -}}
        {{ htmlRenderTypeLink . }}
    {{- end }}
    )
</p>
{{- end }}

{{ $type.Doc | htmlRenderDoc }}

{{ if $type.Members -}}
<table>
    <thead>
        <tr>
            <th>Field</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        {{ if $type.GVK -}}
        <tr>
            <td>
                <code>apiVersion</code><br/>
                string
            </td>
            <td>
                <code>
                    {{ $type.GVK.Group }}/{{ $type.GVK.Version }}
                </code>
            </td>
        </tr>
        <tr>
            <td>
                <code>kind</code><br/>
                string
            </td>
            <td>
                <code>
                    {{ $type.GVK.Kind }}
                </code>
            </td>
        </tr>
        {{ end -}}
        {{ range $type.Members -}}
        {{ template "type_members" . }}
        {{ end -}}
    </tbody>
</table>
{{ end -}}


{{- end -}}
{{- end -}}