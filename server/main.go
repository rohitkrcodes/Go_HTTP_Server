package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	mu sync.Mutex
	db map[string]dollars
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db.db {
		fmt.Fprintf(w, "%s %s\n", item, price)
	}
}

func (db *database) add(w http.ResponseWriter, req *http.Request) {

	db.mu.Lock()
	defer db.mu.Unlock()

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db.db[item]; ok {
		msg := fmt.Sprintf("Dupllicate item %q", item)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	p, err := strconv.ParseFloat(price, 32)

	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	db.db[item] = dollars(p)

	fmt.Fprintf(w, "added %s with price of %s\n", item, db.db[item])
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	db.mu.Lock()
	defer db.mu.Unlock()

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db.db[item]; !ok {
		msg := fmt.Sprintf("Item %q does not exist", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	p, err := strconv.ParseFloat(price, 32)

	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) // 400
		return
	}

	db.db[item] = dollars(p)

	fmt.Fprintf(w, "new price for %s is %s\n", item, db.db[item])
}

func (db *database) fetch(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db.db[item]; !ok {
		msg := fmt.Sprintf("Item %q does not exist", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	fmt.Fprintf(w, "item %s has price %s\n", item, db.db[item])
}

func (db *database) drop(w http.ResponseWriter, req *http.Request) {

	db.mu.Lock()
	defer db.mu.Unlock()

	item := req.URL.Query().Get("item")

	if _, ok := db.db[item]; !ok {
		msg := fmt.Sprintf("Item %q does not exist", item)
		http.Error(w, msg, http.StatusNotFound) // 404
		return
	}

	delete(db.db, item)

	fmt.Fprintf(w, "item %s has been deleted successfully\n", item)
}

func main() {
	db := database{
		db: map[string]dollars{
			"shoes": 50,
			"socks": 5,
		},
	}

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/read", db.fetch)
	http.HandleFunc("/delete", db.drop)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
