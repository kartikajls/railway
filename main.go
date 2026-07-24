package main

import (
	"customer-api/database"
	"database/sql"
	"net/http"
	"os"

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

	router.GET("/customers", GetCustomers)          // melihat keseluruhan database
	router.GET("/customers/:id", GetCustomer)       // melihat database berdasarkan id
	router.POST("/customers", CreateCustomer)       // membuat database
	router.PUT("/customers/:id", UpdateCustomer)    // mengambil database berdasarkan dat yang dibuat di POST
	router.DELETE("/customers/:id", DeleteCustomer) // menghapus database yang diinginkan

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	router.Run(":" + port) // digunakan untuk deploy di port railway
}

// membuat function get customer untuk melihat semua database
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

	id := c.Param("id") // digunakan untuk mengambil nilai parameter pada url customers/id

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

// membuat function create customer untuk membuat database
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

// membuat function update customer untuk mengambil database yang telah dibuat pada function POST
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

// membuat function delete customer untuk menghapus customer
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
