package msq

type Message struct {
	Payload Payload
}

type Payload map[string]interface{}

func (p *Payload) Marshal() ([]byte, error) {

}

func (p *Payload) UnMarshal(data []byte) (*Payload, error) {

}
