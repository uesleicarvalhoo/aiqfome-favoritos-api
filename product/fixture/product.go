package fixture

import "github.com/uesleicarvalhoo/aiqfome/product"

type ProductBuilder struct {
	id          int
	title       string
	price       float32
	description string
	category    string
	imageUrl    string
	rating      product.Rating
}

func AnyProduct() ProductBuilder {
	return ProductBuilder{
		id:          1,
		title:       "Sample Product",
		price:       99.99,
		description: "A sample product description",
		category:    "Sample Category",
		imageUrl:    "http://example.com/image.png",
		rating:      AnyRating().Build(),
	}
}

func (b ProductBuilder) WithID(id int) ProductBuilder {
	b.id = id
	return b
}

func (b ProductBuilder) WithTitle(t string) ProductBuilder {
	b.title = t
	return b
}

func (b ProductBuilder) WithPrice(p float32) ProductBuilder {
	b.price = p
	return b
}

func (b ProductBuilder) WithDescription(desc string) ProductBuilder {
	b.description = desc
	return b
}

func (b ProductBuilder) WithCategory(cat string) ProductBuilder {
	b.category = cat
	return b
}

func (b ProductBuilder) WithImageUrl(url string) ProductBuilder {
	b.imageUrl = url
	return b
}

func (b ProductBuilder) WithRating(r product.Rating) ProductBuilder {
	b.rating = r
	return b
}

func (b ProductBuilder) Build() product.Product {
	return product.Product{
		ID:          b.id,
		Title:       b.title,
		Price:       b.price,
		Description: b.description,
		Category:    b.category,
		ImageUrl:    b.imageUrl,
		Rating:      b.rating,
	}
}
