package cointop

func actionsMap() map[string]bool {
	return map[string]bool{
		"first_page":                      true,
		"help":                            true,
		"last_page":                       true,
		"move_to_page_first_row":          true,
		"move_to_page_last_row":           true,
		"move_to_page_visible_first_row":  true,
		"move_to_page_visible_last_row":   true,
		"move_to_page_visible_middle_row": true,
		"move_up":                         true,
		"move_down":                       true,
		"next_page":                       true,
		"open_link":                       true,
		"page_down":                       true,
		"page_up":                         true,
		"previous_page":                   true,
		"quit":                            true,
		"refresh":                         true,
		"sort_column_1h_change":           true,
		"sort_column_24h_change":          true,
		"sort_column_24h_volume":          true,
		"sort_column_7d_change":           true,
		"sort_column_asc":                 true,
		"sort_column_available_supply":    true,
		"sort_column_desc":                true,
		"sort_column_last_updated":        true,
		"sort_column_market_cap":          true,
		"sort_column_name":                true,
		"sort_column_price":               true,
		"sort_column_rank":                true,
		"sort_column_symbol":              true,
		"sort_column_total_supply":        true,
		"sort_left_column":                true,
		"sort_right_column":               true,
		"toggle_row_chart":                true,
		"open_search":                     true,
		"toggle_favorite":                 true,
		"toggle_show_favorites":           true,
	}
}

func (ct *Cointop) actionExists(action string) bool {
	return ct.actionsmap[action]
}
