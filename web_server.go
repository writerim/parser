package parser

import (
  "fmt"
  "github.com/labstack/echo"
  "net/http"
  "strconv"
)

func (p *Parser) StartWebServer(port int) {

  e := echo.New()

  // Получение списка новостей
  e.GET("/api/news/:offset/:limit/", p.api_get_news_list)

  // Получение списка по заголовку
  e.GET("/api/news/find/:str_find/", p.api_get_finden_news)

  // Получение все ресурсов для парсинга
  e.GET("/api/source/:offset/:limit/", p.api_get_source_list)

  // // Получение правил для ресурса
  e.GET("/api/source/:id/rules/", p.api_get_source_rules)

  e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

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

  sql := fmt.Sprintf("select * from news limit %d,%d", offset, limit)

  res, err := p.query(sql)
  if err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "[]")
  }

  news := []News{}

  defer res.Close()

  for res.Next() {
    news_item := News{}
    err = res.Scan(&news_item.Id, &news_item.SourceId, &news_item.Title)
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
    err = res.Scan(&news_item.Id, &news_item.SourceId, &news_item.Title)
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

func (p *Parser) api_get_source_rules(c echo.Context) error {

  id, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    id = 0
  }

  sql := fmt.Sprintf("select * from attrs_rule_list where source_list_id = %d", id)

  fmt.Println(sql)

  res, err := p.query(sql)
  if err != nil {
    return echo.NewHTTPError(http.StatusNotFound, "[]")
  }

  attrs := []AttrsRulesList{}

  defer res.Close()

  for res.Next() {
    attr := AttrsRulesList{}
    err = res.Scan(&attr.Id, &attr.NewsAttrsId, &attr.SourceListId, &attr.Rule, &attr.GetAttr, &attr.IsMain, &attr.IsUnique)
    if err != nil {
      return c.JSON(http.StatusOK, []AttrsRulesList{})
    }
    attrs = append(attrs, attr)
  }

  return c.JSON(http.StatusOK, attrs)
}
