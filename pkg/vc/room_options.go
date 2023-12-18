package vc

type RoomOptions struct {
	Name                string
	ProgramInstanceId   string
	ProgramLibraryId    int
	Notes               string
	Level               string
	Location            string
	TimeZone            string
	Latitude            string
	Longitude           string
	AddressSetsLocation bool
	UserFile            string
}

func NewRoomOptions(programId int, id string, name string) RoomOptions {
	return RoomOptions{
		Name:              name,
		ProgramLibraryId:  programId,
		ProgramInstanceId: id,
	}
}
