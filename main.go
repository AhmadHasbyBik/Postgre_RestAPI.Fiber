package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akhil/go-fiber-postgres/models"
	"github.com/akhil/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Car struct {
	Name string `json:"name"`
	Year string `json:"year"`
	Type string `json:"type"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateCar(context *fiber.Ctx) error {
	car := Car{}

	err := context.BodyParser(&car)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&car).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldnt create car"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "the car has been added"})
	return nil
}

func (r *Repository) deleteCar(context *fiber.Ctx) error {
	carModel := models.Cars{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cant be empty"})
		return nil
	}

	err := r.DB.Delete(carModel, id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldnt delete book"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "the car has been deleted"})
	return nil
}

func (r *Repository) getCars(context *fiber.Ctx) error {
	carModels := &[]models.Cars{}

	err := r.DB.Find(carModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldn get the car"})
		return nil
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "cars fetched successfully",
			"data":    carModels,
		})
	return nil
}

func (r *Repository) getCarById(context *fiber.Ctx) error {
	carModel := models.Cars{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cant be empty"})
		return nil
	}

	fmt.Println("The ID is", id)

	err := r.DB.Where("id = ?", id).First(carModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "couldn get the car"})
		return nil
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "cars fetched successfully",
			"data":    carModel,
		})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/apis")
	api.Post("/api", r.CreateCar)
	api.Delete("/api/:id", r.deleteCar)
	api.Get("/api/:id", r.getCarById)
	api.Get("/api", r.getCars)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("couldnt load the DB")
	}

	err = models.MigrateCars(db)
	if err != nil {
		log.Fatal("couldnt migrate the DB")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
