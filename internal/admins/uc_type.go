package admins

type CreateCategoryInput struct {
	Name        string
	Description string
}
type VerifyShopInput struct {
	ShopIDs []string
}

type GetCategoriesFilter struct {
	IDs  []string
	Name string
}
