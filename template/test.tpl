{{- define "test" -}}
<html>
    <head>
        <title></title>
    </head>
    <body>
        <form action="/login" method="post">
            用户名:<input type="text" name="username" value={{ tToLower }}>
            密码:<input type="password" name="password">
            <input type="submit" value="登陆">
            <input type="hidden" name="token" value="{{ .Token }}">
            {{- with .Week }}
            {{- range . }}
            <input type="checkbox" name="date" value="{{ .Data|tToLower }}">{{ .Data|tToLower }}
            {{- end }}
            {{- end }}
        </form>
    </body>
</html>
{{end}}