package alertstate

const (
	InputCacheLenDef = 1000
)

var discardCount int = 0
var discardBigCount int = 0
var discardSmallCount int = 0

type classType int32

const (
	AlertTypeBusTransport  classType = 1
	AlertTypeBusPdxp       classType = 2
	AlertTypeIpFragLoss    classType = 3
	AlertTypePdxpLoss      classType = 4
	AlertTypePkgErr        classType = 5
	AlertTypeBandLinkGroup classType = 6
	AlertTypeBandSiteSi    classType = 7
	AlertTypeBandPdxpApp   classType = 8
	AlertTypeFlow5Tuple    classType = 9
)

var classStrInt = map[string]classType{
	"bus_err_transport": AlertTypeBusTransport,
	"bus_err_pdxp":      AlertTypeBusPdxp,
	"pkg_loss_ip":       AlertTypeIpFragLoss,
	"pkg_loss_pdxp":     AlertTypePdxpLoss,
	"pkg_err":           AlertTypePkgErr,
	"flow_line_group":   AlertTypeBandLinkGroup,
	"flow_site_si":      AlertTypeBandSiteSi,
	"flow_pdxp_app":     AlertTypeBandPdxpApp,
	"flow_5tuple":       AlertTypeFlow5Tuple,
}

var classIntStr = map[classType]string{
	AlertTypeBusTransport:  "bus_err_transport",
	AlertTypeBusPdxp:       "bus_err_pdxp",
	AlertTypeIpFragLoss:    "pkg_loss_ip",
	AlertTypePdxpLoss:      "pkg_loss_pdxp",
	AlertTypePkgErr:        "pkg_err",
	AlertTypeBandLinkGroup: "flow_line_group",
	AlertTypeBandSiteSi:    "flow_site_si",
	AlertTypeBandPdxpApp:   "flow_pdxp_app",
	AlertTypeFlow5Tuple:    "flow_5tuple",
}

func Init(width int32) {
	//WindowWidth := 5
	gGlobalCache.Init()
	gIdNameMap.Init()
	GLocalCach.Init(int32(width), 2)
}
