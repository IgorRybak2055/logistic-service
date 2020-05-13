package client

import (
	"log"
	"net/http"
	"strconv"
)

func (c *Client) handleMain(w http.ResponseWriter, r *http.Request) {
	ds, err := c.Deliveries()
	if err != nil {
		log.Println("failed to getting deliveries:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = mainTmpl.Execute(w, ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Client) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := loginTmpl.Execute(w, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		log.Println("Login")

		user := User{
			Email:    r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		user, err := c.Login(user)
		if err != nil {
			log.Println("error login:", err)
			return
		}

		log.Println("End login:", user.Name, user.CompanyId, user.ID)
		http.Redirect(w, r, "/main", http.StatusSeeOther)

		// if user.CompanyId == 1 {
		// 	log.Println("transport")
		// 	http.Redirect(w, r, "/m/main", http.StatusSeeOther)
		// } else {
		// 	log.Println("owner")
		// 	http.Redirect(w, r, "/m/driver/main", http.StatusSeeOther)
		// }
	}
}

func (c *Client) handleAddDelivery(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ds, err := c.Deliveries()
		if err != nil {
			log.Println("failed to getting deliveries:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := addDeliveryTmpl.Execute(w, ds); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		log.Println("add dlv")

		date :=  r.FormValue("shipment_date")
		// tmp := strings.Split(date, ".")

		// if err != nil {
		// 	log.Println("failed to parse date in add dlv: ", err)
		// 	return
		// }

		weight, err := strconv.ParseFloat(r.FormValue("weight"), 10)
		if err != nil {
			log.Println("failed to parse weight in add dlv: ", err)
			return

		}

		volume, err := strconv.ParseFloat(r.FormValue("volume"), 10)
		if err != nil {
			log.Println("failed to parse volume in add dlv: ", err)
			return

		}

		price, err := strconv.ParseFloat(r.FormValue("price"), 10)
		if err != nil {
			log.Println("failed to parse price in add dlv: ", err)
			return

		}

		dlv := Delivery{
			ShipmentDate:   date,
			ShipmentPlace:  r.FormValue("depart"),
			UnloadingPlace: r.FormValue("arrival"),
			Cargo:          r.FormValue("cargo"),
			WeightCargo:    weight,
			VolumeCargo:    volume,
			TrailerType:    r.FormValue("type"),
			Price:          price,
		}

		delivery, err := c.CreateDelivery(dlv)
		if err != nil {
			log.Println("error create delivery:", err)
			return
		}

		log.Println("Added delivery:", delivery.ID)
		http.Redirect(w, r, "/main", http.StatusSeeOther)

		// if user.CompanyId == 1 {
		// 	log.Println("transport")
		// 	http.Redirect(w, r, "/m/main", http.StatusSeeOther)
		// } else {
		// 	log.Println("owner")
		// 	http.Redirect(w, r, "/m/driver/main", http.StatusSeeOther)
		// }
	}
}

func (c *Client) handleActiveTenders(w http.ResponseWriter, r *http.Request) {
	ds, err := c.ActiveDeliveries()
	if err != nil {
		log.Println("failed to getting deliveries:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = activeTenders.Execute(w, ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Client) handleDeliveryInfo(w http.ResponseWriter, r *http.Request) {
	id, err	 :=  strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		return
	}
	ds, err := c.Delivery(id)
	if err != nil {
		log.Println("failed to getting deliveries:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = mainTmpl.Execute(w, ds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}