package medicines

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

func Run() error {

	name := "medicines"

	storage := utils.Storage(name)

	root := colly.NewCollector(colly.AllowURLRevisit())

	root.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	extensions.RandomUserAgent(root)

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	page := root.Clone()

	doc := page.Clone()
	doc.AllowURLRevisit = false

	root.OnHTML("#menu-item-55 ul li a[href]", func(e *colly.HTMLElement) {
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

		s3FileKey, err := utils.UploadToS3(name, tempFile)

		if err != nil {
			panic(err)
		}

		textMessage := fmt.Sprintf("%s collector uploaded %s to S3", name, s3FileKey)

		if err := utils.Notify(textMessage); err != nil {
			log.Fatal(err)
		}
	})

	err := root.Visit("https://proconsumidor.gob.do/")

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
