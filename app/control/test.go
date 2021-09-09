package control

import (
	"chat/library"
	"fmt"
	"html/template"
	"net/http"
)

func TcpTest(w http.ResponseWriter, r *http.Request) {
}

func HttpTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, err := template.ParseFiles("G:/WWW/golang/src/chat/html/test.html")
		if err != nil {
			fmt.Fprintf(w, "parse template error: %s", err.Error())
			return
		}
		t.Execute(w, nil)
	} else {
		username := r.Form["username"]
		password := r.Form["password"]
		fmt.Fprintf(w, "username = %s, password = %s", username, password)
	}
}

func WsTest(ws library.WsConn, data []byte) (result []byte) {
	fmt.Println("发送了：" + string(data))
	return []byte("发送了：" + string(data))
}
