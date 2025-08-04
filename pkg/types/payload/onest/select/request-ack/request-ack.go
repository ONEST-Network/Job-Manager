package requestack

type SelectRequestAck struct {
	Message Message `json:"message"`
	Error   *Error  `json:"error"`
}
type AdditionalDesc struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}
type Media struct {
	Mimetype  string `json:"mimetype"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
	Dsa       string `json:"dsa"`
}
type Images struct {
	URL      string `json:"url"`
	SizeType string `json:"size_type"`
	Width    string `json:"width"`
	Height   string `json:"height"`
}
type Descriptor struct {
	Name           string         `json:"name"`
	Code           string         `json:"code"`
	ShortDesc      string         `json:"short_desc"`
	LongDesc       string         `json:"long_desc"`
	AdditionalDesc AdditionalDesc `json:"additional_desc"`
	Media          []Media        `json:"media"`
	Images         []Images       `json:"images"`
}
type List struct {
	Descriptor Descriptor `json:"descriptor"`
	Value      string     `json:"value"`
}
type Tags struct {
	Descriptor Descriptor `json:"descriptor"`
	List       []List     `json:"list"`
}
type Ack struct {
	Status string `json:"status"`
	Tags   []Tags `json:"tags"`
}
type Message struct {
	Ack Ack `json:"ack"`
}
type Error struct {
	Code    string `json:"code"`
	Paths   string `json:"paths"`
	Message string `json:"message"`
}
