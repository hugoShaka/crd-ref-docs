{{- define "gvDetails" -}}
{{- $gv := . -}}
<h2 id="{{htmlGroupVersionID $gv | htmlRenderAnchorID }}">
    {{- $gv.GroupVersionString -}}
</h2>

<p>{{ $gv.Doc }}</p>

{{- if $gv.Kinds  }}
<p>Resource Types:</p>
<ul>
    {{- range $gv.SortedKinds }}
    <li>
        {{ $gv.TypeForKind . | htmlRenderTypeLink }}
    </li>
    {{- end }}
</ul>
{{ end }}

{{ range $gv.SortedTypes }}
{{ template "type" . }}
{{ end }}
{{- end -}}
