package pagination_test

import (
	"fmt"

	"github.com/evrone/go-clean-template/internal/shared/pagination"
)

func ExampleNew() {
	p := pagination.New(0, -10)
	fmt.Printf("Limit: %d, Offset: %d\n", p.Limit, p.Offset)
	// Output:
	// Limit: 20, Offset: 0
}

func ExampleNewPage() {
	page := pagination.NewPage(15, 0, 100)
	fmt.Printf("Limit: %d, Offset: %d, Total: %d\n", page.Limit, page.Offset, page.Total)
	// Output:
	// Limit: 15, Offset: 0, Total: 100
}
