{{define "container"}}
<head>
  <title>NarouEpub</title>
</head>
<body>
  <h1><a href="/">NarouEpub</a></h1>
  <h2>{{.Container.Title}}</h2>

  <dl>
    <dt>NCode</dt>
    <dd><a href="https://ncode.syosetu.com/{{.Container.NCode}}/" target="_blank">{{.Container.NCode}}</a></dd>
    <dt>Author</dt>
    <dd><a href="https://mypage.syosetu.com/{{.Container.UserID}}/" target="_blank">{{.Container.Author}}</a></dd>
    <dt>GeneralLastUp</dt>
    <dd>{{.Container.GeneralLastUp}}</dd>
    <dt>UpdatedAt</dt>
    <dd>{{.Container.UpdatedAt}}</dd>
    <dt>GeneralAllNo</dt>
    <dd>{{.Container.GeneralAllNo}}</dd>
  </dl>

  <p>
    <span>
      <a
        class="post-form"
        data-method="post"
        href="/containers/{{.Container.NCode}}/fetch"
        onclick="return false;"
      >再読みこみ
      </a></span>
    <span>
      <a
        class="post-form"
        data-method="post"
        href="/containers/{{.Container.NCode}}/build"
        onclick="return false;"
      >epub生成
      </a></span>
    <span>
      <a
        class="post-form"
        data-method="post"
        href="/containers/{{.Container.NCode}}/publish"
        onclick="return false;"
      >publish</a></span>
  </p>

  <table>
    <thead>
      <tr>
        <th>EpisodeNumber</th>
        <th>EpisodeTitle</th>
        <th>CrawledAt</th>
        <th></th>
      </tr>
    </thead>
    <tbody>
    {{range $i, $v := .Container.Episodes}}
      <tr>
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
