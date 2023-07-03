package api

import (
	db "Solvery/db/sqlc"
	"Solvery/util"
	"Solvery/util/tasks/task1"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type task1Request struct {
	Email string `json:"email" binding:"required,email"`
	Array []int  `json:"array" binding:"required,task1Comment"`
}

type task1Response struct {
	Result []int    `json:"result"`
	User   db.User  `json:"user"`
	Entry  db.Entry `json:"entry"`
}

func (s *Server) task1(c *gin.Context) {
	var req task1Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(c, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cost := int32(len(req.Array))

	if user.Credit-cost < s.config.MinCredit {
		if err := util.SendEmail(
			s.config.EmailAddress,
			s.config.EmailPassword,
			req.Email,
			task1Comment,
			fmt.Sprintf("sorry, you have no enough credit\nyour credit: %d, task1 costs: %d",
				user.Credit,
				len(req.Array)),
		); err != nil {
			log.Printf("error sending email: %v", err)
		}

		c.JSON(http.StatusBadRequest, errorResponse(errors.New(lowCredit)))
		return
	}

	arg := db.PaymentTxParams{
		Amount:    -cost,
		UserEmail: req.Email,
		Comment:   fmt.Sprintf("%s, input: %v", task1Comment, req.Array),
	}

	res, err := s.store.PaymentTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := task1Response{
		Result: task1.FindMissingNums(req.Array),
		User:   res.User,
		Entry:  res.Entry,
	}

	if err := util.SendEmail(
		s.config.EmailAddress,
		s.config.EmailPassword,
		req.Email,
		task1Comment,
		fmt.Sprintf("%s, input: %v \nresult: %v \nyour credit: %v",
			task1Comment,
			req.Array,
			resp.Result,
			resp.User.Credit,
		)); err != nil {
		log.Printf("error sending email: %v", err)
	}

	c.JSON(http.StatusOK, resp)
}
