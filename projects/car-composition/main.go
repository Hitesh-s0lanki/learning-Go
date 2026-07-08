package main

import "fmt"

type Engine struct {
	HorsePower int
	FuelType   string
}

func (engine Engine) Start() {
	fmt.Println("Engine started")
}

type Wheel struct {
	Size int
}

type GPS struct {
	Brand string
}

func (gps GPS) Navigate(destination string) {
	fmt.Println("Navigating to", destination, "using", gps.Brand)
}

type Car struct {
	// Embedded struct: Car gets direct access to Engine fields and methods.
	Engine

	// Regular fields: these are accessed using their field names.
	Wheels []Wheel
	GPS    GPS
	Brand  string
	Model  string
}

func NewCar(brand string, model string, horsePower int, fuelType string) Car {
	return Car{
		Brand: brand,
		Model: model,
		Engine: Engine{
			HorsePower: horsePower,
			FuelType:   fuelType,
		},
		Wheels: []Wheel{
			{Size: 17},
			{Size: 17},
			{Size: 17},
			{Size: 17},
		},
		GPS: GPS{Brand: "Google Maps"},
	}
}

func (car Car) Drive(destination string) {
	car.Start()
	car.GPS.Navigate(destination)
	fmt.Println(car.Brand, car.Model, "is driving")
}

func main() {
	car := NewCar("Toyota", "Fortuner", 201, "Diesel")

	fmt.Println("Car:", car.Brand, car.Model)
	fmt.Println("Engine horsepower:", car.HorsePower)
	fmt.Println("Engine fuel type:", car.FuelType)
	fmt.Println("Number of wheels:", len(car.Wheels))
	fmt.Println("First wheel size:", car.Wheels[0].Size)

	car.Drive("Mumbai")
}
