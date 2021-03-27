package textbooks

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
)

func Run() error {

	name := "textbooks"

	storage := utils.Storage(name)

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	doc := root.Clone()
	doc.AllowURLRevisit = false

	root.OnHTML("section.content", func(e *colly.HTMLElement) {
		e.ForEach(".entry .themeform .filetitle a[href]", func(_ int, el *colly.HTMLElement) {
			doc.Visit(el.Attr("href"))
		})
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

	root.Visit("https://proconsumidor.gob.do/monitoreo-de-libros-de-textos/")

	return nil
}
