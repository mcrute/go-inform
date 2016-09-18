package inform

// TODO: Convert string time to time.Time
// Response packet
type NoopMessage struct {
	Type          string `json:"_type"`
	Interval      int    `json:"interval"`
	ServerTimeUTC string `json:"server_time_in_utc"`
}

type AlarmEntry struct {
	Tag   string `json:"tag"`
	Type  string `json:"string"`
	Value string `json:"val"` // float or int observed
}

type AlarmMessage struct {
	Entries []*AlarmEntry `json:"entries"`
	Index   string        `json:"index"`
	Id      string        `json:"sId"`
	Time    int           `json:"time"`
}

type InterfaceMessage struct {
	IP                 string `json:"ip"`
	MacAddress         string `json:"mac"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	ReceivedBytes      int    `json:"rx_bytes"`
	ReceivedDropped    int    `json:"rx_dropped"`
	ReceivedErrors     int    `json:"rx_errors"`
	ReceivedPackets    int    `json:"rx_packets"`
	TransmittedBytes   int    `json:"tx_bytes"`
	TransmittedDropped int    `json:"tx_dropped"`
	TransmittedErrors  int    `json:"tx_errors"`
	TransmittedPackets int    `json:"tx_packets"`
}

type RadioMessage struct {
	Gain             int    `json:"builtin_ant_gain"`
	BuiltinAntenna   bool   `json:"builtin_antenna"`
	MaxTransmitPower int    `json:"max_txpower"`
	Name             string `json:"name"`
	RadioProfile     string `json:"radio"`
	// "scan_table": []
}

type AccessPointMessage struct {
	BasicSSID               string `json:"bssid"`
	ExtendedSSID            string `json:"essid"`
	ClientConnectionQuality int    `json:"ccq"`
	Channel                 int    `json:"channel"`
	Id                      string `json:"id"`
	Name                    string `json:"name"`
	StationNumber           string `json:"num_sta"` // int?
	RadioProfile            string `json:"radio"`
	Usage                   string `json:"usage"`
	ReceivedBytes           int    `json:"rx_bytes"`
	ReceivedDropped         int    `json:"rx_dropped"`
	ReceivedErrors          int    `json:"rx_errors"`
	ReceivedPackets         int    `json:"rx_packets"`
	ReceivedCrypts          int    `json:"rx_crypts"`
	ReceivedFragments       int    `json:"rx_frags"`
	ReceivedNetworkIDs      int    `json:"rx_nwids"`
	TransmittedBytes        int    `json:"tx_bytes"`
	TransmittedDropped      int    `json:"tx_dropped"`
	TransmittedErrors       int    `json:"tx_errors"`
	TransmittedPackets      int    `json:"tx_packets"`
	TransmitPower           int    `json:"tx_power"`
	TransmitRetries         int    `json:"tx_retries"`
}

// TODO: Convert time to time.Time
type IncomingMessage struct {
	Alarms        []*AlarmMessage       `json:"alarm"`
	ConfigVersion string                `json:"cfgversion"`
	Default       bool                  `json:"default"`
	GuestToken    string                `json:"guest_token"`
	Hostname      string                `json:"hostname"`
	InformURL     string                `json:"inform_url"`
	IP            string                `json:"ip"`
	Isolated      bool                  `json:"isolated"`
	LocalVersion  string                `json:"localversion"`
	Locating      bool                  `json:"locating"`
	MacAddress    string                `json:"mac"`
	IsMfi         string                `json:"mfi"` // boolean as string
	Model         string                `json:"model"`
	ModelDisplay  string                `json:"model_display"`
	PortVersion   string                `json:"portversion"`
	Version       string                `json:"version"`
	Serial        string                `json:"serial"`
	Time          int                   `json:"time"`
	Trackable     string                `json:"trackable"` // boolean as string
	Uplink        string                `json:"uplink"`
	Uptime        int                   `json:"uptime"`
	Interfaces    []*InterfaceMessage   `json:"if_table"`
	Radios        []*RadioMessage       `json:"radio_table"`
	AccessPoints  []*AccessPointMessage `json:"vap_table"`
}
