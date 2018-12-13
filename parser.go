package parser

import (
  "database/sql"
  "fmt"
  "github.com/antchfx/xquery/html"
  _ "github.com/go-sql-driver/mysql"
  "golang.org/x/net/html"
  "os"
  "os/signal"
  "time"
)

type Parser struct {
  db_conn *sql.DB
}

func New() *Parser {
  return &Parser{}
}

func (p *Parser) ConnectDb(ip string, port int, login, password, database string) error {
  database_connect := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
    login,
    password,
    ip,
    port,
    database)

  db, err := sql.Open("mysql", database_connect)

  if err != nil {
    return err
  }
  p.db_conn = db
  return nil
}

/*
  "Демон" постоянно повторяющий парсер страниц
*/
func (p *Parser) StartDeamon() {

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func() {
    <-c
    fmt.Println("[SIGINT force quit]")
    p.db_conn.Close()
    os.Exit(0)
  }()

  for {
    p.Parse()
    fmt.Println("stop")
    time.Sleep(1 * time.Minute)

  }

}

/*

  Парсинг по данным пользователя

*/
func (p *Parser) Parse() {
  // Собираем правила сбора
  source_list, err := p.get_source_list()
  if err != nil {
    return // По какой то причине отвалилось. Лучше все остановить
  }
  attrs_rule_list, err := p.get_rules()
  if err != nil {
    return // По какой то причине отвалилось. Лучше все остановить
  }

  for i := 0; i < len(source_list); i++ {
    l := source_list[i]

    doc, err := htmlquery.LoadURL(l.Href)
    if err != nil {
      continue
    }

    // Разбираем Отделаем главную ноду от сотальных
    main_block, nodes := p.find_main_rule_block(attrs_rule_list, l.Id)

    htmlquery.FindEach(doc, main_block.Rule, func(_ int, node *html.Node) {

      news_item := News{}

      for a := 0; a < len(nodes); a++ {
        r := nodes[a]

        htmlquery.FindEach(node, r.Rule, func(i int, node *html.Node) {

          value := ""
          if r.GetAttr != "" {
            value = htmlquery.SelectAttr(node, r.GetAttr)
          } else {
            value = htmlquery.InnerText(node)
          }

          switch r.NewsAttrs {
          case "img":
            news.Img = value
          case "title":
            news.Title = value
          case "href":
            news.Href = value
          }

          return

        })

      }

      if !p.find_double_new(news_item.Title) {
        sql_append := fmt.Sprintf("insert into (title , img) values('%s','%s')", news_item.Title, news_item.Img)
      }

    })

  }
}

/*
  Получение всех ресурсов для сборки
*/
func (p *Parser) get_source_list() ([]SourceList, error) {

  source_list := []SourceList{}

  results, err := p.query("SELECT id, name , href FROM source_list")
  if err != nil {
    return source_list, err
  }
  defer results.Close()

  for results.Next() {
    source := SourceList{}
    err = results.Scan(&source.Id, &source.Name, &source.Href)
    if err != nil {
      return source_list, err
    }
    source_list = append(source_list, source)
  }

  return source_list, nil
}

/*
  Получение всех правил сборки
*/
func (p *Parser) get_rules() ([]Rules, error) {

  attrs_rule_list := []Rules{}

  results, err := p.query("SELECT id, news_attr, source_list_id, rule, get_attr , is_main , is_unique FROM rules")
  if err != nil {
    return attrs_rule_list, err
  }

  defer results.Close()

  for results.Next() {
    attrs := Rules{}
    err = results.Scan(&attrs.Id, &attrs.NewsAttrs, &attrs.SourceListId, &attrs.Rule, &attrs.GetAttr, &attrs.IsMain, &attrs.IsUnique)
    if err != nil {
      return attrs_rule_list, err
    }

    attrs_rule_list = append(attrs_rule_list, attrs)
  }
  return attrs_rule_list, nil
}

/*
  Поиск главного блока среди всех правил
*/
func (p *Parser) find_main_rule_block(attrs_rule_list []Rules, source_id int) (Rules, []Rules) {
  nodes := []Rules{}
  main_block := Rules{}

  for j := 0; j < len(attrs_rule_list); j++ {
    if attrs_rule_list[j].SourceListId == source_id {
      if attrs_rule_list[j].IsMain {
        main_block = attrs_rule_list[j]
      } else {
        nodes = append(nodes, attrs_rule_list[j])
      }
    }
  }
  return main_block, nodes
}

/*
  Поиск дубликата в БД
*/
func (p *Parser) find_double_new(title string) bool {

  count := 0

  if title == "" {
    return true
  }

  sql_find_double := fmt.Sprintf("SELECT count(*) FROM news where title = '%s'", title)
  results, err := p.query(sql_find_double)
  if err != nil {
    return true
  }

  defer results.Close()

  results.Next()

  err = results.Scan(&count)
  if err != nil {
    return true
  }

  return count > 0
}

func (p *Parser) query(sql string) (*sql.Rows, error) {
  err := p.db_conn.Ping()
  if err != nil {
    return nil, err
  }
  return p.db_conn.Query(sql)
}
