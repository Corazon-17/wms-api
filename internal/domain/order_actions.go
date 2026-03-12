package domain

func AllowedActions(status string) []string {

	switch status {

	case WMSReadyToPick:
		return []string{"pick"}

	case WMSPicking:
		return []string{"pack"}

	case WMSPacked:
		return []string{"ship"}

	default:
		return []string{}
	}
}
