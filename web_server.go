package parser

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	e.GET("/api/news/", p.api_get_news)
	e.GET("/api/news_description/:id", p.api_get_news_description)
	e.GET("/api/rule/", p.api_get_rule)
	e.POST("/api/rule/:id", p.api_set_rule)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

}

func (p *Parser) get_index(c echo.Context) error {
	return c.File("src/github.com/writerim/parser/public/index.html")
}

func (p *Parser) get_js(c echo.Context) error {
	return c.File(fmt.Sprintf("src/github.com/writerim/parser/public/js/%s", c.Param("file")))
}

func (p *Parser) get_all_rules() []Rules {
	rules := []Rules{}

	sql := "select * from rules"
	res, err := p.query(sql)
	if err != nil {
		return rules
	}

	defer res.Close()

	for res.Next() {
		rule := Rules{}
		res.Scan(&rule.Id,
			&rule.Name,
			&rule.Link,
			&rule.MainPath,
			&rule.ImgPath,
			&rule.ImgAttr,
			&rule.TitlePath,
			&rule.HrefPath,
			&rule.DescPath)
		rules = append(rules, rule)
	}
	return rules
}

func (p *Parser) get_all_news(filter_title string) []News {
	news := []News{}

	sql := "select * from news"
	if filter_title != "" {
		sql = fmt.Sprintf("select * from news where title like '%%%s%%'", strings.Replace(filter_title, "'", "\\'", -1))
	}

	res, err := p.query(sql)
	if err != nil {
		return news
	}

	defer res.Close()

	for res.Next() {
		news_item := News{}
		res.Scan(&news_item.Id,
			&news_item.Title,
			&news_item.Img)
		news = append(news, news_item)
	}
	return news
}

func (p *Parser) api_get_rule(c echo.Context) error {
	return c.JSON(http.StatusCreated, p.get_all_rules())
}

func (p *Parser) api_get_news(c echo.Context) error {
	return c.JSON(http.StatusCreated, p.get_all_news(c.QueryParam("search")))
}

func (p *Parser) api_set_rule(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	method := c.FormValue("_method")

	if method == "DELETE" {

		if err != nil {
			return nil
		}
		del_sql := fmt.Sprintf("delete from rules where id = %d", id)
		res, err := p.query(del_sql)
		if err == nil {
			res.Close()
			return c.JSON(http.StatusCreated, Rules{})
		}
		return p.get_rule_by_id(id, c)
	}

	model := c.FormValue("model")

	rule := Rules{}

	if err := json.Unmarshal([]byte(model), &rule); err != nil {
		return nil
	}

	if rule.Id != 0 {

		up_sql := fmt.Sprintf(`
			update rules set 
			name='%s',
			link='%s',
			main_path='%s',
			img_path='%s',
			img_attr='%s',
			title_path='%s',
			href_path='%s',
			desc_path='%s'
			where id = %d
		`, strings.Replace(rule.Name, "'", "\\'", -1),
			strings.Replace(rule.Link, "'", "\\'", -1),
			strings.Replace(rule.MainPath, "'", "\\'", -1),
			strings.Replace(rule.ImgPath, "'", "\\'", -1),
			strings.Replace(rule.ImgAttr, "'", "\\'", -1),
			strings.Replace(rule.TitlePath, "'", "\\'", -1),
			strings.Replace(rule.HrefPath, "'", "\\'", -1),
			strings.Replace(rule.DescPath, "'", "\\'", -1),
			rule.Id)
		res, err := p.query(up_sql)
		if err != nil {
			return c.JSON(http.StatusCreated, Rules{})
		}
		defer res.Close()
		return p.get_rule_by_id(rule.Id, c)
	}

	insert_sql := fmt.Sprintf(`insert into rules
		(name, link, main_path, img_path, img_attr, title_path, href_path, desc_path) 
		values('%s','%s','%s','%s','%s','%s','%s','%s')`,
		strings.Replace(rule.Name, "'", "\\'", -1),
		strings.Replace(rule.Link, "'", "\\'", -1),
		strings.Replace(rule.MainPath, "'", "\\'", -1),
		strings.Replace(rule.ImgPath, "'", "\\'", -1),
		strings.Replace(rule.ImgAttr, "'", "\\'", -1),
		strings.Replace(rule.TitlePath, "'", "\\'", -1),
		strings.Replace(rule.HrefPath, "'", "\\'", -1),
		strings.Replace(rule.DescPath, "'", "\\'", -1))

	res, err := p.query(insert_sql)
	if err != nil {
		return nil
	}
	res.Close()

	res, err = p.query("SELECT LAST_INSERT_ID()")
	if err != nil {
		return nil
	}
	defer res.Close()
	res.Next()
	id_add := 0
	res.Scan(&id_add)

	return p.get_rule_by_id(id_add, c)
}

func (p *Parser) add_news(news News) News {
	insert_sql := fmt.Sprintf(`insert into news (title, img) values('%s','%s')`,
		strings.Replace(news.Title, "'", "\\'", -1),
		strings.Replace(news.Img, "'", "\\'", -1))

	res, err := p.query(insert_sql)
	if err != nil {
		return News{}
	}
	res.Close()

	res, err = p.query("SELECT LAST_INSERT_ID()")
	if err != nil {
		return News{}
	}
	res.Next()
	id_add := 0
	res.Scan(&id_add)
	res.Close()

	description := strings.Replace(news.Description, "'", "\\'", -1)
	sql := fmt.Sprintf("insert into news_description(news_id,description) values(%d,'%s')", id_add, description)

	res, err = p.query(sql)
	if err != nil {
		return News{}
	}
	defer res.Close()

	return p.news_by_id(id_add)
}

func (p *Parser) news_by_id(id int) News {
	sql := fmt.Sprintf("select * from news where id = %d", id)
	res, err := p.query(sql)
	if err != nil {
		return News{}
	}
	defer res.Close()

	news := News{}
	res.Next()
	res.Scan(&news.Id, &news.Title,
		&news.Img)
	defer res.Close()
	return news
}

func (p *Parser) get_rule_by_id(id int, c echo.Context) error {
	sql := fmt.Sprintf("select * from rule where id = %s", id)
	res, err := p.query(sql)
	if err != nil {
		return nil
	}
	defer res.Close()

	rule := Rules{}
	res.Next()
	res.Scan(&rule.Id, &rule.Name,
		&rule.Link,
		&rule.MainPath,
		&rule.ImgPath,
		&rule.ImgAttr,
		&rule.TitlePath,
		&rule.HrefPath,
		&rule.DescPath)

	return c.JSON(http.StatusCreated, rule)

}

func (p *Parser) api_get_news_description(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	sql := fmt.Sprintf("select * from news_description where news_id = %d", id)
	res, err := p.query(sql)
	if err != nil {
		return nil
	}
	defer res.Close()

	desc := NewsDecription{}
	res.Next()
	res.Scan(&desc.Id, &desc.NewsId,
		&desc.Desciption)

	return c.JSON(http.StatusCreated, desc)

}
