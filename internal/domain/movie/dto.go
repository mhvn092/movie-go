package movie

type MovieGetAllResponse struct {
	Id             int    `json:"id"              db:"id"`
	Title          string `json:"title"           db:"title"`
	ProductionYear int    `json:"production_year" db:"production_year"`
}

type MovieGetDetailResponse struct {
	Id             int                      `json:"id"              db:"id"`
	Title          string                   `json:"title"           db:"title"`
	ProductionYear int                      `json:"production_year" db:"production_year"`
	DirectorId     int                      `json:"director_id"     db:"director_id"`
	DirectorName   string                   `json:"director_name"   db:"director_name"`
	GenreId        int                      `json:"genre_id"        db:"genre_id"`
	GenreName      string                   `json:"genre_name"      db:"genre_name"`
	Description    string                   `json:"description"     db:"description"`
	Staffs         []MovieStaffBaseResponse `json:"movie_staffs"    db:"movie_staffs"`
}

type MovieStaffBaseResponse struct {
	StaffId        int    `json:"staff_id"         db:"staff_id"`
	StaffName      string `json:"staff_name"       db:"staff_name"`
	StaffTypeId    int    `json:"staff_type_id"    db:"staff_type_id"`
	StaffTypeTitle string `json:"staff_type_title" db:"staff_type_title"`
}

type movieStaffGetDetailBaseRow struct {
	StaffId        *int    `db:"staff_id"`
	StaffName      *string `db:"staff_name"`
	StaffTypeId    *int    `db:"staff_type_id"`
	StaffTypeTitle *string `db:"staff_type_title"`
}

type MovieUpsertPayload struct {
	Title          string                        `json:"title"           db:"title"           validate:"required,is_string"`
	ProductionYear int                           `json:"production_year" db:"production_year" validate:"required,is_int,is_valid_year"`
	DirectorId     int                           `json:"director_id"     db:"director_id"     validate:"required,is_int"`
	GenreId        int                           `json:"genre_id"        db:"genre_id"        validate:"required,is_int"`
	Description    string                        `json:"description"     db:"description"     validate:"required,is_string"`
	Staffs         []movieStaffUpsertBasePayload `json:"movie_staffs"    db:"movie_staffs"    validate:"required"`
}

type movieStaffUpsertBasePayload struct {
	StaffId     int `json:"staff_id"      db:"staff_id"      validate:"required, is_int"`
	StaffTypeId int `json:"staff_type_id" db:"staff_type_id" validate:"required, is_int"`
}
