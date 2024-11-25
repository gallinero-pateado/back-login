package models

type Usuario_empresa struct {
	Id_empresa               uint   `gorm:"primaryKey;autoIncrement"`
	Firebase_usuario_empresa string `gorm:"type:text;uniqueIndex"`
	Nombre_empresa           string `json:"Nombre_empresa"`
	Correo_empresa           string `json:"Correo_empresa"`
	Sector                   string `json:"Sector"`
	Descripcion              string `json:"Descripcion"`
	Direccion                string `json:"Direccion"`
	Persona_contacto         string `json:"Persona_contacto"`
	Correo_contacto          string `json:"Correo_contacto"`
	Telefono_contacto        int    `json:"Telefono_contacto"`
	Estado_verificacion      uint   `json:"Estado_verificacion"`
	Perfil_Completado        bool   `json:"Perfil_Completado"`
	Rol                      string `json:"Rol"`
}

// TableName establece el nombre de la tabla para GORM
func (Usuario_empresa) TableName() string {
	return "Usuario_empresa"
}
