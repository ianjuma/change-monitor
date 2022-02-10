package sa

type ProductEvent struct {
	Operation string `json:"operation"`

	Old struct {
		Active    bool   `json:"active"`
		ProductID int    `json:"product_id"`
		Name      string `json:"name"`
		Quantity  int    `json:"quantity,omitempty"`
		Price     int    `json:"price"`
	} `json:"old"`
	New struct {
		Active    bool   `json:"active"`
		ProductID int    `json:"product_id"`
		Name      string `json:"name"`
		Quantity  int    `json:"quantity,omitempty"`
		Price     int    `json:"price"`
	} `json:"new"`
}
