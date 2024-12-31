package api

import "github.com/gin-gonic/gin"

// ---------------------------------------------- API Models ----------------------------------------------
type paymentReq struct {
	OrderID int64  `json:"order_id" binding:"required,min=1"`
	Source  string `json:"source" binding:"required,oneof=balance card"`
}

// ---------------------------------------------- API Handlers ----------------------------------------------
func (sv *Server) makePayment(c *gin.Context) {

}
