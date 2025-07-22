package notifier

type Notifier struct {
	PfdChangeNotifier *PfdChangeNotifier
}

func NewNotifier() (*Notifier, error) {
	var err error
	n := &Notifier{}
	if n.PfdChangeNotifier, err = NewPfdChangeNotifier(); err != nil {
		return nil, err
	}
	return n, nil
}
