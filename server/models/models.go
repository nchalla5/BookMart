package models

type CredsStruct struct {
	EmailOrPhone string `json:"emailOrPhone"`
	Password     string `json:"password"`
}

type DetailsStruct struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
