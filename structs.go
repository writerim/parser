package parser

// В таблице все новости
// Новость не может быть без названия.
// Остальные поля вынесем для масштабирования аттрибутов
type News struct {
  Id       int    `json:"id"`
  SourceId int    `json:"source_id"`
  Title    string `json:"title"`
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
  Id    int    `json:"id"`
  Name  string `json:"name"`
  Ident string `json:"ident"`
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
