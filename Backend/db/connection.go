package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func DbConnection() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Errorf("error cargando .env: %w", err))
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
	)
	fmt.Println(dsn)

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("Error al conectarse a la DB: %v", err)
	}

	log.Println("DB Conectada")
}

func InitMockDb() *gorm.DB {
	// se ctrea un mock de la db
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error al crear el mock de la base de datos: %v", err)
	}

	// Crear una conexión GORM al mock
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al abrir GORM: %v", err)
	}

	//se debe crear un mocck para cada consulta
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"stocks\" WHERE ticker ILIKE \\$1 OR company ILIKE \\$2 OR brokerage ILIKE \\$3").
		WithArgs("%BSBR%", "%BSBR%", "%BSBR%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2)) // se simula que hay 2 registros

	mock.ExpectQuery("SELECT \\* FROM \"stocks\" WHERE ticker ILIKE \\$1 OR company ILIKE \\$2 OR brokerage ILIKE \\$3 LIMIT \\$4").
		WithArgs("%BSBR%", "%BSBR%", "%BSBR%", 10). // simulamos el limit en 10
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticker", "target_from", "target_to", "company", "action", "brokerage", "rating_from", "rating_to", "time"}).
			AddRow(1, "BSBR", "$4.20", "$4.70", "Banco Santander (Brasil)", "upgraded by", "The Goldman Sachs Group", "Sell", "Neutral", time.Date(2025, 1, 13, 0, 30, 5, 0, time.UTC)).
			AddRow(2, "VYGR", "$11.00", "$9.00", "Voyager Therapeutics", "reiterated by", "Wedbush", "Outperform", "Outperform", time.Date(2025, 1, 14, 0, 30, 5, 0, time.UTC)))

	mock.ExpectQuery("SELECT \\* FROM \"stocks\" WHERE ticker ILIKE \\$1 OR company ILIKE \\$2 OR brokerage ILIKE \\$3 LIMIT \\$4 OFFSET \\$5").
		WithArgs("%BSBR%", "%BSBR%", "%BSBR%", 10, 10). // para la segunda page
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticker", "target_from", "target_to", "company", "action", "brokerage", "rating_from", "rating_to", "time"}).
			AddRow(3, "MSFT", "$300", "$310", "Microsoft", "reiterated by", "Morgan Stanley", "Buy", "Buy", time.Date(2025, 1, 15, 0, 30, 5, 0, time.UTC)))

	return gormDB
}
