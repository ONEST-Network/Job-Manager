package response

type SelectResponse struct {
	Context Context `json:"context"`
	Message Message `json:"message"`
}
type City struct {
	Code string `json:"code"`
}
type Country struct {
	Code string `json:"code"`
}
type Location struct {
	City    City    `json:"city"`
	Country Country `json:"country"`
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
type Provider struct {
	ID string `json:"id"`
}
type Fulfillments struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Descriptor struct {
	Code string `json:"code"`
}
type List struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}
type Tags struct {
	Descriptor Descriptor `json:"descriptor"`
	List       []List     `json:"list"`
}
type Items struct {
	ID             string   `json:"id"`
	FulfillmentIds []string `json:"fulfillment_ids"`
	Tags           []Tags   `json:"tags"`
}
type Price struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}
type Item struct {
	ID    string `json:"id"`
	Price Price  `json:"price"`
	Title string `json:"title"`
}
type Breakup struct {
	Item Item `json:"item"`
}
type Quote struct {
	Price   Price     `json:"price"`
	Breakup []Breakup `json:"breakup"`
	TTL     string    `json:"ttl"`
}
type Order struct {
	Provider     Provider       `json:"provider"`
	Fulfillments []Fulfillments `json:"fulfillments"`
	Items        []Items        `json:"items"`
	Quote        Quote          `json:"quote"`
}
type Message struct {
	Order Order `json:"order"`
}
