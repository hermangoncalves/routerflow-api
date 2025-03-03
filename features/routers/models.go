package routers

// Router represents a router stored in the database
type Router struct {
	ID       int
	IP       string
	Username string
	Password string
}

type RouterRegisterRequest struct {
	IP       string `json:"ip" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RouterStatus represents the structured system status response
type RouterStatus struct {
	Uptime               string `json:"uptime"`
	Version              string `json:"version"`
	BuildTime            string `json:"build_time"`
	FactorySoftware      string `json:"factory_software"`
	FreeMemory           string `json:"free_memory"`
	TotalMemory          string `json:"total_memory"`
	CPU                  string `json:"cpu"`
	CPUCount             string `json:"cpu_count"`
	CPUFrequency         string `json:"cpu_frequency"`
	CPULoad              string `json:"cpu_load"`
	FreeHDDSpace         string `json:"free_hdd_space"`
	TotalHDDSpace        string `json:"total_hdd_space"`
	WriteSectSinceReboot string `json:"write_sect_since_reboot"`
	WriteSectTotal       string `json:"write_sect_total"`
	ArchitectureName     string `json:"architecture_name"`
	BoardName            string `json:"board_name"`
	Platform             string `json:"platform"`
}
