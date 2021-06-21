package app

import (
	"fmt"

	"github.com/cosmostation/cosmostation-cosmos/client"
	"github.com/cosmostation/cosmostation-cosmos/custom"
	"github.com/cosmostation/cosmostation-cosmos/db"
	mdschema "github.com/cosmostation/mintscan-database/schema"
	"go.uber.org/zap"

	mblconfig "github.com/cosmostation/mintscan-backend-library/config"

	mp "github.com/cosmostation/mintscan-prometheus/prometheus"
)

type App struct {
	Config          *mblconfig.Config
	Client          *client.Client
	DB              *db.Database
	RawDB           *db.RawDatabase
	ExporterMetrics mp.ExporterMetrics
	ChainNumMap     map[int]string
	ChainIDMap      map[string]int
	MessageIDMap    map[int]string
	MessageTypeMap  map[string]int
}

func init() {
	if !custom.IsSetAppConfig() {
		panic(fmt.Errorf("appconfig was not set"))
	}

	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)
	defer l.Sync()
}

// NewApp
func NewApp(fileBaseName string) *App {
	app := new(App)
	app.Config = mblconfig.ParseConfig(fileBaseName)

	app.Client = client.NewClient(&app.Config.Client)

	app.DB = db.Connect(&app.Config.DB)
	err := app.DB.Ping()
	if err != nil {
		panic(err)
	}

	if fileBaseName == "chain-exporter" {
		app.RawDB = db.RawDBConnect(&app.Config.RAWDB)
		err = app.RawDB.Ping()
		if err != nil {
			panic(err)
		}
		app.ExporterMetrics = mp.NewMetricsForExporter(app.Config.Prometheus.Namespace)
		mp.RegisterMetricForExporter(&app.ExporterMetrics)
		go mp.StartMetricsScraping(app.Config.Prometheus.Path, app.Config.Prometheus.Port)
	}
	mdschema.SetCommonSchema(app.Config.DB.CommonSchema)
	mdschema.SetChainSchema(app.Config.DB.ChainSchema)

	// app.DB.AddQueryHook(dbLogger{})    // debugging 용
	// app.RawDB.AddQueryHook(dbLogger{}) // debugging 용

	return app
}

// SetChainID ChainID를 할당하고, DB에서 InsertSelect()하여 맵을 구성
func (a *App) SetChainID() {
	a.ChainIDMap = make(map[string]int)
	a.ChainNumMap = make(map[int]string)
	if a.Config.Chain.ChainID == "" {
		chainID, err := a.Client.RPC.GetNetworkChainID()
		if err != nil {
			panic(err)
		}
		a.Config.Chain.ChainID = chainID
	}

	exist, err := a.DB.ExistChainID(a.Config.Chain.ChainID)
	if err != nil {
		panic(err)
	}

	if !exist {
		// insert db
		if err := a.DB.InsertChainID(a.Config.Chain.ChainID); err != nil {
			panic(err)
		}
	}

	chainInfo, err := a.DB.GetChainInfo()
	if err != nil {
		panic(err)
	}

	for _, c := range chainInfo {
		a.ChainNumMap[int(c.ID)] = c.ChainID
		a.ChainIDMap[c.ChainID] = int(c.ID)
	}

	fmt.Println("ChainIDMap :", a.ChainIDMap)
	fmt.Println("ChainNumMap :", a.ChainNumMap)
}

func (a *App) SetMessageInfo() {
	a.MessageIDMap = make(map[int]string)
	a.MessageTypeMap = make(map[string]int)
	messageInfo, err := a.DB.GetMessageInfo()
	if err != nil {
		panic(err)
	}
	for _, m := range messageInfo {
		a.MessageIDMap[int(m.ID)] = m.Type
		a.MessageTypeMap[m.Type] = int(m.ID)
	}
}
