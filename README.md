# parser

```
Нужно написать агрегатор новостей.
Пользователь подает на вход адрес новостного сайта или его RSS-ленты и правило парсинга (формат правила — на усмотрение разработчика).
База данных агрегатора начинает автоматически пополняться новостями с этого сайта.
У пользователя есть возможность просматривать список новостей из базы и искать их по подстроке в заголовке новости.
В качестве примера требуется подключить два любых новостных сайта на выбор.
Результат — исходный код агрегатора, а также рабочие адреса и правила парсинга, которые можно подать ему на вход.  
Язык —Golang. Хранилище — любая реляционная база данных.
```


## Example
```golang
package main

import (
  "github.com/parser"
)

func main() {

  p := parser.New()
  err := p.ConnectDb("127.0.0.1", 3306, "root", "", "parser")
  if err != nil {
    return
  }

  go p.StartWebServer(8080)
  p.StartDeamon()

}


```


Писать правила надо по стандарту xPath  
Пример : 
```
//*[@class='class_name']  
//*[contains(@class,'Link-root')]@href  
//*[contains(@class,'Link-root')]  
```