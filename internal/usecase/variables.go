package usecase

const (
	UpdateInfoQuery     = "info"
	UpdatePasswordQuery = "password"
	UpdateSessionQuery  = "session"
	UniqueEmailErr      = "UNIQUE constraint failed: users.email"
	UniqueNameErr       = "UNIQUE constraint failed: users.name"
	DateAndTimeFormat   = "2006-01-02 15:04:05"
	DateFormat          = "2006-01-02"
	UserGenderMale      = "Male"
	UserGenderFemale    = "Female"
)
