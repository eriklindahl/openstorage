package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/libopenstorage/gossip/types"
	"github.com/libopenstorage/openstorage/api"
	clusterclient "github.com/libopenstorage/openstorage/api/client/cluster"
	"github.com/libopenstorage/openstorage/cluster"
	"github.com/stretchr/testify/assert"
)

func TestClusterEnumerateSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		Enumerate().
		Return(api.Cluster{
			Id:            "cluster-dummy-id",
			Status:        api.Status_STATUS_OK,
			ManagementURL: "mgmturl:1234/mgmt-endpoint",
			Nodes: []api.Node{
				api.Node{
					Hostname: "node1-hostname",
					Id:       "1",
				},
				api.Node{
					Hostname: "node2-hostname",
					Id:       "2",
				},
				api.Node{
					Hostname: "node3-hostname",
					Id:       "3",
				},
			},
		}, nil)
	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp, err := restClient.Enumerate()

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.EqualValues(t, "cluster-dummy-id", resp.Id)
}

func TestInspectNodeSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	nodeID := "dummy-node-id-121"
	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		Inspect(nodeID).
		Return(api.Node{
			Id:       nodeID,
			Hostname: "dummy-hostname",
			Status:   api.Status_STATUS_OK,
		}, nil)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp, err := restClient.Inspect(nodeID)

	assert.NoError(t, err)
	assert.EqualValues(t, nodeID, resp.Id)
	assert.EqualValues(t, api.Status_STATUS_OK, resp.Status)

	// mock the cluster response with IP
	nodeIP := "192.168.1.1"

	tc.MockCluster().
		EXPECT().
		Inspect(nodeIP).
		Return(api.Node{
			Id:       nodeID,
			Hostname: "dummy-hostname",
			Status:   api.Status_STATUS_OK,
		}, nil)

	// make the REST call
	restClient = clusterclient.ClusterManager(c)
	resp, err = restClient.Inspect(nodeIP)

	assert.NoError(t, err)
	assert.EqualValues(t, nodeID, resp.Id)
	assert.EqualValues(t, api.Status_STATUS_OK, resp.Status)
}

func TestGossipStateSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		GetGossipState().
		Return(&cluster.ClusterState{
			NodeStatus: []types.NodeValue{
				{
					GenNumber: uint64(1234),
					Id:        "node1-id",
					Status:    types.NODE_STATUS_UP,
				},
				{
					GenNumber: uint64(4567),
					Id:        "node2-id",
					Status:    types.NODE_STATUS_UP,
				},
				{
					GenNumber: uint64(7890),
					Id:        "node3-id",
					Status:    types.NODE_STATUS_UP,
				},
			},
		})

		// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.GetGossipState()

	assert.NotNil(t, resp)

	assert.Len(t, resp.NodeStatus, 3)
	assert.EqualValues(t, "node1-id", resp.NodeStatus[0].Id)
}

func TestGossipStateFailed(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		GetGossipState().
		Return(&cluster.ClusterState{})

		// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.GetGossipState()

	assert.NotNil(t, resp)

	assert.Len(t, resp.NodeStatus, 0)

}

func TestPeerStatusSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	listenerName := "pxd"
	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		PeerStatus(listenerName).
		Return(map[string]api.Status{
			"node-1": api.Status_STATUS_OK,
			"node-2": api.Status_STATUS_OK,
		}, nil)

		// make the REST call
	restClient := clusterclient.ClusterManager(c)

	statusMap, err := restClient.PeerStatus(listenerName)
	assert.NoError(t, err)
	assert.Equal(t, len(statusMap), 2)

	for _, v := range statusMap {
		assert.Equal(t, v, api.Status_STATUS_OK)
	}
}

func TestNodeHealthSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		NodeStatus().
		Return(api.Status_STATUS_OK, nil)

		// make the REST call
	restClient := clusterclient.ClusterManager(c)

	status, err := restClient.NodeStatus()
	assert.NoError(t, err)
	assert.Equal(t, api.Status_STATUS_OK, status)

}
func TestClusterNodeStatusSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	restClient, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// Set expections
	tc.MockCluster().
		EXPECT().
		NodeStatus().
		Return(api.Status_STATUS_OK, nil).
		Times(1)

	// Check status
	status, err := clusterclient.ClusterManager(restClient).NodeStatus()
	assert.NoError(t, err)
	assert.Equal(t, api.Status_STATUS_OK, status)
}

func TestNodeRemoveSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	nodeId := "dummy-node-id-121"
	secondNodeId := "dummy-node-id-131"

	nodes := []api.Node{
		{Id: nodeId},
		{Id: secondNodeId},
	}

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		Remove(nodes, false).
		Return(nil)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.Remove(nodes, false)

	assert.NoError(t, resp)
}

func TestNodeRemoveFailed(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	nodeId := ""

	nodes := []api.Node{
		{Id: nodeId},
	}

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		Remove(nodes, false).
		Return(fmt.Errorf("error in removing node"))

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.Remove(nodes, false)

	assert.Error(t, resp)

	assert.Contains(t, resp.Error(), "error in removing node")

}

func TestEnableGossipSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		EnableUpdates().
		Return(nil)

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.EnableUpdates()

	assert.NoError(t, resp)

}

func TestDisableGossipSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		DisableUpdates().
		Return(nil)

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.DisableUpdates()

	assert.NoError(t, resp)

}
func TestEnumerateAlertsSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// time frame is exactly 24 hrs from current time.
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		EnumerateAlerts(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&api.Alerts{
			Alert: []*api.Alert{
				&api.Alert{
					AlertType: 1,
					Id:        123,
					Resource:  api.ResourceType_RESOURCE_TYPE_NODE,
				},
			},
		}, nil)

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp, err := restClient.EnumerateAlerts(startTime, endTime, api.ResourceType_RESOURCE_TYPE_NODE)

	assert.NoError(t, err)

	assert.Len(t, resp.Alert, 1)
	assert.EqualValues(t, 123, resp.Alert[0].GetId())
	assert.EqualValues(t, api.ResourceType_RESOURCE_TYPE_NODE, resp.Alert[0].GetResource())
}

func TestClearAlertSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// alertId
	alertID := int64(12345)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		ClearAlert(api.ResourceType_RESOURCE_TYPE_NODE, gomock.Any()).
		Return(nil)

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.ClearAlert(api.ResourceType_RESOURCE_TYPE_NODE, alertID)
	assert.NoError(t, resp)
}

func TestClearAlertFailed(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// alertId
	alertID := int64(12345)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		ClearAlert(api.ResourceType_RESOURCE_TYPE_NODE, gomock.Any()).
		Return(fmt.Errorf("Error in clearing alert"))

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.ClearAlert(api.ResourceType_RESOURCE_TYPE_NODE, alertID)
	assert.Error(t, resp)
	assert.Contains(t, resp.Error(), "Error in clearing alert")
}

func TestEraseAlertSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// alertId
	alertID := int64(12345)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		EraseAlert(gomock.Any(), gomock.Any()).
		Return(nil)

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.EraseAlert(api.ResourceType_RESOURCE_TYPE_NODE, alertID)
	assert.NoError(t, resp)
}

func TestEraseAlertFailed(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// alertId
	alertID := int64(12345)

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		EraseAlert(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("Error in Erasing alert"))

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	resp := restClient.EraseAlert(api.ResourceType_RESOURCE_TYPE_NODE, alertID)
	assert.Error(t, resp)
	assert.Contains(t, resp.Error(), "Error in Erasing alert")
}

func TestGetNodeIdFromIpSuccess(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	nodeIP := "192.168.1.10"
	nodeID := "dummy-node-id-ip"

	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		GetNodeIdFromIp(nodeIP).
		Return(nodeID, nil)

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	id, err := restClient.GetNodeIdFromIp(nodeIP)

	assert.NoError(t, err)
	assert.EqualValues(t, nodeID, id)
}

func TestGetNodeIdFromIpFailed(t *testing.T) {

	// Create a new global test cluster
	ts, tc := testClusterServer(t)
	defer ts.Close()
	defer tc.Finish()

	// create a cluster client to make the REST call
	c, err := clusterclient.NewClusterClient(ts.URL, "v1")
	assert.NoError(t, err)

	nodeIP := "192.168.1.10"
	// mock the cluster response
	tc.MockCluster().
		EXPECT().
		GetNodeIdFromIp(nodeIP).
		Return(nodeIP, fmt.Errorf("Failed to locate IP in this cluster."))

	// make the REST call
	restClient := clusterclient.ClusterManager(c)
	id, err := restClient.GetNodeIdFromIp(nodeIP)

	assert.EqualValues(t, nodeIP, id)
	// Client code fix should be able to error respons: todo
	assert.Contains(t, err.Error(), "error 500")
}
