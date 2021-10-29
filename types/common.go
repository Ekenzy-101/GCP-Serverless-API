package types

type EmailField struct {
	Email string `json:"email" validate:"email,max=255"`
}

type NameField struct {
	Name string `json:"name" validate:"required,name,max=100"`
}

type PasswordField struct {
	Password string `json:"password" validate:"required,min=8,max=128,password"`
}

type LoginRequestBody struct {
	EmailField
	PasswordField
}

type RegisterRequestBody struct {
	NameField
	LoginRequestBody
}

type M map[string]interface{}
