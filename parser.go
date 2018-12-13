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

func (p *Parser) Parse() {
	// Соберем все правла
	rules := p.get_all_rules()

	news := []News{}

	for i := 0; i < len(rules); i++ {

		rule := rules[i]

		doc, err := htmlquery.LoadURL(rule.Link)
		if err != nil {
			continue
		}

		// Пробуем защитится от невенрных тегов
		last_title := ""

		htmlquery.FindEach(doc, rule.MainPath, func(_ int, node *html.Node) {

			news_item := News{}
			htmlquery.FindEach(node, rule.TitlePath, func(i int, node_title *html.Node) {
				value := htmlquery.InnerText(node_title)
				if last_title == value {
					return
				}
				last_title = value
				news_item.Title = last_title

				if news_item.Title != "" {

					htmlquery.FindEach(node, rule.ImgPath, func(i int, node_img *html.Node) {

						get_attr := "src"
						if rule.ImgAttr != "" {
							get_attr = rule.ImgAttr
						}

						value := htmlquery.SelectAttr(node_img, get_attr)
						news_item.Img = value
					})

					if rule.HrefPath != "" {
						htmlquery.FindEach(node, rule.HrefPath, func(i int, node_href *html.Node) {
							news_item.Href = htmlquery.SelectAttr(node_href, "href")
						})

						doc_desc, err := htmlquery.LoadURL(rule.Link + rule.HrefPath)
						if err == nil {
							htmlquery.FindEach(doc_desc, rule.DescPath, func(i int, node_desc *html.Node) {
								value := htmlquery.InnerText(node_desc)
								fmt.Println(value)
							})
						}
					}
					if !p.find_double_new(news_item.Title) {
						news = append(news, news_item)
					}

				}
			})

		})

	}

	for i := 0; i < len(news); i++ {
		p.add_news(news[i])
	}

	fmt.Println(news)

}

func (p *Parser) query(sql string) (*sql.Rows, error) {
	err := p.db_conn.Ping()
	if err != nil {
		return nil, err
	}
	return p.db_conn.Query(sql)
}

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
