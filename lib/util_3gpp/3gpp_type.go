package util_3gpp

type Dnn []uint8

func (d *Dnn) MarshalBinary() (data []byte, err error) {

	data = append(data, uint8(len(*d)))
	data = append(data, (*d)...)

	return data, nil
}

func (d *Dnn) UnmarshalBinary(data []byte) error {

	(*d) = data[1:]
	return nil
}
