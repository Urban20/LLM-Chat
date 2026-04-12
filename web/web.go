package web

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Buscar(url string) (string, error) {

	resp, resperr := http.Get(url)

	if resperr != nil {

		return "", resperr

	}

	if resp.StatusCode != http.StatusOK {

		return "", fmt.Errorf("codigo de estado : %d", resp.StatusCode)

	}

	doc, docerr := goquery.NewDocumentFromReader(resp.Body)

	if docerr != nil {
		return "", docerr
	}

	if err := resp.Body.Close(); err != nil {
		return "", err
	}

	return doc.Find("body").Text(), nil
}
