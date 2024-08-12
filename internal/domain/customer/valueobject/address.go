package valueobject

type Address struct {
	Street  string
	City    string
	State   string
	Country Country
}
