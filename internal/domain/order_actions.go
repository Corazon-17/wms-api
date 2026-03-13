package domain

func AllowedActions(status string) *string {

	switch status {

	case WMSReadyToPick:
		result := "pickup"
		return &result

	case WMSPicking:
		result := "pack"
		return &result

	case WMSPacked:
		result := "ship"
		return &result

	default:
		return nil
	}
}
