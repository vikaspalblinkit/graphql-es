package graph

import "compensation-api/internal/elastic"

type Resolver struct {
	ES *elastic.Client
}
