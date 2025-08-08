package fixture

import "github.com/uesleicarvalhoo/aiqfome/product"

type RatingBuilder struct {
	rate  float32
	count int
}

func AnyRating() RatingBuilder {
	return RatingBuilder{
		rate:  4.0,
		count: 5,
	}
}

func (b RatingBuilder) WithRate(r float32) RatingBuilder {
	b.rate = r
	return b
}

func (b RatingBuilder) WithCount(c int) RatingBuilder {
	b.count = c
	return b
}

func (b RatingBuilder) Build() product.Rating {
	return product.Rating{
		Rate:  b.rate,
		Count: b.count,
	}
}
