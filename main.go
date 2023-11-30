package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Name     string `json:"name"`
	Price    string `json:"price"`
	ImageUrl string `json:"imageUrl"`
}

func main() {
	// Permitir domínios para realizar o Scrap
	call := colly.NewCollector(
		colly.AllowedDomains("j2store.net"),
	)

	// Array de Itens
	var items []item

	// Busca dos dados pelos elementos
	call.OnHTML("div.col-sm-9 div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
		item := item{
			Name:     h.ChildText("h2.product-title"),
			Price:    h.ChildText("div.sale-price"),
			ImageUrl: h.ChildAttr("img", "src"),
		}
		
		items = append(items, item)
	})

	// Buscar a próxima página
	call.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		call.Visit(next_page)
	})

	call.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	// Fetch na URL
	call.Visit("http://j2store.net/demo/index.php/shop")

	// Transformar os dados do Web Scraping em JSON
	content, err := json.Marshal(items)

	if err != nil {
		log.Fatal(err)
	}

	// Vai gerar os produtos encontrados dentro um arquivo products.json internamente
	os.WriteFile("products.json", content, 0644)
}
