package webhook

import "logbook/cmd/objectives/database"

type eventOutdate struct {
	oid           database.ObjectiveId
	current, next database.VersionId
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
