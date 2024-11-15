package reg

type RegeditInfo struct {
	ExePath     string
	InstallDate string
	InstallPath string
	Version     string
}

type AppName = string

const (
	AI                   AppName = "AI"
	Answer                       = "answertool"
	Betterme                     = "betterme"
	Daemon                       = "daemon"
	Infiniti                     = "infiniti"
	Recordplayeroffline          = "recordplayeroffline"
	Recordplayerteaching         = "recordplayerteaching"
	Yunketang                    = "yunketang"
	Yunketangshell               = "yunketangshell"
)
