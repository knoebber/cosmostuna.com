package main

type state struct {
	abr  string
	full string
}

// Currently only allow shipping to lower 48 states; Alaska and Hawaii are absent from this list.
var states = []state{
	{abr: "AL", full: "Alabama"},
	{abr: "AZ", full: "Arizona"},
	{abr: "AR", full: "Arkansas"},
	{abr: "CA", full: "California"},
	{abr: "CO", full: "Colorado"},
	{abr: "CT", full: "Connecticut"},
	{abr: "DE", full: "Delaware"},
	{abr: "DC", full: "District Of Columbia"},
	{abr: "FL", full: "Florida"},
	{abr: "GA", full: "Georgia"},
	{abr: "ID", full: "Idaho"},
	{abr: "IL", full: "Illinois"},
	{abr: "IN", full: "Indiana"},
	{abr: "IA", full: "Iowa"},
	{abr: "KS", full: "Kansas"},
	{abr: "KY", full: "Kentucky"},
	{abr: "LA", full: "Louisiana"},
	{abr: "ME", full: "Maine"},
	{abr: "MD", full: "Maryland"},
	{abr: "MA", full: "Massachusetts"},
	{abr: "MI", full: "Michigan"},
	{abr: "MN", full: "Minnesota"},
	{abr: "MS", full: "Mississippi"},
	{abr: "MO", full: "Missouri"},
	{abr: "MT", full: "Montana"},
	{abr: "NE", full: "Nebraska"},
	{abr: "NV", full: "Nevada"},
	{abr: "NH", full: "New Hampshire"},
	{abr: "NJ", full: "New Jersey"},
	{abr: "NM", full: "New Mexico"},
	{abr: "NY", full: "New York"},
	{abr: "NC", full: "North Carolina"},
	{abr: "ND", full: "North Dakota"},
	{abr: "OH", full: "Ohio"},
	{abr: "OK", full: "Oklahoma"},
	{abr: "OR", full: "Oregon"},
	{abr: "PA", full: "Pennsylvania"},
	{abr: "RI", full: "Rhode Island"},
	{abr: "SC", full: "South Carolina"},
	{abr: "SD", full: "South Dakota"},
	{abr: "TN", full: "Tennessee"},
	{abr: "TX", full: "Texas"},
	{abr: "UT", full: "Utah"},
	{abr: "VT", full: "Vermont"},
	{abr: "VA", full: "Virginia"},
	{abr: "WA", full: "Washington"},
	{abr: "WV", full: "West Virginia"},
	{abr: "WI", full: "Wisconsin"},
	{abr: "WY", full: "Wyoming"},
}

func shippable(stateName *string) bool {
	if stateName == nil {
		return false
	}
	for _, s := range states {
		if *stateName == s.abr || *stateName == s.full {
			return true
		}
	}
	return false
}
