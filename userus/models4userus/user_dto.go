package models4userus

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

type WithUserIDs struct {
	UserIDs map[string]string `json:"userIDs,omitempty" firestore:"userIDs,omitempty"`
}

func (v *WithUserIDs) SetUserID(teamID string, userID string) {
	if v.UserIDs == nil {
		v.UserIDs = map[string]string{teamID: userID}
	} else {
		v.UserIDs[teamID] = userID
	}
}

// UserDto is a record that hold information about user
type UserDto struct {
	briefs4contactus.ContactBase
	dbmodels.WithCreated
	dbmodels.WithPreferredLocale
	botsfwmodels.WithBotUserIDs

	IsAnonymous bool `json:"isAnonymous" firestore:"isAnonymous"`
	//Title       string `json:"title,omitempty" firestore:"title,omitempty"`

	Timezone *dbmodels.Timezone `json:"timezone,omitempty" firestore:"timezone,omitempty"`

	Defaults *UserDefaults `json:"defaults,omitempty" firestore:"defaults,omitempty"`

	Email         string `json:"email,omitempty"  firestore:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified"  firestore:"emailVerified"`

	// List of teams a user belongs to
	Teams   map[string]*UserTeamBrief `json:"teams,omitempty"   firestore:"teams,omitempty"`
	TeamIDs []string                  `json:"teamIDs,omitempty" firestore:"teamIDs,omitempty"`

	Created dbmodels.CreatedInfo `json:"created" firestore:"created"`
	// TODO: Should this be moved to company members?
	//models.DatatugUser
}

func (v *UserDto) GetFullName() string {
	return v.Name.GetFullName()
}

// SetTeamBrief sets team brief and adds teamID to the list of team IDs if needed
func (v *UserDto) SetTeamBrief(teamID string, brief *UserTeamBrief) (updates []dal.Update) {
	if v.Teams == nil {
		v.Teams = map[string]*UserTeamBrief{teamID: brief}
	} else {
		v.Teams[teamID] = brief
	}
	updates = append(updates, dal.Update{Field: "teams." + teamID, Value: brief})
	if !slice.Contains(v.TeamIDs, teamID) {
		v.TeamIDs = append(v.TeamIDs, teamID)
		updates = append(updates, dal.Update{Field: "teamIDs", Value: v.TeamIDs})
	}
	return
}

// GetTeamBriefByType returns the first team brief that matches a specific type
func (v *UserDto) GetTeamBriefByType(t core4teamus.TeamType) (teamID string, teamBrief *UserTeamBrief) {
	for id, brief := range v.Teams {
		if brief.Type == t {
			return id, brief
		}
	}
	return "", nil
}

// Validate validates user record
func (v *UserDto) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return err
	}
	//if v.Avatar != nil {
	//	if err := v.Avatar.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("avatar", err.Error())
	//	}
	//}
	//if v.Title != "" {
	//	if err := v.Name.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.validateEmails(); err != nil {
		return err
	}
	if err := v.validateTeams(); err != nil {
		return err
	}
	if err := dbmodels.ValidateGender(v.Gender, true); err != nil {
		return err
	}
	//if v.Datatug != nil {
	//	if err := v.Datatug.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}
	return nil
}

func (v *UserDto) validateEmails() error {
	if strings.TrimSpace(v.Email) != v.Email {
		return validation.NewErrBadRecordFieldValue("email", "contains leading or closing spaces")
	}
	if strings.Contains(v.Email, " ") {
		return validation.NewErrBadRecordFieldValue("email", "contains space")
	}
	if v.Email != "" {
		if _, err := mail.ParseAddress(v.Email); err != nil {
			return validation.NewErrBadRecordFieldValue("email", err.Error())
		}
		if len(v.Emails) == 0 {
			return validation.NewErrBadRecordFieldValue("emails", "user record has 'email' value but 'emails' are empty")
		}
	}
	primaryEmailInEmails := false
	for i, email := range v.Emails {
		if err := email.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("emails[%v]", i), err.Error())
		}
		if email.Address == v.Email {
			primaryEmailInEmails = true
		}
	}
	if v.Email != "" && !primaryEmailInEmails {
		return validation.NewErrBadRecordFieldValue("emails", "user's primary email is not in 'emails' field")
	}
	return nil
}

func (v *UserDto) validateTeams() error {
	if len(v.Teams) != len(v.TeamIDs) {
		return validation.NewErrBadRecordFieldValue("teamIDs",
			fmt.Sprintf("len(v.Teams) != len(v.TeamIDs): %v != %v", len(v.Teams), len(v.TeamIDs)))
	}
	if len(v.Teams) > 0 {
		teamIDs := make([]string, 0, len(v.Teams))
		teamTitles := make([]string, 0, len(v.Teams))
		for teamID, t := range v.Teams {
			if teamID == "" {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("teams['%v']", teamID), "holds empty id")
			}
			if !slice.Contains(v.TeamIDs, teamID) {
				return validation.NewErrBadRecordFieldValue("teamIDs", "missing team ID: "+teamID)
			}
			if err := t.Validate(); err != nil {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("teams[%s]{title=%v}", teamID, t.Title), err.Error())
			}
			if len(v.Teams) > i {
				for i, title := range teamTitles {
					if t.Title == title {
						return validation.NewErrBadRecordFieldValue("teams",
							fmt.Sprintf("at least 2 teams (%s & %s) with same title: %s", teamID, teamIDs[i], title))
					}
				}
			}
			teamIDs = append(teamIDs, teamID)
			teamTitles = append(teamIDs, t.Title)
		}
	}
	return nil
}

// GetUserTeamInfoByID returns team info specific to the user by team ID
func (v *UserDto) GetUserTeamInfoByID(teamID string) *UserTeamBrief {
	return v.Teams[teamID]
}
