package cmd

import "github.com/ceph/go-ceph/rgw/admin"

func convertUserCapSpec(input []admin.UserCapSpec) []UserCapSpec {
	var output []UserCapSpec

	for _, capSpec := range input {
		userCap := UserCapSpec{
			Type: capSpec.Type,
			Perm: capSpec.Perm,
		}
		output = append(output, userCap)
	}

	return output
}
