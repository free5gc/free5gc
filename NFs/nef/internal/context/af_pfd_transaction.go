package context

import (
	"github.com/sirupsen/logrus"
)

type AfPfdTransaction struct {
	TransID   string
	ExtAppIDs map[string]struct{}
	Log       *logrus.Entry
}

func (a *AfPfdTransaction) GetExtAppIDs() []string {
	ids := make([]string, 0, len(a.ExtAppIDs))
	for id := range a.ExtAppIDs {
		ids = append(ids, id)
	}
	return ids
}

func (a *AfPfdTransaction) AddExtAppID(appID string) {
	a.ExtAppIDs[appID] = struct{}{}
	a.Log.Infof("appID[%s] is added", appID)
}

func (a *AfPfdTransaction) DeleteExtAppID(appID string) {
	delete(a.ExtAppIDs, appID)
	a.Log.Infof("appID[%s] is deleted", appID)
}

func (a *AfPfdTransaction) DeleteAllExtAppIDs() {
	a.ExtAppIDs = make(map[string]struct{})
}
