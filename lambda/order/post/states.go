package main

import (
	"strings"
)

type state struct {
	abr  string
	full string
}

// Currently only allow shipping to lower 48 states; Alaska and Hawaii are absent from this list.
var states = []state{
	{abr: "al", full: "alabama"},
	{abr: "az", full: "arizona"},
	{abr: "ar", full: "arkansas"},
	{abr: "ca", full: "california"},
	{abr: "co", full: "colorado"},
	{abr: "ct", full: "connecticut"},
	{abr: "de", full: "delaware"},
	{abr: "dc", full: "district of columbia"},
	{abr: "fl", full: "florida"},
	{abr: "ga", full: "georgia"},
	{abr: "id", full: "idaho"},
	{abr: "il", full: "illinois"},
	{abr: "in", full: "indiana"},
	{abr: "ia", full: "iowa"},
	{abr: "ks", full: "kansas"},
	{abr: "ky", full: "kentucky"},
	{abr: "la", full: "louisiana"},
	{abr: "me", full: "maine"},
	{abr: "md", full: "maryland"},
	{abr: "ma", full: "massachusetts"},
	{abr: "mi", full: "michigan"},
	{abr: "mn", full: "minnesota"},
	{abr: "ms", full: "mississippi"},
	{abr: "mo", full: "missouri"},
	{abr: "mt", full: "montana"},
	{abr: "ne", full: "nebraska"},
	{abr: "nv", full: "nevada"},
	{abr: "nh", full: "new hampshire"},
	{abr: "nj", full: "new jersey"},
	{abr: "nm", full: "new mexico"},
	{abr: "ny", full: "new york"},
	{abr: "nc", full: "north carolina"},
	{abr: "nd", full: "north dakota"},
	{abr: "oh", full: "ohio"},
	{abr: "ok", full: "oklahoma"},
	{abr: "or", full: "oregon"},
	{abr: "pa", full: "pennsylvania"},
	{abr: "ri", full: "rhode island"},
	{abr: "sc", full: "south carolina"},
	{abr: "sd", full: "south dakota"},
	{abr: "tn", full: "tennessee"},
	{abr: "tx", full: "texas"},
	{abr: "ut", full: "utah"},
	{abr: "vt", full: "vermont"},
	{abr: "va", full: "virginia"},
	{abr: "wa", full: "washington"},
	{abr: "wv", full: "west virginia"},
	{abr: "wi", full: "wisconsin"},
	{abr: "wy", full: "wyoming"},
}

func shippable(stateName *string) bool {
	if stateName == nil {
		return false
	}

	val := strings.TrimSpace(strings.ToLower(*stateName))
	for _, s := range states {
		if val == s.abr || val == s.full {
			return true
		}
	}
	return false
}
