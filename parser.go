package main

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

// В таблице все новости
// Новость не может быть без названия.
// Остальные поля вынесем для масштабирования аттрибутов
type News struct {
  Id       int `json:"id"`
  SourceId int `json:"source_id"`
}

// Дополнительные аттрибуты новостей
// Ссылки, атвор, дата выпуска, канртинка  и прочее
type NewsAttrsValue struct {
  Id     int    `json:"id"`
  NewsId int    `json:"news_id"`
  AttrId int    `json:"news_attrs_id"`
  Value  string `json:"value"`
}

// Описание всех дополнительных полей новости
type NewsAttrs struct {
  Id   int    `json:"id"`
  Name string `json:"name"`
}

// Источники откуда будем брать
type SourceList struct {
  Id   int    `json:"id"`
  Name string `json:"name"`
  Href string `json:"href"`
}

type AttrsRulesList struct {
  Id           int    `json:"id"`
  NewsAttrsId  int    `json:"news_attr_id"`   // Для какого аттрибута правило
  SourceListId int    `json:"source_list_id"` // Для какого урла будем разбирать
  Rule         string `json:"rule"`           // Правило
  GetAttr      string `json:"get_attr"`       // Собирать из тега или содержимого. Если тега то какого?
  IsMain       bool   `json:"is_main"`        // Является ли блок главным в который надо входить?
  IsUnique     bool   `json:"is_unique"`      // Является ли поле уникальным чтобы по нему отслеживать (Title)
}

func main() {

  db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/parser")

  if err != nil {
    panic(err.Error())
  }

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func() {
    <-c
    fmt.Println("[SIGINT force quit]")
    db.Close()
    os.Exit(0)
  }()

  for {

    // Собираем правила сбора
    source_list, err := get_source_list(db)
    if err != nil {
      break // По какой то причине отвалилось. Лучше все остановить
    }
    attrs_rule_list, err := get_attrs_rule_list(db)
    if err != nil {
      break // По какой то причине отвалилось. Лучше все остановить
    }

    for i := 0; i < len(source_list); i++ {
      l := source_list[i]

      doc, err := htmlquery.LoadURL(l.Href)
      if err != nil {
        continue
      }

      // Разбираем Отделаем главную ноду от сотальных
      main_block, nodes := find_main_rule_block(attrs_rule_list, l.Id)

      htmlquery.FindEach(doc, main_block.Rule, func(_ int, node *html.Node) {

        title := ""
        sql_buff := []string{}

        for a := 0; a < len(nodes); a++ {
          r := nodes[a]

          htmlquery.FindEach(node, r.Rule, func(i int, node *html.Node) {

            value := ""
            if r.GetAttr != "" {
              value = htmlquery.SelectAttr(node, r.GetAttr)
            } else {
              value = htmlquery.InnerText(node)
            }

            if r.IsUnique {
              title = value
              sql_buff = append(sql_buff, fmt.Sprintf("INSERT INTO news (source_id,title) values(%d,'%s')", l.Id, title))
            } else {
              sql_add_attr := fmt.Sprintf("INSERT INTO news_attrs_value (news_id,news_attr_id , value) values( LAST_INSERT_ID() , %d , '%s')", r.NewsAttrsId, value)
              sql_buff = append(sql_buff, sql_add_attr)
            }

            return

          })

        }

        // Нашли дубликат. Выходим и з парсинга сущности
        if find_double_new(db, title) {
          return
        }

        // Перебираем все запросы на добавление
        for c := 0; c < len(sql_buff); c++ {
          results, err := db.Query(sql_buff[c])
          if err != nil {
            return
          }
          defer results.Close()
        }

      })

    }

    fmt.Println("stop")
    time.Sleep(1 * time.Minute)

  }

}

/*
  Получение всех ресурсов для сборки
*/
func get_source_list(db *sql.DB) ([]SourceList, error) {

  source_list := []SourceList{}

  results, err := db.Query("SELECT id, name , href FROM source_list")
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
func get_attrs_rule_list(db *sql.DB) ([]AttrsRulesList, error) {

  attrs_rule_list := []AttrsRulesList{}

  results, err := db.Query("SELECT id, news_attr_id, source_list_id, rule, get_attr , is_main , is_unique FROM attrs_rule_list")
  if err != nil {
    return attrs_rule_list, err
  }

  defer results.Close()

  for results.Next() {
    attrs := AttrsRulesList{}
    err = results.Scan(&attrs.Id, &attrs.NewsAttrsId, &attrs.SourceListId, &attrs.Rule, &attrs.GetAttr, &attrs.IsMain, &attrs.IsUnique)
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
func find_main_rule_block(attrs_rule_list []AttrsRulesList, source_id int) (AttrsRulesList, []AttrsRulesList) {
  nodes := []AttrsRulesList{}
  main_block := AttrsRulesList{}

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
func find_double_new(db *sql.DB, title string) bool {

  count := 0

  if title == "" {
    return true
  }

  sql_find_double := fmt.Sprintf("SELECT count(*) FROM news where title = '%s'", title)
  results, err := db.Query(sql_find_double)
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
