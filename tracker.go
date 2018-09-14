package main

// hexdump -e '16/1 "0x%02x, " "\n"' pixel.png

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type Data struct {
	Title string
	Dest  string
	Check time.Time
}

var IDS = map[string]*Data{}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/track/", trackHandler)
	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/", formHandler)
	fmt.Println("Listening on 7777")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.EscapedPath()
	aurl := strings.Split(url, "/")
	id := aurl[len(aurl)-1]
	id = strings.Replace(id, ".png", "", 1)
	IDS[id].Check = time.Now()
	w.Header().Set("Cache-control", "private, max-age=0, no-cache")
	w.Header().Set("Content-Type", "image/png")
	sendEmail("kidandcat@gmail.com", "jairo@galax.be", IDS[id].Title+" to "+IDS[id].Dest+" --- "+IDS[id].Check.Format("_2 Monday January 15:04:05 2006"))
	// 1x1 PNG
	w.Write([]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x06, 0x62, 0x4b, 0x47, 0x44, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0xa0,
		0xbd, 0xa7, 0x93, 0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59, 0x73, 0x00, 0x00, 0x0b, 0x13, 0x00,
		0x00, 0x0b, 0x13, 0x01, 0x00, 0x9a, 0x9c, 0x18, 0x00, 0x00, 0x00, 0x07, 0x74, 0x49, 0x4d, 0x45,
		0x07, 0xe2, 0x09, 0x0e, 0x0b, 0x35, 0x1b, 0x5e, 0x96, 0xe1, 0xd0, 0x00, 0x00, 0x00, 0x1d, 0x69,
		0x54, 0x58, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x43,
		0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x47, 0x49, 0x4d, 0x50,
		0x64, 0x2e, 0x65, 0x07, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8,
		0xff, 0xff, 0x3f, 0x03, 0x00, 0x08, 0xfc, 0x02, 0xfe, 0x5c, 0x9f, 0xcf, 0xda, 0x00, 0x00, 0x00,
		0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	res := ""
	w.Header().Set("Content-Type", "text/plain")
	for _, v := range IDS {
		res += v.Title + " to " + v.Dest + " - - - " + v.Check.Format("_2 Monday January 15:04:05 2006")
	}
	w.Write([]byte(res))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    <form action="/new">
        <p>
            <input name="title" type="text" placeholder="Titulo">
        </p>
        <p>
            <input name="dest" type="text" placeholder="Destinatario">
        </p>
        <input type="submit">
    </form>
</body>
</html>
	`))
}

func sendEmail(to, from, body string) {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	c.Mail(from)
	c.Rcpt(to)
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(body)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := randStringRunes(15)

	IDS[id] = &Data{
		Title: r.Form["title"][0],
		Dest:  r.Form["dest"][0],
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("https://track.galax.be/track/" + id + ".png"))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
