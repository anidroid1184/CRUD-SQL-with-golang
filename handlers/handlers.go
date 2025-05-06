package handlers

import (
	"bufio"
	"context"
	"database/sql" //libreria necesaria para el sql
	"fmt"
	"os"
	"proyectoMySQL/connect"
	"proyectoMySQL/models"
	"strconv"
	"strings"
	"time"
)

// creamos una structura para los clientes
type ClienteHandler struct{
	// asegurarse de que todas las instancias de la estructura usen la misma conexión ya creada
	db *sql.DB	// todo cliente tendrá un campo db de tipo puntero sql.DB
}

func NuevoClienteHandler() *ClienteHandler{
	connect.Connecting()
	if connect.Db == nil {
		panic("Error: connect.DB no está inicializado")
	}
	return &ClienteHandler{db: connect.Db}	// se accede al puntero para modificar el struct
	// el {db: connect.Db} es equivalente a decir "usa la conexión que ya creamos", de ahí el puntero

}


// se agrega un h *ClienteHandler para modificar el struct original, y actualizar la lista de clientes facilmente

// creamos un metodo para el struct  ClienteHandler
func (h *ClienteHandler) ListarClientes() {
	

	listSQL := "SELECT id, nombre, correo, telefono FROM clientes ORDER BY id ASC"

	rows, err := h.db.Query(listSQL)
	if err != nil{
		panic(fmt.Errorf("error al obtener datos: %w", err))
	}
	defer rows.Close()

	var clientes models.Clientes // se usa para la lista de clientes
	for rows.Next() {
		var cliente models.Cliente
		if err := rows.Scan(
			&cliente.Id,
			&cliente.Nombre,
			&cliente.Correo,
			&cliente.Telefono,
		); err != nil {
			panic(fmt.Errorf("error al escanear fila: %w", err))
		}
		// agregar nuevo cliente
		clientes = append(clientes, cliente)
	}
	printClientes(clientes)
}

// funcion reusable para imprimir clientes
func printClientes(clientes models.Clientes) {
	fmt.Printf("----\nRegistros encontrados: %d\n----", len(clientes))
	for _, cliente := range clientes {
		fmt.Printf(
			"ID: %-4d | Nombre: %-20s | Correo: %-30s | Teléfono: %s\n",
			cliente.Id,
			cliente.Nombre,
			cliente.Correo,
			cliente.Telefono,
		)
	}
}


// pasa a ser un metodo del struct
func (h* ClienteHandler) BuscarPorId(id int) (*models.Cliente, error){
	findSQL := "SELECT id, nombre, correo, telefono FROM clientes WHERE id = ?"
	
	// se usa QueryRow para obtener un único resultado
	row := h.db. QueryRow(findSQL, id)
	
	var cliente models.Cliente

	// escanear el registro
	err := row.Scan(
		&cliente.Id,
		&cliente.Nombre,
		&cliente.Correo,
		&cliente.Telefono,
	)

	// manejo de errores
	if err != nil {
		if err == sql.ErrNoRows {
			// no se encontro registro
			return nil, fmt.Errorf("cliente con ID %d no encontrado", id)
		}
		// otro tipo de error
		return nil, fmt.Errorf("Error buscando cliente: %w", err)
	}
	
	printCliente(cliente)
	return &cliente, nil
	
}

func printCliente(cliente models.Cliente){
	fmt.Printf("----\nRegistro encontrado\n----")
	fmt.Printf(
		"ID: %-4d | Nombre: %-20s | Correo: %-30s | Teléfono: %s\n",
		cliente.Id,
		cliente.Nombre,
		cliente.Correo,
		cliente.Telefono,
	)
	
}

func (h* ClienteHandler) CrearCliente(cliente models.Cliente) (int64, error){
	
	// Validación de datos basica

	if cliente.Nombre == "" || cliente.Correo == "" {
		return 0, fmt.Errorf("Nombre y correo son campos OBLIGATORIOS")
	}

	// Permite la cancelación luego de 5 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// si no se entran datos se cancelará
	defer cancel()

	// query 
	insertSQL := "INSERT INTO clientes(nombre, correo, telefono) VALUES(?, ?, ?)"

	// se usa el ExecContenxt 
	result, err := h.db.ExecContext(
		ctx,
		insertSQL,
		cliente.Nombre,
		cliente.Correo,
		cliente.Telefono,
		)
	
		if err != nil {
			return 0, fmt.Errorf("Error al insertar valores en registro: %w", err)
		}

		// verificar si fue por el timeOUt
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				return 0, fmt.Errorf("La operación tardó demasiado")
			}
		}

		// obtener ID autogenerado
		id, err := result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("error al obtener ID generado: %w", err)
		}
		
	

	fmt.Printf("Se ha creado el cliente con éxito. ID: %d\n", id)

	return id, nil
}

func (h *ClienteHandler) EditarCliente (id int, cliente models.Cliente) (error) {

	// permite la cancelación luego de 5 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// sii no se entran datos
	defer cancel()

	//query
	updateSQL := "UPDATE clientes SET nombre = ?, correo = ?, telefono = ? WHERE id = ?"


	result, err := h.db.ExecContext(
		ctx,
		updateSQL,
		cliente.Nombre,
		cliente.Correo,
		cliente.Telefono,
		id, // importante para denotar donde se hará el cambio
	)
	if err != nil {
		return fmt.Errorf("\nError al realizar la actualización del usuario: %w", err)
	}

	// verificar si hay error por el timeOut
	if err != nil{
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("La operación tardó demasiado")
		}
	}

	_, err = result.RowsAffected()

	fmt.Printf("\nSe ha actualizado con exito el cliente.ID: %d | NOMBRE: %s | CORREO: %s | TELÉFONO: %s\n",
		id,
		cliente. Nombre,
		cliente.Correo,
		cliente.Telefono,
	)

	return nil


}

func (h *ClienteHandler) EliminarCliente (id int) (error) {

	// permite la cancelación luego de 5 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// sii no se entran datos
	defer cancel()

	//query
	updateSQL := "DELETE FROM clientes WHERE id = ?"


	result, err := h.db.ExecContext(
		ctx,
		updateSQL,
		id, // importante para denotar donde se hará el cambio
	)

	if err != nil {
		return fmt.Errorf("\nError al realizar al eliminar el usuario: %w", err)
	}

	// verificar si hay error por el timeOut
	if err != nil{
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("La operación tardó demasiado")
		}
	}

	_, err = result.RowsAffected()

	fmt.Printf("\nSe ha eliminado con exito  con exito el cliente. ID: %d\n", id)

	return nil


}

func (h* ClienteHandler) Menu() {
	scanner := bufio.NewScanner(os.Stdin)
	var opcion string
	limpiarPantalla()

	for {
		mostrarMenu()

		if !scanner.Scan(){
			break // si hay errores de lectura abortar
		}

		opcion = scanner.Text()

		switch opcion {
		case "1":
			limpiarPantalla()
			h.ListarClientes()
			pausa(scanner)

		case "2":
			limpiarPantalla()
			id := leerID(scanner)
			h.BuscarPorId(id)
			pausa(scanner)
		
		case "3":
			limpiarPantalla()
			h.crearClienteInterno(scanner)
			pausa(scanner)
		
		case "4":
			limpiarPantalla()
			h.editarClienteInterno(scanner)
			pausa(scanner)
		
		case "5":
			limpiarPantalla()
			h.eliminarClienteInterno(scanner)
			pausa(scanner)
		
		case "6":
			fmt.Println("\nSaliento del programa.")
			// al retornar en vacio salimos del programa
			return

		default:
			fmt.Println("\nOpción no válida. Intente nuevamente")
			time.Sleep(1 * time.Second)

		}

		limpiarPantalla()

	}
}

// funciones auxiliares para el menu
func mostrarMenu() {
	fmt.Println(`
	*** MENÚ PRINCIPAL ***
	1) Listar clientes
	2) Buscar cliente por ID
	3) Crear cliente
	4) Editar cliente
	5) Eliminar cliente
	6) Salir
	**********************`)
    fmt.Print("Seleccione una opción: ")
}

func limpiarPantalla(){
	fmt.Print("\033[H\033[2J")
}
func pausa(scanner *bufio.Scanner){
	fmt.Print("\nPresione Enter para continuar...")
	scanner.Scan()
}

func leerID(scanner *bufio.Scanner) int {
	for {
		fmt.Print("Ingrese el ID del cliente: ")
		scanner.Scan()
		idStr := scanner.Text()

		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			fmt.Println("ID inválido. Debe ser mayor a 0")
			continue
		}
		return id
	}
}

func (h *ClienteHandler) crearClienteInterno(scanner *bufio.Scanner){
	fmt.Println("*** NUEVO CLIENTE ***")

	cliente := models.Cliente{
		Nombre: leerEntrada(scanner, "Nombre"),
		Correo: leerEntrada(scanner, "Correo electrónico"),
		Telefono: leerEntrada(scanner, "Teléfono"),
	}

	if _, err := h.CrearCliente(cliente); err != nil {
		fmt.Printf("\nError creando cliente: %v\n", err)
	}
}

func (h *ClienteHandler) editarClienteInterno(scanner *bufio.Scanner) {
	fmt.Println("*** EDITAR CLIENTE ***")
	id := leerID(scanner)

	// buscar cliente existente
	clienteExistente, err := h.BuscarPorId(id)
	if err != nil{
		fmt.Printf("\nError: %v\n", err)
		return
	}

	// Pedir nuevos datos
	nuevoCliente := models.Cliente{
		Nombre: leerEntradaOpcional(scanner, "Nombre", clienteExistente.Nombre),
		Correo: leerEntradaOpcional(scanner, "Correo", clienteExistente.Correo),
		Telefono: leerEntradaOpcional(scanner, "Telefono", clienteExistente.Telefono),
	}
	if err := h.EditarCliente(id, nuevoCliente); err != nil {
		fmt.Printf("\nError actualizando cliente: %v\n", err)
	}
}

func (h *ClienteHandler) eliminarClienteInterno(scanner *bufio.Scanner) {
	fmt.Println("*** ELIMINAR CLIENTE ***")
	id := leerID(scanner)

	// confirmamos la solicitud
	fmt.Printf("¿Está seguro de que desea eliminar al cliente %d? (s/n)", id)

	scanner.Scan()
	if strings.ToLower(scanner.Text()) != "s" {
		fmt.Println("Operación cancelada")
		return 
	}

	if err := h.EliminarCliente(id); err != nil {
		fmt.Printf("\nError eliminando cliente: %v\n", err)
	}
}

func leerEntrada(scanner *bufio.Scanner, campo string) string {
	for {
		fmt.Printf("%s: ", campo)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			return input
		}
		fmt.Println("Este campo es obligatorio")
	}
}

func leerEntradaOpcional(scanner *bufio.Scanner, campo, valorActual string) string {
	fmt.Printf("%s [Actual: %s] (Enter para mantener valor actual): ", campo, valorActual)
	scanner.Scan()
	// sirve para quitar espacios en blanco
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return valorActual
	}
	return input
}

