package apify

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/matryer/is"
)

func Test_Params(t *testing.T) {
	assert := is.New(t)
	router := mux.NewRouter()

	router.HandleFunc("/kitchens/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			UUID string `params:"uuid"`
		}

		if err := Params(r, &params); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		assert.Equal(params.UUID, "ad6e7682-19b4-4295-add4-a409687d41ca")

		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			ID int `params:"id"`
		}

		if err := Params(r, &params); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		assert.Equal(params.ID, 100500)

		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/kitchens", func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Limit  int64 `params:"limit,query"`
			Offset int64 `params:"offset,query"`
		}

		params.Limit = 500

		if err := Params(r, &params); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		assert.Equal(params.Limit, int64(500))
		assert.Equal(params.Offset, int64(100))

		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(router)

	defer func() {
		server.Close()
	}()

	urls := []string{
		fmt.Sprintf("%s/kitchens/ad6e7682-19b4-4295-add4-a409687d41ca", server.URL),
		fmt.Sprintf("%s/kitchens?offset=100", server.URL),
		fmt.Sprintf("%s/users/%d", server.URL, 100500),
	}

	for _, url := range urls {
		res, err := http.Get(url)

		assert.NoErr(err)

		defer func() {
			res.Body.Close()
		}()
	}
}
