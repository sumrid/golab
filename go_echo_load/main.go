package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sumrid/golab/go_echo/model"
	"github.com/tjarratt/babble"
)

func main() {
	t := time.Now()
	c := make(chan int, 500)

	for i := 0; i < cap(c); i++ {
		go GetLoad(c, i)
		// go PostLoad(c, i)
	}

	for i := 0; i < cap(c); i++ {
		fmt.Println(<-c) // รอให้ทุกตัวส่งค่ากลับมา
	}
	fmt.Println("time:", time.Since(t)) // จบการทำงาน พร้อมกับแสดงเวลาที่ใช้ไป
}

// GetLoad ทำการเรียกไปยัง endpoint GET:/book
// Ref: https://www.thepolyglotdeveloper.com/2017/07/consume-restful-api-endpoints-golang-application/
func GetLoad(c chan int, i int) {
	// Call endpoint
	res, err := http.Get("http://localhost/book")
	if err != nil {
		panic(err)
	}

	// อ่านข้อมูลที่ได้จากการเรียก endpoint
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	c <- i // ส่งค่าผ่าน channel
}

// PostLoad ทำการเรียกไปยัง endpoint POST:/book
func PostLoad(c chan int, i int) {
	// Create and generate data
	book := model.Book{}
	generateBook(&book)

	bookJSON, err := json.Marshal(book)
	if err != nil {
		panic(err)
	}

	// Request ไปยัง endpoint พร้อมกับรับ response
	res, err := http.Post("http://localhost/book", "application/json", bytes.NewBuffer(bookJSON))
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	c <- i // ส่งค่าผ่าน channel
}

func generateBook(b *model.Book) {
	bb := babble.NewBabbler()
	bb.Separator = " "
	b.Title = bb.Babble()
	bb.Count = 10
	b.Description = bb.Babble()
	bb.Count = 2
	b.Author = bb.Babble()
	b.Publisher = bb.Babble()
}
