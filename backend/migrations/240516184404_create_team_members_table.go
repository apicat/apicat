package migrations

import (
	"time"

	"github.com/apicat/apicat/v2/backend/model"
	"github.com/apicat/apicat/v2/backend/model/team"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "240516184404",
		Migrate: func(tx *gorm.DB) error {

			type TeamMember struct {
				ID              uint      `gorm:"primarykey"`
				TeamID          string    `gorm:"type:varchar(24);uniqueIndex:ukey;not null;comment:team id"`
				UserID          uint      `gorm:"type:bigint;uniqueIndex:ukey;not null;comment:user id"`
				Role            team.Role `gorm:"type:varchar(32);comment:team member role"`
				Status          string    `gorm:"type:varchar(32);default:active;comment:team member status"`
				InvitationToken string    `gorm:"type:varchar(32);index;comment:invitation code"`
				InvitedBy       uint      `gorm:"type:bigint;default:0;comment:invited by member id"`
				LastActiveAt    time.Time `gorm:"type:datetime;not null;comment:last active time"`
				model.TimeModel
			}

			if tx.Migrator().HasTable(&TeamMember{}) {
				return nil
			}
			return tx.Migrator().CreateTable(&TeamMember{})
		},
	}

	MigrationHelper.Register(m)
}
