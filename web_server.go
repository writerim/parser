package parser

import (
  "fmt"
  "github.com/labstack/echo"
  "html/template"
  "io"
  "net/http"
  "strconv"
)

type TemplateRenderer struct {
  templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, nil)
}

func (p *Parser) StartWebServer(port int) {

  e := echo.New()

  e.GET("/", p.get_index)
  e.GET("/js/:file", p.get_js)

  // Получение списка новостей
  e.GET("/api/news/:offset/:limit/", p.api_get_news_list)

  // Получение списка значений по свойствам
  // e.GET("/api/news_attrs_value/:offset/:limit/", p.api_get_news_attrs_value_list)

  // Получение списка по заголовку
  e.GET("/api/news/find/:str_find/", p.api_get_finden_news)

  // Получение все ресурсов для парсинга
  e.GET("/api/source/:offset/:limit/", p.api_get_source_list)

  // // Получение правил для ресурса
  e.GET("/api/news_rule_list/", p.api_get_rules)

  e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

}

func (p *Parser) get_index(c echo.Context) error {
  return c.File("src/github.com/parser/public/index.html")
}

func (p *Parser) get_js(c echo.Context) error {
  return c.File(fmt.Sprintf("src/github.com/parser/public/js/%s", c.Param("file")))
}

func (p *Parser) api_get_news_list(c echo.Context) error {

  offset, err := strconv.Atoi(c.Param("offset"))
  if err != nil {
    offset = 0
  }
  limit, err := strconv.Atoi(c.Param("limit"))
  if err != nil {
    limit = 15
  }

  sql := fmt.Sprintf("select * from news order by id desc limit %d,%d", offset, limit)

  res, err := p.query(sql)
  if err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "[]")
  }

  news := []News{}

  defer res.Close()

  for res.Next() {
    news_item := News{}
    err = res.Scan(&news_item.Id, &news_item.Title, &news_item.Img)
    if err != nil {
      return c.JSON(http.StatusOK, []News{})
    }
    news = append(news, news_item)
  }

  return c.JSON(http.StatusOK, news)
}

func (p *Parser) api_get_finden_news(c echo.Context) error {

  sql := fmt.Sprintf("select * from news where title like '%%%s%%'", c.Param("str_find"))

  res, err := p.query(sql)
  if err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "[]")
  }

  news := []News{}

  defer res.Close()

  for res.Next() {
    news_item := News{}
    err = res.Scan(&news_item.Id, &news_item.Title, &news_item.Img)
    if err != nil {
      return c.JSON(http.StatusOK, []News{})
    }
    news = append(news, news_item)
  }

  return c.JSON(http.StatusOK, news)
}

func (p *Parser) api_get_source_list(c echo.Context) error {

  offset, err := strconv.Atoi(c.Param("offset"))
  if err != nil {
    offset = 0
  }
  limit, err := strconv.Atoi(c.Param("limit"))
  if err != nil {
    limit = 15
  }

  sql := fmt.Sprintf("select * from source_list limit %d,%d", offset, limit)

  res, err := p.query(sql)
  if err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "[]")
  }

  sources := []SourceList{}

  defer res.Close()

  for res.Next() {
    source := SourceList{}
    err = res.Scan(&source.Id, &source.Name, &source.Href)
    if err != nil {
      return c.JSON(http.StatusOK, []SourceList{})
    }
    sources = append(sources, source)
  }

  return c.JSON(http.StatusOK, sources)
}

func (p *Parser) api_get_rules(c echo.Context) error {

  res, err := p.query("select * from attrs_rule_list")
  if err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "[]")
  }

  attrs := []Rules{}

  defer res.Close()

  for res.Next() {
    attr := Rules{}
    err = res.Scan(&attr.Id, &attr.NewsAttrs, &attr.SourceListId, &attr.Rule, &attr.GetAttr, &attr.IsMain, &attr.IsUnique)
    if err != nil {
      return c.JSON(http.StatusOK, []Rules{})
    }
    attrs = append(attrs, attr)
  }

  return c.JSON(http.StatusOK, attrs)
}
