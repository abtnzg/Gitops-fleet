{{/*
Expand the name of the chart.
*/}}
{{- define "service-a.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some K8s name fields are limited.
*/}}
{{- define "service-a.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Chart name and version label value.
*/}}
{{- define "service-a.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels — applied to every resource.
*/}}
{{- define "service-a.labels" -}}
helm.sh/chart: {{ include "service-a.chart" . }}
{{ include "service-a.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: {{ .Values.partOf | default "pilot-team-platform-validation" }}
{{- end }}

{{/*
Selector labels — used in matchLabels. Must be stable across upgrades.
*/}}
{{- define "service-a.selectorLabels" -}}
app.kubernetes.io/name: {{ include "service-a.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: workload
{{- end }}

{{/*
Service account name.
*/}}
{{- define "service-a.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "service-a.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Image reference. Always rendered as <registry>/<repo>:<tag>.
If registry is empty, use Docker Hub (most common case in dev).
*/}}
{{- define "service-a.image" -}}
{{- if .Values.image.registry }}
{{- printf "%s/%s:%s" .Values.image.registry .Values.image.repository .Values.image.tag }}
{{- else }}
{{- printf "%s:%s" .Values.image.repository .Values.image.tag }}
{{- end }}
{{- end }}
