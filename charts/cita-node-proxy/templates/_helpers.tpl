{{/*
Expand the name of the chart.
*/}}
{{- define "cita-node-proxy.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cita-node-proxy.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cita-node-proxy.labels" -}}
helm.sh/chart: {{ include "cita-node-proxy.chart" . }}
{{ include "cita-node-proxy.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cita-node-proxy.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cita-node-proxy.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
