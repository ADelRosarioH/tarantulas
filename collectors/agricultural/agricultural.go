package agricultural

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
)

func Run() error {

	name := "agricultural"

	storage := utils.Storage(name)

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	page := root.Clone()
	page.AllowURLRevisit = false

	doc := page.Clone()

	root.OnHTML("#menu-item-7437 ul li a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Root", e.Attr("href"))
		page.Visit(e.Attr("href"))
	})

	page.OnHTML(".filetitle a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Doc", e.Attr("href"))
		doc.Visit(e.Attr("href"))
	})

	doc.OnResponse(func(r *colly.Response) {
		fmt.Println("Downloading", r.FileName())

		tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*_%s", r.FileName()))

		if err != nil {
			panic(err)
		}

		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		r.Save(tempFile.Name())

		if err := utils.UploadToS3(name, tempFile); err != nil {
			panic(err)
		}
	})

	root.Visit("https://proconsumidor.gob.do/")

	return nil
}
