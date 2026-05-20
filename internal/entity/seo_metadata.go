package entity

type SeoMetadata struct {
	ID              string `json:"id"`
	OwnerType       string `json:"owner_type"`
	OwnerID         string `json:"owner_id"`
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	Slug            string `json:"slug"`
	AltText         string `json:"alt_text,omitempty"`
}
