package admins

type CreateCategoryInput struct {
	Name        string
	Description string
}
type VerifyShopInput struct {
	ShopIDs []string
}
