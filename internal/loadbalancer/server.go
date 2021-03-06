package loadbalancer

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/NodeFactoryIo/vedran/internal/actions"
	"github.com/NodeFactoryIo/vedran/internal/auth"
	"github.com/NodeFactoryIo/vedran/internal/configuration"
	"github.com/NodeFactoryIo/vedran/internal/controllers"
	"github.com/NodeFactoryIo/vedran/internal/models"
	"github.com/NodeFactoryIo/vedran/internal/prometheus"
	"github.com/NodeFactoryIo/vedran/internal/repositories"
	"github.com/NodeFactoryIo/vedran/internal/router"
	"github.com/NodeFactoryIo/vedran/internal/schedule/checkactive"
	schedulepayout "github.com/NodeFactoryIo/vedran/internal/schedule/payout"
	"github.com/NodeFactoryIo/vedran/internal/schedule/penalize"
	"github.com/asdine/storm/v3"
	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

func StartLoadBalancerServer(
	props configuration.Configuration,
	privateKey string,
) {
	configuration.Config = props

	// set auth secret
	err := auth.SetAuthSecret(props.AuthSecret)
	if err != nil {
		// terminate app: no auth secret provided
		log.Fatalf("Unable to start vedran load balancer: %v", err)
	}

	// init database
	database, err := storm.Open(path.Join(props.RootDir, "vedran-load-balancer.db"))
	if err != nil {
		// terminate app: unable to start database connection
		log.Fatalf("Unable to start vedran load balancer: %v", err)
	}
	log.Debug("Successfully connected to database")

	// initialize repos
	var repos = &repositories.Repos{}
	repos.PingRepo = repositories.NewPingRepo(database)
	repos.MetricsRepo = repositories.NewMetricsRepo(database)
	repos.RecordRepo = repositories.NewRecordRepo(database)
	repos.NodeRepo = repositories.NewNodeRepo(database)
	repos.DowntimeRepo = repositories.NewDowntimeRepo(database)
	repos.PayoutRepo = repositories.NewPayoutRepo(database)
	repos.FeeRepo = repositories.NewFeeRepo(database)
	err = repos.PingRepo.ResetAllPings()
	if err != nil {
		log.Fatalf("Failed reseting pings because of: %v", err)
	}

	// save initial payout if there isn't any saved payouts
	p, err := repos.PayoutRepo.GetAll()
	if err != nil {
		log.Fatalf("Failed creating initial payout because of: %v", err)
	} else if len(*p) == 0 {
		err := repos.PayoutRepo.Save(&models.Payout{
			Timestamp:      time.Now(),
			PaymentDetails: nil,
		})
		if err != nil {
			log.Fatalf("Failed creating initial payout because of: %v", err)
		}
	}

	penalizedNodes, err := repos.NodeRepo.GetPenalizedNodes()
	if err != nil {
		log.Fatalf("Failed fetching penalized nodes because of: %v", err)
	}

	for _, node := range *penalizedNodes {
		go penalize.ScheduleCheckForPenalizedNode(node, *repos)
	}

	// starts task that checks active nodes
	checkactive.StartScheduledTask(repos)

	// start scheduled payout if auto payout enabled
	if props.PayoutConfiguration != nil {
		schedulepayout.StartScheduledPayout(
			*props.PayoutConfiguration,
			privateKey,
			*repos)
	}

	// start server
	log.Infof("Starting vedran load balancer on port :%d...", props.Port)
	apiController := controllers.NewApiController(
		props.WhitelistEnabled, *repos, actions.NewActions(),
	)
	r := router.CreateNewApiRouter(apiController, privateKey)
	prometheus.RecordMetrics(*repos)
	if props.CertFile != "" {
		tlsConfig := &tls.Config{MinVersion: 0}
		server := &http.Server{
			Addr:      fmt.Sprintf(":%d", props.Port),
			Handler:   handlers.CORS()(r),
			TLSConfig: tlsConfig,
		}
		err = server.ListenAndServeTLS(props.CertFile, props.KeyFile)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%d", props.Port), handlers.CORS()(r))
	}
	if err != nil {
		log.Error(err)
	}

	// close database connection
	err = database.Close()
	log.Error(err)
}
