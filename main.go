package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type Product struct {
	ID    uint64  `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float64 `json:"price"`
}

func main() {
	var err error
	dsn := "host=192.168.88.18 user=postgres password=postgres dbname=dbtest port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/product/:id", getItemById)
	router.GET("/products", getAllItems)
	router.POST("/product/:id", updateItemById)
	router.POST("/products", addItem)
	router.DELETE("/product/:id", delItemById)

	router.Run()
}

func addItem(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// INSERT INTO products (name, brand, price) VALUES ('a', 'b', 1.23);
	result := db.Create(&product)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &product)
}

func delItemById(c *gin.Context) {
	id := c.Param("id")

	// DELETE FROM products WHERE id=<id>;
	result := db.Delete(&Product{}, id)

	if result.RowsAffected == 0 || result.Error != nil {
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, &gin.H{
				"error": "no record",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "record deleted",
	})
}

func updateItemById(c *gin.Context) {
	var product Product
	id := c.Param("id")

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// UPDATE products SET name='a', brand='b', price=1.23 WHERE id=<id>;
	result := db.Where("id=?", id).Updates(&product)

	if result.RowsAffected == 0 || result.Error != nil {
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, &gin.H{
				"error": "no record",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, &product)
}

func getItemById(c *gin.Context) {
	var product Product
	id := c.Param("id")

	// SELECT * FROM products WHERE id=<id>;
	result := db.Where("id=?", id).Find(&product)

	if result.RowsAffected == 0 || result.Error != nil {
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, &gin.H{
				"error": "no result",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, &product)
}

func getAllItems(c *gin.Context) {
	var product []Product

	// SELECT * FROM products ORDER BY id;
	result := db.Order("id").Find(&product)

	if result.RowsAffected == 0 || result.Error != nil {
		if result.RowsAffected == 0 {
			c.JSON(http.StatusOK, &gin.H{
				"error": "no result",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, &product)
}
