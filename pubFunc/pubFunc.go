package pubfunc

func IsValidPermission(permission string) bool {
	switch permission {
	case
		"read",
		"write",
		"update",
		"delete":
		return true
	}
	return false
}
func IsValidPtype(ptype string) bool {
	switch ptype {
	case "p", "g":
		return true
	default:
		return false
	}
}
