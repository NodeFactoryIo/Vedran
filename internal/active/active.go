package active

import (
	"fmt"
	"github.com/NodeFactoryIo/vedran/internal/models"
	"github.com/NodeFactoryIo/vedran/internal/repositories"
	log "github.com/sirupsen/logrus"
	"time"
)

const IntervalFromLastPing = 10 * time.Second

func CheckIfNodeActive(node models.Node, repos *repositories.Repos) (bool, error) {
	lastPing, err := repos.PingRepo.FindByNodeID(node.ID)
	if err != nil {
		log.Error(err)
	}

	fmt.Printf("%s NODE: %s\n", node.ID, lastPing.Timestamp.String())

	// more than 10 seconds passed from last ping
	if lastPing.Timestamp.Add(IntervalFromLastPing).Before(time.Now()) {
		// log.Infof("PENALIZE NODE %s", node.ID)
		// actions.PenalizeNode(node, repos.NodeRepo)
		return false, nil
	}

	// node's latest and best block lag behind the best in the pool by more than 10 blocks
	metrics, err := repos.MetricsRepo.FindByID(node.ID)
	if err != nil {
		return false, err
	}
	latestBlockMetrics, err := repos.MetricsRepo.GetLatestBlockMetrics()
	if err != nil {
		return false, err
	}
	if metrics.BestBlockHeight <= latestBlockMetrics.BestBlockHeight-10 &&
		metrics.FinalizedBlockHeight <= latestBlockMetrics.FinalizedBlockHeight-10 {
		// log.Infof("PENALIZE NODE %s", node.ID)
		// actions.PenalizeNode(node, repos.NodeRepo)
		return false, nil
	}
	return true, nil
}