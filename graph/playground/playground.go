package playground

import (
	"html/template"
	"net/http"
)

var apage = template.Must(template.ParseFiles("graph/playground/index.html"))

func Handler(title string, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := apage.Execute(w, map[string]string{
			"title":      title,
			"endpoint":   endpoint,
			"version":    "1.7.20",
			"cssSRI":     "sha256-cS9Vc2OBt9eUf4sykRWukeFYaInL29+myBmFDSa7F/U=", // I don't know what these do...
			"faviconSRI": "sha256-GhTyE+McTU79R4+pRO6ih+4TfsTOrpPwD8ReKFzb3PM=", // I don't know what these do...
			"jsSRI":      "sha256-4QG1Uza2GgGdlBL3RCBCGtGeZB6bDbsw8OltCMGeJsA=", // I don't know what these do...
		})
		if err != nil {
			panic(err)
		}
	}
}
