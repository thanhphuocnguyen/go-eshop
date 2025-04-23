package api

type ProductParam struct {
	ID string `uri:"id" binding:"required"`
}
