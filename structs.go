package parser

// Таблица новостей
// Описание в отдульной таблице
type News struct {
  Id    int    `json:"id"`
  Title string `json:"title"`
  Img   string `json:"img"`
}

type NewsDecription struct {
  Id         int    `json:"id"`
  NewsId     string `json:"news_id"`
  Desciption string `json:"description"`
}

// Источники откуда будем брать
type SourceList struct {
  Id   int    `json:"id"`
  Name string `json:"name"`
  Href string `json:"href"`
}

type Rules struct {
  Id           int    `json:"id"`
  NewsAttrs    string `json:"news_attr"`      // Для какого аттрибута правило
  SourceListId int    `json:"source_list_id"` // Для какого урла будем разбирать
  Rule         string `json:"rule"`           // Правило
  GetAttr      string `json:"get_attr"`       // Собирать из тега или содержимого. Если тега то какого?
  IsMain       bool   `json:"is_main"`        // Является ли блок главным в который надо входить?
  IsUnique     bool   `json:"is_unique"`      // Является ли поле уникальным чтобы по нему отслеживать (Title)
}
