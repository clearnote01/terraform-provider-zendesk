package models

type TicketForm struct {
	ID                 int64   `json:"id,omitempty"`
	URL                string  `json:"url,omitempty"`
	Name               string  `json:"name"`
	RawName            string  `json:"raw_name,omitempty"`
	DisplayName        string  `json:"display_name,omitempty"`
	RawDisplayName     string  `json:"raw_display_name,omitempty"`
	Position           int64   `json:"position"`
	Active             bool    `json:"active,omitempty"`
	EndUserVisible     bool    `json:"end_user_visible,omitempty"`
	Default            bool    `json:"default,omitempty"`
	TicketFieldIDs     []int64 `json:"ticket_field_ids,omitempty"`
	InAllBrands        bool    `json:"in_all_brands,omitempty"`
	RestrictedBrandIDs []int64 `json:"restricted_brand_ids,omitempty"`
}

type TicketFormListOptions struct {
	PageOptions
	Active            bool `url:"active,omitempty"`
	EndUserVisible    bool `url:"end_user_visible,omitempty"`
	FallbackToDefault bool `url:"fallback_to_default,omitempty"`
	AssociatedToBrand bool `url:"associated_to_brand,omitempty"`
}

type PageOptions struct {
	PerPage int `url:"per_page,omitempty"`
	Page    int `url:"page,omitempty"`
}
