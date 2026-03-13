package domain

func AllowedAction(status string) *string {

	switch status {

	case WMSReadyToPick:
		result := "pick"
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
