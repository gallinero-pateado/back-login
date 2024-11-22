package models

type Usuario struct {
	Id                uint   `gorm:"primaryKey;autoIncrement"`
	Firebase_usuario  string `gorm:"type:text;uniqueIndex"`
	Correo            string `json:"Correo"`
	Nombres           string `json:"Nombres"`
	Apellidos         string `json:"Apellidos"`
	Fecha_nacimiento  string `json:"Fecha_Nacimiento"`
	Ano_ingreso       string `json:"Ano_Ingreso"`
	Id_carrera        uint   `json:"Id_carrera"`
	Id_estado_usuario bool   `json:"Id_Estado_Usuario"`
	Foto_perfil       string `json:"Foto_Perfil"`
	Rol               string `json:"Rol"`
	PerfilCompletado  bool   `json:"PerfilCompletado"`
}

// TableName establece el nombre de la tabla para GORM
func (Usuario) TableName() string {
	return "Usuario"
}
