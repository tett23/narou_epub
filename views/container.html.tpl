{{define "container"}}
<head>
  <title>NarouEpub</title>
</head>
<body>
  <h1>{{.Container.Title}}</h1>

  <dl>
    <dt>NCode</dt>
    <dd><a href="https://ncode.syosetu.com/{{.Container.NCode}}/" target="_blank">{{.Container.NCode}}</a></dd>
    <dt>Author</dt>
    <dd><a href="https://mypage.syosetu.com/{{.Container.UserID}}/" target="_blank">{{.Container.Author}}</a></dd>
    <dt>UpdatedAt</dt>
    <dd>{{.Container.UpdatedAt}}</dd>
    <dt>GeneralAllNo</dt>
    <dd>{{.Container.GeneralAllNo}}</dd>
  </dl>

  <p>
    <span><a href="/containers/{{.Container.NCode}}/fetch">epub生成</a></span>
    <span><a href="/containers/{{.Container.NCode}}/publish">publish</a></span>
  </p>

  <table>
    <thead>
      <tr>
        <th>EpisodeNumber</th>
        <th>EpisodeTitle</th>
        <th>UpdatedAt</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
    {{range $i, $v := .Container.Episodes}}
      <tr>
        <td>{{$v.EpisodeNumber}}</td>
        <td><a href="https://ncode.syosetu.com/{{$v.NCode}}/{{$v.EpisodeNumber}}/" target="_blank">{{$v.EpisodeTitle}}</a></td>
        <td>{{$v.UpdatedAt}}</td>
        <td>
          <span><a href="/containers/{{$v.NCode}}/episode/{{$v.EpisodeNumber}}/fetch">epub生成</a></span>
          <span><a href="/containers/{{$v.NCode}}/episode/{{$v.EpisodeNumber}}/publish">publish</a></span>
          <span><a href="/containers/{{$v.NCode}}/episode/{{$v.EpisodeNumber}}">JSON</a></span>
        </td>
      </tr>
    {{end}}
    </tbody>
  </table>

</body>
{{end}}
