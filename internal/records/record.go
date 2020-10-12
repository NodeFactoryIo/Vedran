package record

import (
	"time"

	"github.com/NodeFactoryIo/vedran/internal/models"
	log "github.com/sirupsen/logrus"
)

// FailedRequest should be called when rpc response is invalid to penalize node.
// It does not return value as it should be called in separate goroutine
func FailedRequest(node models.Node, nodeRepo models.NodeRepository, recordRepo models.RecordRepository) {
	nodeRepo.PenalizeNode(node)

	err := recordRepo.Save(&models.Record{
		NodeId:    node.ID,
		Timestamp: time.Now(),
		Status:    "failed",
	})
	if err != nil {
		log.Errorf("Failed saving failed request because of: %v", err)
	}
}

// SuccessfulRequest should be called when rpc response is valid to reward node.
// It does not return value as it should be called in separate goroutine
func SuccessfulRequest(node models.Node, nodeRepo models.NodeRepository, recordRepo models.RecordRepository) {
	nodeRepo.RewardNode(node)

	err := recordRepo.Save(&models.Record{
		NodeId:    node.ID,
		Timestamp: time.Now(),
		Status:    "successful",
	})
	if err != nil {
		log.Errorf("Failed saving successful request because of: %v", err)
	}
}