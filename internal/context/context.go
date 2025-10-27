package context

import (
	"os"
	"sync"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type NFContext struct {
	NfId        string
	Name        string
	UriScheme   models.UriScheme
	BindingIPv4 string
	SBIPort     int

	SpyFamilyData map[string]string

	MessageRecord []string
	MessageMu     sync.Mutex

	Tasks      []Task
	TaskMutex  sync.RWMutex
	NextTaskID uint64

	Messages       []Message
	DragonBallData map[string]int32

	Fortunes     []string
	FortuneMutex sync.RWMutex

	AttendanceData []string

	TimeZoneData map[string]string
}

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Time    string `json:"time"`
}

var nfContext = NFContext{}

func InitNfContext() {
	cfg := factory.NfConfig

	nfContext.NfId = uuid.New().String()
	nfContext.Name = "ANYA"

	nfContext.UriScheme = cfg.Configuration.Sbi.Scheme
	nfContext.SBIPort = cfg.Configuration.Sbi.Port
	nfContext.BindingIPv4 = os.Getenv(cfg.Configuration.Sbi.BindingIPv4)
	if nfContext.BindingIPv4 != "" {
		logger.CtxLog.Info("Parsing ServerIPv4 address from ENV Variable.")
	} else {
		nfContext.BindingIPv4 = cfg.Configuration.Sbi.BindingIPv4
		if nfContext.BindingIPv4 == "" {
			logger.CtxLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
			nfContext.BindingIPv4 = "0.0.0.0"
		}
	}
	nfContext.SpyFamilyData = map[string]string{
		"Loid":   "Forger",
		"Anya":   "Forger",
		"Yor":    "Forger",
		"Bond":   "Forger",
		"Becky":  "Blackbell",
		"Damian": "Desmond",
		"Franky": "Franklin",
		"Fiona":  "Frost",
		"Sylvia": "Sherwood",
		"Yuri":   "Briar",
		"Millie": "Manis",
		"Ewen":   "Egeburg",
		"Emile":  "Elman",
		"Henry":  "Henderson",
		"Martha": "Marriott",
	}
	nfContext.AttendanceData = []string{}

	nfContext.MessageRecord = []string{}

	nfContext.Tasks = make([]Task, 0)
	nfContext.NextTaskID = 0

	nfContext.Messages = make([]Message, 0)

	nfContext.DragonBallData = map[string]int32{
		"Goku":    7,
		"Vegeta":  6,
		"Gohan":   5,
		"Trunks":  4,
		"Piccolo": 3,
		"Krillin": 2,
		"Yamcha":  1,
	}

	nfContext.Fortunes = []string{
		"大吉: All your endeavors will be successful.",
		"中吉: You will have good luck, but be cautious.",
		"小吉: A small amount of luck is coming your way.",
		"吉: Good fortune is with you.",
		"末吉: Your luck is gradually improving.",
		"凶: Be careful, misfortune may be ahead.",
		"大凶: A great misfortune is coming. Be prepared.",
	}

	nfContext.TimeZoneData = map[string]string{
		"Taipei":  "UTC+8",
		"Tokyo":   "UTC+9",
		"Seoul":   "UTC+9",
		"NewYork": "UTC-5",
		"Paris":   "UTC+2",
		"London":  "UTC+1",
		"Berlin":  "UTC+2",
		"Sydney":  "UTC+10",
		"Moscow":  "UTC+3",
		"Dubai":   "UTC+4",
	}
}

func GetSelf() *NFContext {
	return &nfContext
}
