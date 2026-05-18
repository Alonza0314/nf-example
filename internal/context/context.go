package context

import (
	"os"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
)

type NFContext struct {
	NfId        string
	Name        string
	UriScheme   models.UriScheme
	BindingIPv4 string
	SBIPort     int
	NetworkData map[string]string
	SpyFamilyData map[string]string
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
	nfContext.NetworkData = map[string]string{
		"Router-Core-01":   "Core Router",
		"Switch-Access-01": "Access Switch",
		"Switch-Dist-01":   "Distribution Switch",
		"AP-Floor-01":      "Access Point",
		"Firewall-01":      "Firewall",
		"gNB-Site-01":      "5G Base Station",
		"UPF-Node-01":      "User Plane Function",
		"AMF-Node-01":      "Access Mobility Function",
	}
}

func GetSelf() *NFContext {
	return &nfContext
}
