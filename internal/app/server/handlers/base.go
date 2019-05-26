package handlers

import (
	"log"

	"go.ybk.im/homepage/internal/app/database"
	"go.ybk.im/homepage/internal/app/types"
)

type handlerBase struct{}

func (b handlerBase) args(args ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	site := new(types.Site)
	err := database.GetSite(site)
	if err != nil {
		log.Panicf("Failed to get site data: %s\n", err)
	}
	result["Site"] = site

	for _, arg := range args {
		for k, v := range arg {
			result[k] = v
		}
	}

	return result
}
