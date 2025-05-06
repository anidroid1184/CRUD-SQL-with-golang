package connect

import (
	"database/sql"
	"os"
	"fmt"
	"sync"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"  // añadimos el driver por si se necesita
)

// funcion para conectarse a la BD
var (
	Db *sql.DB
	once sync.Once	// Esto garantiza que se crea una sola conexión a la BD durante la ejecucion del programa
)
func Connecting() {
	// once.Do garantiza que una funcion solo se ejecute unicamente una vez, así sea llamada multiples veces
	once.Do(
		// El metodo once.DO recibe una funcion sin argumento como parametros
		func(){
			// obtener valores de .env, %w sirve para el error string
			if err := godotenv.Load(); err != nil {
				panic(fmt.Errorf("error loading .env file %w", err))
			}

			// se crea un string para simplificar la conexión entre BD y golang
			connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_SERVER"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),	
			)

			var err error

			// intentar conexión
			Db, err = sql.Open("mysql", connectionString)
			if err != nil {
				panic(fmt.Errorf("error opening database: %w", err))
			}

			// verificar conexión
			if err = Db.Ping(); err != nil {
				panic(fmt.Errorf("error connecting to database: %w", err))
			}
			// si no se produjeron errores se establecera la conexion
			fmt.Println("Conexión existosa a la base de datos :D")
		})
	
}

// cerrar conexión - damos el retorno del error
func CloseConnection () error {
	if Db != nil {
		fmt.Printf("\nLa conexión a la BD ha sido cerrada\n")
		return Db.Close()	
	}
	return nil
}