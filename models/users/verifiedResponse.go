package users

// VerifiedUserResponse representa la respuesta con información adicional para usuarios verificados
type VerifiedUserResponse struct {
	Email          string `json:"email"`
	VerifiedStatus string `json:"verified_status"`
	Nama           string `json:"nama,omitempty"`
	IdKelas        *int   `json:"id_kelas,omitempty"`
	Kelas          string `json:"kelas,omitempty"`
}
