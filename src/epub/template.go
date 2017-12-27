package epub

const htmlTemplate = `
{{define "base"}}
<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="ja" lang="ja" class="vrtl">
  <head>
    <title>{{.title}}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <link href="../stylesheet.css" rel="stylesheet" type="text/css"/>
    <link href="../page_styles.css" rel="stylesheet" type="text/css"/>
  </head>
  <body id="E9OE0-ab176f2d8da64abb8d9f8e088e8b6a8f" class="calibre">
    <h2 class="calibre7" id="calibre_pb_0">{{.title}}</h2>
    {{range .lines}}
      {{.}}<p class="calibre6" style="margin:0pt; border:0pt; height:0pt"> </p>
    {{end}}
  </body>
</html>
{{end}}
`
