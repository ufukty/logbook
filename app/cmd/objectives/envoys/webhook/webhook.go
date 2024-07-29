package webhook

import "logbook/models/columns"

type eventOutdate struct {
	oid           columns.ObjectiveId
	current, next columns.VersionId
}

type Webhook struct {
	queue chan eventOutdate
}

func (wh *Webhook) handleOutdateEvent(e eventOutdate) error {

	return nil
}

func (wh *Webhook) loop() {
	for item := range wh.queue {
		wh.handleOutdateEvent(item)
	}
}

func (wh *Webhook) Event(e eventOutdate) {
	wh.queue <- e
}

func New() *Webhook {
	wh := &Webhook{}
	go wh.loop()
	return wh
}
