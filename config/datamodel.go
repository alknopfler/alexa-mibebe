package config

type Record struct {
	Email  string    `json:"email"`
	Fecha  string	 `json:"fecha"`
	Nombre string    `json:"nombre"`
	Peso   float64   `json:"peso"`
	Toma   float64   `json:"toma"`
}
