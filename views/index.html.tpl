{{define "index"}}
<head>
  <title>NarouEpub</title>
</head>
<body>
  <h1><a href="/">NarouEpub</a></h1>

  <h2>latest</h2>
  <table>
    <thead>
      <tr>
        <th>NCode</th>
        <th>EpisodeNumber</th>
        <th>EpisodeTitle</th>
        <th>CrawledAt</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
      {{range $i,$v :=.Latests}}
      <tr>
        <td>
          <a href="/containers/{{$v.NCode | ToUpper}}">{{$v.NCode | ToUpper}}</a>
        </td>
        <td>{{$v.EpisodeNumber}}</td>
        <td><a href="https://ncode.syosetu.com/{{$v.NCode}}/{{$v.EpisodeNumber}}/" target="_blank">{{$v.EpisodeTitle}}</a></td>
        <td>{{$v.CrawledAt}}</td>
        <td>
          <span>
            <a
              class="post-form"
              data-method="post"
              href="/containers/{{$v.NCode}}/episode/{{$v.EpisodeNumber}}/fetch"
              onclick="return false;"
            >epub生成</a></span>
          <span>
            <a
              class="post-form"
              data-method="post"
              href="/containers/{{$v.NCode}}/episode/{{$v.EpisodeNumber}}/publish"
              onclick="return false;"
            >publish</a></span>
          <span><a href="/containers/{{$v.NCode}}/episode/{{$v.EpisodeNumber}}">JSON</a></span>
        </td>
      </tr>
      {{end}}
    </tbody>
  </table>

  {{range $i, $v := .Containers}}
  <h2>
    <a href="/containers/{{$v.NCode}}">{{$v.Title}}</a>
  </h2>

  <dl>
    <dt>NCode</dt>
    <dd><a href="https://ncode.syosetu.com/{{$v.NCode}}/" target="_blank">{{$v.NCode}}</a></dd>
    <dt>Author</dt>
    <dd><a href="https://mypage.syosetu.com/{{$v.UserID}}/" target="_blank">{{$v.Author}}</a></dd>
    <dt>GeneralLastUp</dt>
    <dd>{{$v.GeneralLastUp}}</dd>
    <dt>UpdatedAt</dt>
    <dd>{{$v.UpdatedAt}}</dd>
    <dt>GeneralAllNo</dt>
    <dd>{{$v.GeneralAllNo}}</dd>
  </dl>
  {{end}}

  <script>
    const links = document.querySelectorAll("a.post-form")
    links.forEach((_, i) => {
      const link = links[i]

      link.addEventListener("click", () => {
        const form = document.createElement("form");
        form.setAttribute("action", link.href);
        form.setAttribute("method", link.dataset.method);
        form.style.display = "none";
        document.body.appendChild(form);
        form.submit();

        return false;
      })

    })
  </script>
</body>
{{end}}
