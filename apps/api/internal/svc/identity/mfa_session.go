package identity

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type mfaSessionData struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

var mfaSessions sync.Map

func SaveMfaSessionData(token string, userID uuid.UUID) {
	mfaSessions.Store(token, mfaSessionData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})
}

func GetMfaSessionData(token string) (uuid.UUID, bool) {
	val, ok := mfaSessions.Load(token)
	if !ok {
		return uuid.Nil, false
	}
	data := val.(mfaSessionData)
	if time.Now().After(data.ExpiresAt) {
		mfaSessions.Delete(token)
		return uuid.Nil, false
	}
	return data.UserID, true
}

func DeleteMfaSessionData(token string) {
	mfaSessions.Delete(token)
}
