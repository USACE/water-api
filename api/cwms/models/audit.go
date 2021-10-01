package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditInfo struct {
	Creator    *uuid.UUID `json:"creator"`
	CreateDate time.Time  `json:"create_date"`
	Updater    *uuid.UUID `json:"updater"`
	UpdateDate *time.Time `json:"update_date"`
}
