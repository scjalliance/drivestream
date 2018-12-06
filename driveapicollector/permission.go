package driveapicollector

import (
	"fmt"

	"github.com/scjalliance/drivestream/resource"
	drive "google.golang.org/api/drive/v3"
)

// MarshalPermission marshals the given permission as a resource.
func MarshalPermission(perm *drive.Permission) (resource.Permission, error) {
	expiration, err := parseRFC3339(perm.ExpirationTime)
	if err != nil {
		return resource.Permission{}, fmt.Errorf("invalid expiration time: %v", err)
	}

	return resource.Permission{
		ID:           perm.Id,
		Type:         perm.Type,
		EmailAddress: perm.EmailAddress,
		Domain:       perm.Domain,
		Role:         perm.Role,
		DisplayName:  perm.DisplayName,
		Expiration:   expiration,
	}, nil
}
