{{- define "letsagree.labels" -}}
app.kubernetes.io/managed-by: Helm
meta.helm.sh/release-name: {{ .Release.Name }}
meta.helm.sh/release-namespace: {{ .Release.Namespace }}
{{- end }}

{{- define "letsagree.backend.image" -}}
{{ .Values.backend.image }}:{{ .Chart.AppVersion }}
{{- end }}

{{- define "letsagree.frontend.image" -}}
{{ .Values.frontend.image }}:{{ .Chart.AppVersion }}
{{- end }}
