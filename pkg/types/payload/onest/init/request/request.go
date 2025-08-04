package request

type InitRequest struct {
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
type TagsDescriptor struct {
	Code string `json:"code"`
}
type List struct {
	Code  string `json:"code"`
	Value string `json:"value"`
}
type Tags struct {
	Descriptor TagsDescriptor `json:"descriptor"`
	List       []List         `json:"list"`
}
type Items struct {
	ID             string   `json:"id"`
	FulfillmentIds []string `json:"fulfillment_ids"`
	Tags           []Tags   `json:"tags"`
}
type Skills struct {
	Name string `json:"name"`
}
type Languages struct {
	Name string `json:"name"`
}
type CredsDescriptor struct {
	Name      string `json:"name"`
	ShortDesc string `json:"short_desc"`
	LongDesc  string `json:"long_desc"`
}
type Creds struct {
	ID         string          `json:"id"`
	Descriptor CredsDescriptor `json:"descriptor"`
	URL        string          `json:"url"`
	Type       string          `json:"type"`
}
type Person struct {
	Name      string      `json:"name"`
	Gender    string      `json:"gender"`
	Age       string      `json:"age"`
	Skills    []Skills    `json:"skills"`
	Languages []Languages `json:"languages"`
	Creds     []Creds     `json:"creds"`
	Tags      []Tags      `json:"tags"`
}
type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}
type Customer struct {
	Person  Person  `json:"person"`
	Contact Contact `json:"contact"`
}
type Fulfillments struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Customer Customer `json:"customer"`
}
type Order struct {
	Provider     Provider       `json:"provider"`
	Items        []Items        `json:"items"`
	Fulfillments []Fulfillments `json:"fulfillments"`
}
type Message struct {
	Order Order `json:"order"`
}
