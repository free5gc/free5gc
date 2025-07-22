package context

func GetULCLGroupNameFromSUPI(supi string) string {
	ulclGroups := smfContext.ULCLGroups
	for name, group := range ulclGroups {
		for _, member := range group {
			if member == supi {
				return name
			}
		}
	}
	return ""
}
