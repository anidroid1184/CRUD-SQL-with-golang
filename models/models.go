package models

// registro individual
type Cliente struct{
	Id int
	Nombre string
	Correo string
	Telefono string
}
// multiples registros
type Clientes []Cliente