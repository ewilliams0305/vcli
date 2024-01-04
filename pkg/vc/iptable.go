package vc

const (
	IPTABLEBYID = "IpTableByPID"
)

type VcIpTableApi interface {
	GetIpTable(roomId string) ([]IpTableEntry, VirtualControlError)
}

func (v *VC) GetIpTable(roomId string) ([]IpTableEntry, VirtualControlError) {
	return getIpTable(v, roomId)
}

func getIpTable(server *VC, roomId string) ([]IpTableEntry, VirtualControlError) {

	var results IpTableResponse
	err := server.getBody(IPTABLEBYID+"/"+roomId, &results)

	if err != nil {
		return make([]IpTableEntry, 0), NewServerError(500, err)
	}
	return results.IpTableDevice.IpTablePrograms.IPTableByPID, nil
}

type IpTableResponse struct {
	IpTableDevice `json:"Device"`
}

type IpTableDevice struct {
	IpTablePrograms `json:"Programs"`
}

type IpTablePrograms struct {
	IPTableByPID []IpTableEntry `json:"IpTableByPID"`
}

type IpTableEntry struct {
	UniqueID           int    `json:"UniqueId"`
	ProgramInstanceID  string `json:"ProgramInstanceId"`
	ProgramIPID        int    `json:"ProgramIpId"`
	Model              string `json:"Model"`
	Description        string `json:"Description"`
	RemoteIP           string `json:"remote_ip"`
	Status             string `json:"Status"`
	DeviceType         int    `json:"device_type"`
	MacAddress         string `json:"MacAddress"`
	DeviceID           int    `json:"DeviceId"`
	Hostname           string `json:"Hostname"`
	SupportAssociation bool   `json:"SupportAssociation"`
}
