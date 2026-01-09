package messages

type ClientName struct {
	Name string `json:"name"`
}
type ClientNames struct {
	Names []ClientName `json:"names"`
}
