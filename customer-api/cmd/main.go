package main

import (
	"customer-api/database"
	"customer-api/model"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	// Connect Database
	database.ConnectDB() //connect ke database DB
	db = database.DB

	router := gin.Default()

	router.GET("/api/v1/customers", GetCustomers)          // membuat API untuk melihat semua data customers
	router.GET("/api/v1/customers/:id", GetCustomerByID)   // membuat API untuk melihat data customers by id
	router.POST("/api/v1/customers", CreateCustomer)       // untuk create database customer
	router.PUT("/api/v1/customers/:id", UpdateCustomer)    // untuk mengupdate database customer jika ingin diganti
	router.DELETE("/api/v1/customers/:id", DeleteCustomer) // untuk menghapus database customer jika ingin dihapus

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func GetCustomers(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT
			id,
			name,
			email,
			phone,
			created_at,
			updated_at
		FROM customers
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			ResponseCode:    "500",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	defer rows.Close()

	var customers []model.Customer

	for rows.Next() {

		var customer model.Customer

		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.CreatedAt,
			&customer.UpdateAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{
				ResponseCode:    "500",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
			return
		}

		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, model.Response{
		ResponseCode:    "200",
		ResponseMessage: "Successfully retrieved customers.",
		ResponseData:    customers,
	})
}

func GetCustomerByID(c *gin.Context) {

	id := c.Param("id")

	var customer model.Customer

	err := database.DB.QueryRow(`
		SELECT
			id,
			name,
			email,
			phone,
			created_at,
			update_at
		FROM customers
		WHERE id = ?
	`, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdateAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, model.Response{
			ResponseCode:    "404",
			ResponseMessage: "Customer not found",
			ResponseData:    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			ResponseCode:    "500",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		ResponseCode:    "200",
		ResponseMessage: "Customer found.",
		ResponseData:    customer,
	})
}

func CreateCustomer(c *gin.Context) {

	var req model.CreateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			ResponseCode:    "400",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	stmt, err := db.Prepare(`
		INSERT INTO customers(name,email,phone)
		VALUES(?,?,?)
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			ResponseCode:    "500",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		req.Name,
		req.Email,
		req.Phone,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			ResponseCode:    "500",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, model.Response{
		ResponseCode:    "201",
		ResponseMessage: "Customer created successfully",
		ResponseData: gin.H{
			"id": id,
		},
	})
}

func UpdateCustomer(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			ResponseCode:    "400",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	var req model.UpdateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			ResponseCode:    "400",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	result, err := database.DB.Exec(`
		UPDATE customers
		SET
			name = ?,
			email = ?,
			phone = ?
		WHERE id = ?
	`, req.Name, req.Email, req.Phone, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			ResponseCode:    "500",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, model.Response{
			ResponseCode:    "404",
			ResponseMessage: "Customer not found",
			ResponseData:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		ResponseCode:    "200",
		ResponseMessage: "Customer updated successfully",
		ResponseData:    nil,
	})
}

func DeleteCustomer(c *gin.Context) {

	id := c.Param("id")

	result, err := database.DB.Exec(`
		DELETE FROM customers
		WHERE id = ?
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			ResponseCode:    "500",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, model.Response{
			ResponseCode:    "404",
			ResponseMessage: "Customer not found",
			ResponseData:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		ResponseCode:    "200",
		ResponseMessage: "Customer deleted successfully",
		ResponseData:    nil,
	})
}
