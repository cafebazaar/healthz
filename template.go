package healthz

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="refresh" content="10">
		<title>Overall Health: {{.Root.OverallHealth|HealthTitle}}</title>
		<link rel="stylesheet" href="min.css">
	</head>
	<body>
		<div class="container">
			<div class="overall hl{{.Root.OverallHealth}}">
				<h1>Overall: {{.Root.OverallHealth|HealthTitle}}</h1>
				<p style="text-align: right"><a href="json">JSON</a>, <a href="liveness">Liveness</a>, <a href="readiness">Readiness</a></p>
			</div>
			<h3>Service: <span class="service">{{ .ServiceSignature }}</span></h3>
			<h3>Hostname: <span class="hostname">{{ .Hostname }}</span></h3>
			<h3>Uptime: <span class="uptime">{{ .Uptime }}</span></h3>
			{{ range .Root.Subcomponents }}<div class="row hl{{.OverallHealth}}">
				<div class="col c4"><strong>{{ .Name }}</strong></div>
				<div class="col c4">{{.OverallHealth|HealthTitle}}</div>
				<div class="col c4">{{ .Severity|SeverityTitle }}</div>
				{{ if .Subcomponents}}<div class="subcomponents">
					{{ range .Subcomponents }}<div class="row hl{{.OverallHealth}}">
						<div class="col c4"><strong>{{ .Name }}</strong></div>
						<div class="col c4">{{.OverallHealth|HealthTitle}}</div>
						<div class="col c4">{{ .Severity|SeverityTitle }}</div>
					</div>{{ end }}
				</div>{{ end }}
			</div>{{ else }}
				<div><strong>No components is registered, or showing details is disabled.</strong></div>
			{{ end }}
	</body>
</html>`
