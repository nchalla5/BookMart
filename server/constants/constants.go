package constants

// ProductStatus represents various states a product can have.
type ProductStatus string

const (
	// Available indicates that a product is available for sale.
	Available ProductStatus = "available"

	// Sold indicates that a product has been sold.
	Sold ProductStatus = "sold"

	// Unavailable indicates that a product is no longer available for sale.
	Unavailable ProductStatus = "unavailable"
)
