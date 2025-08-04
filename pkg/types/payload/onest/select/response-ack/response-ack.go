package responseack

type SelectResponseAck struct {
	Message Message `json:"message"`
	Error   Error   `json:"error"`
}
type List struct {
	Descriptor string `json:"descriptor"`
	Value      string `json:"value"`
}
type Tags struct {
	Descriptor string `json:"descriptor"`
	List       []List `json:"list"`
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
