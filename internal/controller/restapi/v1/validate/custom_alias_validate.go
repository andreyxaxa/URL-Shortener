package validate

func IsValidAlias(alias string) bool {
	if len(alias) == 0 {
		return false
	}

	for _, r := range alias {
		if !((r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' || r == '-') {
			return false
		}
	}

	return true
}

func IsValidAliasLength(alias string) bool {
	if len(alias) < 2 || len(alias) > 50 {
		return false
	}

	return true
}
