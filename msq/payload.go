package msq

import "encoding/json"

type Payload map[string]interface{}

func (p *Payload) Marshal() ([]byte, error) {

	marshalledData, err := json.Marshal(p)

	if err != nil {
		return []byte{}, err
	}

	return marshalledData, nil
}

func (p *Payload) UnMarshal(data []byte) (*Payload, error) {

	err := json.Unmarshal(data, p)

	if err != nil {
		return p, err
	}

	return p, nil
}
