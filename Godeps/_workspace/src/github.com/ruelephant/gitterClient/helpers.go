package gitterClient
import (
	"net/http"
	"log"
	"io"
)

func getRequest(url string) (io.ReadCloser) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("%s", err)
	}
	return response.Body
}
