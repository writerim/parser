package parser

// Таблица новостей
// Описание в отдульной таблице
type News struct {
	Id          int    `json:"id"`
	Title       string `json:"title"` // Название новости
	Img         string `json:"img"`   // Путь до картинки
	Href        string
	Description string
}

// Сдержания новостей
type NewsDecription struct {
	Id         int    `json:"id"`
	NewsId     string `json:"news_id"`     // К какой новости
	Desciption string `json:"description"` // Сожаржание
}

// Правила разбора
type Rules struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`       // Название правила
	Link      string `json:"link"`       // Ссылка откуда брать информацию
	MainPath  string `json:"main_path"`  // Путь до блока с новостью
	ImgPath   string `json:"img_path"`   // Путь до картинки
	ImgAttr   string `json:"img_attr"`   // Из какого атрибута брать картинку
	TitlePath string `json:"title_path"` // Путь до названия новости
	HrefPath  string `json:"href_path"`  // Путь до ссылки новости. Если пусто то на той же странице
	DescPath  string `json:"desc_path"`  // Путь до содержиого новости
}

type BackboneRequest struct {
	Model  string `json:"model"`
	Method string `json:"_method"`
}
