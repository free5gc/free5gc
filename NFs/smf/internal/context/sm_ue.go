package context

import "sync"

type UeData struct {
	PduSessionCount   int    // store number of PDU Sessions for each UE
	SdmSubscriptionId string // store SDM Subscription ID per UE
}

type Ues struct {
	ues map[string]UeData // map to store UE data with SUPI as key
	mu  sync.Mutex        // mutex for concurrent access
}

func InitSmfUeData() *Ues {
	return &Ues{
		ues: make(map[string]UeData),
	}
}

// IncrementPduSessionCount increments the PDU session count for a given UE.
func (u *Ues) IncrementPduSessionCount(ueId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	ueData := u.ues[ueId]
	ueData.PduSessionCount++
	u.ues[ueId] = ueData
}

// DecrementPduSessionCount decrements the PDU session count for a given UE.
func (u *Ues) DecrementPduSessionCount(ueId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	ueData := u.ues[ueId]
	if ueData.PduSessionCount > 0 {
		ueData.PduSessionCount--
		u.ues[ueId] = ueData
	}
}

// SetSubscriptionId sets the SDM subscription ID for a given UE.
func (u *Ues) SetSubscriptionId(ueId, subscriptionId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	ueData := u.ues[ueId]
	ueData.SdmSubscriptionId = subscriptionId
	u.ues[ueId] = ueData
}

// GetSubscriptionId returns the SDM subscription ID for a given UE.
func (u *Ues) GetSubscriptionId(ueId string) string {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.ues[ueId].SdmSubscriptionId
}

// GetUeData returns the data for a given UE.
func (u *Ues) GetUeData(ueId string) UeData {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.ues[ueId]
}

// DeleteUe deletes a UE.
func (u *Ues) DeleteUe(ueId string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	delete(u.ues, ueId)
}

// UeExists checks if a UE already exists.
func (u *Ues) UeExists(ueId string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	_, exists := u.ues[ueId]
	return exists
}

// IsLastPduSession checks if it is the last PDU session for a given UE.
func (u *Ues) IsLastPduSession(ueID string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.ues[ueID].PduSessionCount == 1
}

// GetPduSessionCount returns the number of sessions for a given UE.
func (u *Ues) GetPduSessionCount(ueId string) int {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.ues[ueId].PduSessionCount
}
