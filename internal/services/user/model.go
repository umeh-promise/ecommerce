package user

type User struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	PhoneNumber    string `json:"phone_number"`
	DOB            string `json:"dob"`
	Gender         string `json:"gender"`
	ProfilePicture string `json:"profile_picture"`
	Version        string `json:"-"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
