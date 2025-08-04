package response

import "time"

type SearchResponse struct {
	Context Context `json:"context"`
	Message Message `json:"message"`
}
type ContextCity struct {
	Code string `json:"code"`
}
type Country struct {
	Code string `json:"code"`
}
type Location struct {
	City    ContextCity `json:"city"`
	Country Country     `json:"country"`
}
type Context struct {
	Domain        string   `json:"domain"`
	Action        string   `json:"action"`
	Version       string   `json:"version"`
	BapID         string   `json:"bap_id"`
	BapURI        string   `json:"bap_uri"`
	BppID         string   `json:"bpp_id"`
	BppURI        string   `json:"bpp_uri"`
	TransactionID string   `json:"transaction_id"`
	MessageID     string   `json:"message_id"`
	Location      Location `json:"location"`
	Timestamp     string   `json:"timestamp"`
	TTL           string   `json:"ttl"`
}
type CatalogDescriptor struct {
	Name string `json:"name"`
}
type Images struct {
	URL string `json:"url"`
}
type ProvidersDescriptor struct {
	Name      string   `json:"name"`
	ShortDesc string   `json:"short_desc"`
	Images    []Images `json:"images"`
}
type Fulfillments struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type City struct {
	Code string `json:"code"`
}
type State struct {
	Code string `json:"code"`
}
type AreaCode struct {
	Code string `json:"code"`
}
type Locations struct {
	ID       string   `json:"id"`
	GPS      string   `json:"gps,omitempty"`
	Address  string   `json:"address,omitempty"`
	Street   string   `json:"street,omitempty"`
	AreaCode AreaCode `json:"areaCode"`
	City     City     `json:"city"`
	State    State    `json:"state"`
}
type Media struct {
	Mimetype string `json:"mimetype"`
	URL      string `json:"url"`
}
type ItemsDescriptor struct {
	Name      string   `json:"name"`
	LongDesc  string   `json:"long_desc"`
	ShortDesc string   `json:"short_desc"`
	Images    []Images `json:"images"`
	Media     []Media  `json:"media"`
}
type Available struct {
	Count int `json:"count"`
}
type Quantity struct {
	Available Available `json:"available"`
}
type Range struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
type Time struct {
	Range Range `json:"range"`
}
type CreatorDescriptor struct {
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	ShortDesc string   `json:"short_desc"`
	LongDesc  string   `json:"long_desc"`
	Media     []Media  `json:"media"`
	Images    []Images `json:"images"`
}
type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}
type Creator struct {
	Descriptor CreatorDescriptor `json:"descriptor"`
	Address    string            `json:"address"`
	State      State             `json:"state"`
	City       City              `json:"city"`
	Contact    Contact           `json:"contact"`
}
type TagsDescriptor struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
type List struct {
	Descriptor TagsDescriptor `json:"descriptor"`
	Value      string         `json:"value"`
}
type Tags struct {
	Descriptor TagsDescriptor `json:"descriptor"`
	List       []List         `json:"list"`
}
type Price struct {
	Currency     string `json:"currency"`
	Value        string `json:"value"`
	MaximumValue string `json:"maximum_value"`
}
type Items struct {
	ID             string          `json:"id"`
	CategoryIds    []string        `json:"category_ids"`
	Descriptor     ItemsDescriptor `json:"descriptor"`
	Quantity       Quantity        `json:"quantity"`
	Time           Time            `json:"time"`
	LocationIds    []string        `json:"location_ids"`
	FulfillmentIds []string        `json:"fulfillment_ids"`
	Creator        Creator         `json:"creator"`
	Tags           []Tags          `json:"tags"`
	Price          Price           `json:"price"`
}
type Providers struct {
	ID           string              `json:"id"`
	Descriptor   ProvidersDescriptor `json:"descriptor"`
	Fulfillments []Fulfillments      `json:"fulfillments"`
	Locations    []Locations         `json:"locations"`
	Items        []Items             `json:"items"`
}
type Catalog struct {
	Descriptor CatalogDescriptor `json:"descriptor"`
	Providers  []Providers       `json:"providers"`
}
type Message struct {
	Catalog Catalog `json:"catalog"`
}
