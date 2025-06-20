package eventsource

import "encoding/json"

type Marshaler interface {
	MarshalES() ([]byte, error)
}

type Unmarshaler interface {
	UnMarshalES(b []byte, object any) error
}

func MarshalES(object any) ([]byte, error) {
	if s, implements := object.(Marshaler); implements {
		return s.MarshalES()
	}

	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func UnmarshalES(b []byte, object any) error {
	if s, implements := object.(Unmarshaler); implements {
		return s.UnMarshalES(b, object)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return err
	}

	return nil
}
