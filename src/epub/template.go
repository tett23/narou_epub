package epub

const htmlTemplate = `{{define "base"}}<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="ja" lang="ja" class="vrtl">
  <head>
    <title>{{.title}}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <link href="../stylesheet.css" rel="stylesheet" type="text/css"/>
    <link href="../page_styles.css" rel="stylesheet" type="text/css"/>
  </head>
  <body id="main" class="calibre">
    {{range .preface}}
      {{.}}<p class="calibre6" style="margin:0pt; border:0pt; height:0pt"> </p>
    {{end}}

    <h2 class="calibre7" id="calibre_pb_0">{{.title}}</h2>
    {{range .body}}
      {{.}}<p class="calibre6" style="margin:0pt; border:0pt; height:0pt"> </p>
    {{end}}

    {{if .postscript}}
      <hr />
      {{range .postscript}}
        {{.}}<p class="calibre6" style="margin:0pt; border:0pt; height:0pt"> </p>
      {{end}}
    {{end}}
  </body>
</html>
{{end}}
`

const overviewTemplate = `{{define "overview"}}<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="ja" lang="ja" class="vrtl">
  <head>
    <title>{{.title}}</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <link href="../stylesheet.css" rel="stylesheet" type="text/css"/>
    <link href="../page_styles.css" rel="stylesheet" type="text/css"/>
  </head>
  <body id="main" class="calibre">
    <h2 class="calibre7" id="calibre_pb_0">{{.title}}</h2>
    {{if .episodeTitle}}
    <h3 class="calibre3">{{.episodeTitle}}</h3>
    {{end}}
    <p class="calibre6">NCode {{.nCode}}</p>
    <p class="calibre6">Author {{.author}}</p>
    <p class="calibre6">CreatedAt {{.date}}</p>

    <nav id="toc" xmlns:epub="http://www.idpf.org/2007/ops" epub:type="toc">
      <ol>
      {{range $i, $v := .items}}
        <li><a href="section_{{$v.EpisodeNumber}}.html">{{$v.Name}}</a></li>
      {{end}}
      </ol>
    </nav>

    <p class="pagebreak"></p>
  </body>
</html>
{{end}}
`

const contentOpfTemplate = `{{define "opf"}}<?xml version='1.0' encoding='utf-8'?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="uuid_id" version="2.0">
  <metadata xmlns:calibre="http://calibre.kovidgoyal.net/2009/metadata" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:opf="http://www.idpf.org/2007/opf" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <dc:date>{{.date}}</dc:date>
    <!-- <dc:publisher></dc:publisher> -->
    <!-- <dc:rights></dc:rights> -->
    <dc:creator opf:file-as="Unknown" opf:role="aut">{{.author}}</dc:creator>
    <dc:language>ja</dc:language>
    <dc:identifier id="uuid_id" opf:scheme="uuid">{{.uuid}}</dc:identifier>
    <dc:title>{{.title}}</dc:title>
    <meta content="vertical-rl" name="primary-writing-mode"/>
  </metadata>
  <manifest>
    <item href="page_styles.css" id="page_css" media-type="text/css"/>
    <item href="stylesheet.css" id="css" media-type="text/css"/>
    <item href="toc.ncx" id="ncx" media-type="application/x-dtbncx+xml"/>
    {{range $i, $v := .items}}
    {{if eq $v.Path "body/overview.html"}}
    <item href="{{$v.Path}}" properties="nav" id="id_{{$v.Order}}" media-type="application/xhtml+xml"/>
    {{else}}
    <item href="{{$v.Path}}" id="id_{{$v.Order}}" media-type="application/xhtml+xml"/>
    {{end}}
    {{end}}
  </manifest>
  <spine toc="ncx" page-progression-direction="rtl">
    {{range $i, $v := .items}}
    <itemref idref="id_{{$v.Order}}"/>
    {{end}}
  </spine>
  <guide/>
</package>
{{end}}
`

const tocNcxTemplate = `{{define "ncx"}}<?xml version='1.0' encoding='utf-8'?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="jpn">
  <head>
    <meta content="{{.uuid}}" name="dtb:uid"/>
    <meta content="2" name="dtb:depth"/>
    <meta content="narou_epub" name="dtb:generator"/>
    <meta content="0" name="dtb:totalPageCount"/>
    <meta content="0" name="dtb:maxPageNumber"/>
  </head>
  <docTitle>
    <text>{{.title}}</text>
  </docTitle>
  <navMap>
    {{range $i, $v := .items}}
    <navPoint class="chapter" id="id_{{$v.Order}}" playOrder="{{$v.Order}}">
      <navLabel>
        <text>{{$v.Name}}</text>
      </navLabel>
      <content src="{{$v.Path}}#main"/>
    </navPoint>
    {{end}}
  </navMap>
</ncx>
{{end}}
`
