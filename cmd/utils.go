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
func bytesToKB(bytes int64) float64 {
	const KB = 1024
	return float64(bytes) / float64(KB)
}
