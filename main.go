package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", "127.0.0.1:8000", "address to listen to")
var statusURL = flag.String("url", "https://www.gran-turismo.com/us/api/gt7/server/status", "GT7 Status page")

func main() {
	flag.Parse()
	http.HandleFunc("/text", textHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/pretty", prettyHandler)
	http.HandleFunc("/", prettyHandler)
	fmt.Printf("Starting server on http://%s\n", *addr)
	log.Panicln(http.ListenAndServe(*addr, nil))
}

func getStatus() (string, error) {
	res, err := http.Post(*statusURL, "text/html", nil)
	if err != nil {
		log.Println(err)
		return "error", err
	}
	if res.StatusCode < 500 {
		return "Online", nil
	}
	return "Offline", nil
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	status, err := getStatus()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, status)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	status, err := getStatus()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}
	out, err := json.Marshal(struct{ Status string }{status})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}
	fmt.Fprintf(w, "%s", out)
}

func prettyHandler(w http.ResponseWriter, r *http.Request) {
	status, _ := getStatus()
	err := tpl.Execute(w, struct {
		Status string
	}{
		status,
	})
	if err != nil {
		log.Println(err)
	}
}

var tpl = template.Must(template.New("tpl").Parse(`
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
    <title>GT7 Status: {{.Status}}</title>
  </head>
  <body>
    <h1>GT7 Status: {{.Status}}</h1>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
  </body>
</html>`))
