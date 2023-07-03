package api

import (
	db "Solvery/db/sqlc"
	"Solvery/util"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"log"
	"net/http"
)

type createUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Class    string `json:"class" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type updateUserCreditRequest struct {
	Amount int32  `json:"amount" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

type getUserRequest struct {
	Email string `uri:"email" binding:"required,email"`
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:  req.FullName,
		Class: req.Class,
		Email: req.Email,
	}

	user, err := s.store.CreateUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) getUser(c *gin.Context) {
	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
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

	c.JSON(http.StatusOK, user)
}

func (s *Server) listUsers(c *gin.Context) {
	var req listUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := s.store.ListUsers(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *Server) updateUserCredit(c *gin.Context) {
	var req updateUserCreditRequest
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

	if user.Credit+req.Amount < s.config.MinCredit {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New(lowCredit)))
		return
	}

	arg := db.PaymentTxParams{
		Amount:    req.Amount,
		UserEmail: req.Email,
		Comment:   updateCreditComment,
	}

	res, err := s.store.PaymentTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.SendEmail(
		s.config.EmailAddress,
		s.config.EmailPassword,
		req.Email,
		updateCreditComment,
		fmt.Sprintf("your credit has been updated from %d to %d",
			user.Credit, res.User.Credit),
	); err != nil {
		log.Printf("error sending email: %v", err)
	}

	c.JSON(http.StatusOK, res)
}
