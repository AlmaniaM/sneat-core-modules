package briefs4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/validation"
	"strings"
)

// ContactBrief needed as ContactBase is used in models4contactus.ContactDto and in dto4contactus.CreatePersonRequest
// Status is not part of ContactBrief as we keep in briefs only active contacts
type ContactBrief struct {
	dbmodels.WithUserID
	dbmodels.WithOptionalRelatedAs // This is used in `Related` field of `ContactDto`
	dbmodels.WithOptionalCountryID
	dbmodels.WithRoles

	Type       ContactType     `json:"type" firestore:"type"` // "person", "company", "location"
	Gender     dbmodels.Gender `json:"gender,omitempty" firestore:"gender,omitempty"`
	Name       *dbmodels.Name  `json:"name,omitempty" firestore:"name,omitempty"`
	Title      string          `json:"title,omitempty" firestore:"title,omitempty"`
	ShortTitle string          `json:"shortTitle,omitempty" firestore:"shortTitle,omitempty"` // Not supposed to be used in models4contactus.ContactDto
	ParentID   string          `json:"parentID" firestore:"parentID"`                         // Intentionally not adding `omitempty` so we can search root contacts only

	// Number of active invites to join a team
	InvitesCount int `json:"activeInvitesCount,omitempty" firestore:"activeInvitesCount,omitempty"`

	// AgeGroup is deprecated?
	AgeGroup string `json:"ageGroup,omitempty" firestore:"ageGroup,omitempty"` // TODO: Add validation
	PetKind  string `json:"species,omitempty" firestore:"species,omitempty"`

	// Avatar holds a photo of a member
	Avatar *dbprofile.Avatar `json:"avatar,omitempty" firestore:"avatar,omitempty"`
}

func (v *ContactBrief) SetNames(first, last string) {
	v.Name.First = first
	v.Name.Last = last
	//v.User = user
}

func (v *ContactBrief) IsTeamMember() bool {
	return v.HasRole(const4contactus.TeamMemberRoleMember)
}

// GetUserID returns UserID field value
func (v *ContactBrief) GetUserID() string {
	return v.UserID
}

// Equal returns true if 2 instances are equal
func (v *ContactBrief) Equal(v2 *ContactBrief) bool {
	return v.Type == v2.Type &&
		v.WithUserID == v2.WithUserID &&
		v.Gender == v2.Gender &&
		v.WithOptionalCountryID == v2.WithOptionalCountryID &&
		v.Name.Equal(v2.Name) &&
		v.WithOptionalRelatedAs.Equal(v2.WithOptionalRelatedAs) &&
		v.Avatar.Equal(v2.Avatar)
}

// Validate returns error if not valid
func (v *ContactBrief) Validate() error {
	if err := ValidateContactType(v.Type); err != nil {
		return err
	}
	if err := dbmodels.ValidateGender(v.Gender, false); err != nil {
		return err
	}
	if strings.TrimSpace(v.Title) == "" && v.Name == nil {
		return validation.NewErrRecordIsMissingRequiredField("name|title")
	} else if err := v.Name.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("name", err.Error())
	}
	if v.UserID != "" {
		if !core.IsAlphanumericOrUnderscore(v.UserID) {
			return validation.NewErrBadRecordFieldValue("userID", "is not alphanumeric: "+v.UserID)
		}
	}
	switch v.Type {
	case ContactTypeLocation:
		if v.ParentID == "" {
			return validation.NewErrRecordIsMissingRequiredField("parentID")
		}
	}
	if err := v.WithOptionalCountryID.Validate(); err != nil {
		return err
	}
	if err := v.WithRoles.Validate(); err != nil {
		return err
	}
	if err := v.WithUserID.Validate(); err != nil {
		return err
	}
	if v.PetKind != "" {
		if !const4contactus.IsKnownPetPetKind(v.PetKind) {
			return validation.NewErrBadRecordFieldValue("species", "unknown value: "+v.PetKind)
		}
	}
	return nil
}

// GetTitle return full name of a person
func (v *ContactBrief) GetTitle() string {
	if v.Title != "" {
		return v.Title
	}
	if v.Name.Full != "" {
		return v.Name.Full
	}
	if v.Name.First != "" && v.Name.Last != "" && v.Name.Middle != "" {
		return v.Name.First + " " + v.Name.Middle + " " + v.Name.Full
	}
	if v.Name.First != "" && v.Name.Last != "" {
		return v.Name.First + " " + v.Name.Full
	}
	if v.Name.First != "" {
		return v.Name.First
	}
	if v.Name.Last != "" {
		return v.Name.Last
	}
	if v.Name.Middle != "" {
		return v.Name.Middle
	}
	return ""
}

func (v *ContactBrief) DetermineShortTitle(title string, contacts map[string]*ContactBrief) string {
	if v.Name.First != "" && IsUniqueShortTitle(v.Name.First, contacts, const4contactus.TeamMemberRoleMember) {
		v.ShortTitle = v.Name.First
	} else if v.Name.Nick != "" && IsUniqueShortTitle(v.Name.First, contacts, const4contactus.TeamMemberRoleMember) {
		return v.Name.Nick
	} else if v.Name.Full != "" {
		return getShortTitle(v.Name.Full, contacts)
	} else if title != "" {
		return getShortTitle(title, contacts)
	}
	return ""
}

func getShortTitle(title string, members map[string]*ContactBrief) string {
	shortNames := GetShortNames(title)
	for _, short := range shortNames {
		isUnique := true
		for _, m := range members {
			if m.ShortTitle == short.Name {
				isUnique = false
				break
			}
		}
		if isUnique {
			return short.Name
		}
	}
	return ""
}

type ShortName struct {
	Name string `json:"name" firestore:"name"`
	Type string `json:"type" firestore:"type"`
}

// GetShortNames returns short names from a title
func GetShortNames(title string) (shortNames []ShortName) {
	title = CleanTitle(title)
	names := strings.Split(title, " ")
	shortNames = make([]ShortName, 0, len(names))
NAMES:
	for _, s := range names {
		name := strings.TrimSpace(s)
		if name == "" {
			continue
		}
		for _, sn := range shortNames {
			if sn.Name == name {
				continue NAMES
			}
		}
		shortNames = append(shortNames, ShortName{
			Name: name,
			Type: "unknown",
		})
	}
	return shortNames
}

// CleanTitle cleans title from spaces
func CleanTitle(title string) string {
	title = strings.TrimSpace(title)
	for strings.Contains(title, "  ") {
		title = strings.Replace(title, "  ", " ", -1)
	}
	return title
}
