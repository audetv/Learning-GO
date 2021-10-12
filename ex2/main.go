package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Pizza struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Pizzas []Pizza

func (ps Pizzas) findById(id int) (Pizza, error) {
	for _, pizza := range ps {
		if pizza.Id == id {
			return pizza, nil
		}
	}
	return Pizza{}, fmt.Errorf("couldn't find pizza with Id: %d", id)
}

type Order struct {
	PizzaID  int `json:"pizza_id"`
	Quantity int `json:"quantity"`
	Total    int `json:"total"`
}

type Orders []Order

type pizzasHandler struct {
	pizzas *Pizzas
}

func (ph pizzasHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if len(*ph.pizzas) == 0 {
			http.Error(w, "Error: no pizzas found", http.StatusNotFound)
			return
		}
		err := json.NewEncoder(w).Encode(ph.pizzas)
		if err != nil {
			log.Fatal(err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type ordersHandler struct {
	pizzas *Pizzas
	orders *Orders
}

func (oh ordersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var o Order

		if len(*oh.pizzas) == 0 {
			http.Error(w, "Error: no pizzas found", http.StatusNotFound)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&o)
		if err != nil {
			http.Error(w, "Can't decode body", http.StatusBadRequest)
			return
		}

		p, err := oh.pizzas.findById(o.PizzaID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error %s:", err), http.StatusBadRequest)
		}

		o.Total = p.Price * o.Quantity
		*oh.orders = append(*oh.orders, o)
		err = json.NewEncoder(w).Encode(o)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodGet:
		err := json.NewEncoder(w).Encode(oh.orders)
		if err != nil {
			log.Fatal(err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	var orders Orders
	pizzas := Pizzas{
		Pizza{
			Id:    1,
			Name:  "Pepperoni",
			Price: 12,
		},
		Pizza{
			Id:    2,
			Name:  "Capricciosa",
			Price: 11,
		},
		Pizza{
			Id:    3,
			Name:  "Margherita",
			Price: 10,
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/pizzas", pizzasHandler{&pizzas})
	mux.Handle("/orders", ordersHandler{&pizzas, &orders})

	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
