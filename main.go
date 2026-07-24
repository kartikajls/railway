package main

import (
	"customer-api/database"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
}

func main() {
	database.ConnectDB()
	router := gin.Default()

	router.GET("/customers", GetCustomers)
	router.GET("/customers/:id", GetCustomer)
	router.POST("/customers", CreateCustomer)
	router.PUT("/customers/:id", UpdateCustomer)
	router.DELETE("/customers/:id", DeleteCustomer)

	router.Run(":8080")
}

func GetCustomers(c *gin.Context) {

	rows, err := database.DB.Query("SELECT id,name,email,phone,created_at,update_at FROM customers")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var customer Customer

		rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.CreatedAt,
			&customer.UpdateAt,
		)

		customers = append(customers, customer)

	}

	c.JSON(http.StatusOK, customers)
}

func GetCustomer(c *gin.Context) {

	id := c.Param("id")

	var customer Customer

	err := database.DB.QueryRow("SELECT	id,name,email,phone,created_at,update_at FROM customers WHERE id=?",
		id,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdateAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Customer not Found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}

	c.JSON(http.StatusOK, customer)
}

func CreateCustomer(c *gin.Context) {

	var customer Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(
		"INSERT INTO customers(name,email,phone) VALUES(?,?,?)",
		customer.Name,
		customer.Email,
		customer.Phone,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created",
	})
}

func UpdateCustomer(c *gin.Context) {

	id := c.Param("id")

	var customer Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(
		"UPDATE customers SET name=?,email=?,phone=? WHERE id=?",
		customer.Name,
		customer.Email,
		customer.Phone,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer updated",
	})
}

func DeleteCustomer(c *gin.Context) {

	id := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM customers WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer deleted",
	})
}
