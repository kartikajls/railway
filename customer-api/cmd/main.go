package main

import (
	"customer-api/database"
	"customer-api/model"
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// Connect Database
	database.ConnectDB() //connect ke database DB

	router := gin.Default()
	router.SetTrustedProxies(nil) //digunakan agar GIN dapat berjalan disemua proxy

	api := router.Group("/api/v1") // didefinisikan untuk server localhost sudah beraada pada /api/v1
	{
		api.GET("/customers", GetCustomers)          // membuat API untuk melihat semua data customers
		api.GET("/customers/:id", GetCustomerByID)   // membuat API untuk melihat data customers by id
		api.POST("/customers", CreateCustomer)       // untuk create database customer
		api.PUT("//customers/:id", UpdateCustomer)   // untuk mengupdate database customer jika ingin diganti
		api.DELETE("/customers/:id", DeleteCustomer) // untuk menghapus database customer jika ingin dihapus
	}

	port := os.Getenv("PORT") //di setting untuk di deploy pada railway
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

func GetCustomers(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT
			id,
			name,
			email,
			phone,
			created_at,
			update_at
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

	c.JSON(http.StatusOK, customers)
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

	c.JSON(http.StatusOK, customer)
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

	result, err := database.DB.Exec(`
		INSERT INTO customers(name,email,phone)
		VALUES(?,?,?)
	`, req.Name, req.Email, req.Phone)

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

	id := c.Param("id")

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
