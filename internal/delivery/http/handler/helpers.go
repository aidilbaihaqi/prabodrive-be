package handler

import "github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"

func clampPage(page, limit int) (int, int) {
	if page < 1 {
		page = constants.DefaultPage
	}
	if limit < 1 {
		limit = constants.DefaultLimit
	}
	if limit > constants.MaxLimit {
		limit = constants.MaxLimit
	}
	return page, limit
}
