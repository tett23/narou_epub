{{define "index"}}
<head>
  <title>NarouEpub</title>
</head>
<body>
  <h1>NarouEpub</h1>

  {{range $i, $v := .Containers}}
  <h2>
    <a href="/containers/{{$v.NCode}}">{{$v.Title}}</a>
  </h2>

  <dl>
    <dt>NCode</dt>
    <dd><a href="https://ncode.syosetu.com/{{$v.NCode}}/" target="_blank">{{$v.NCode}}</a></dd>
    <dt>Author</dt>
    <dd><a href="https://mypage.syosetu.com/{{$v.UserID}}/" target="_blank">{{$v.Author}}</a></dd>
    <dt>UpdatedAt</dt>
    <dd>{{$v.UpdatedAt}}</dd>
    <dt>GeneralAllNo</dt>
    <dd>{{$v.GeneralAllNo}}</dd>
  </dl>
  {{end}}
</body>
{{end}}
