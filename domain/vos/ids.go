package vos

type SubscriberID string

func (s SubscriberID) String() string {
	return string(s)
}
