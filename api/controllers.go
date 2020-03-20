package api

import (
	"encoding/json"
	"fmt"
	"github.com/heeeeeng/node_stewpot/stewpot"
	"github.com/heeeeeng/node_stewpot/types"
	"html/template"
	"net/http"
	"strconv"
)

type StewpotController struct {
	stewpot *stewpot.Stewpot
}

func NewSewpotController(s *stewpot.Stewpot) *StewpotController {
	return &StewpotController{stewpot: s}
}

func (c *StewpotController) MainPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/index.html")
	t.Execute(w, "hello!")
}

func (c *StewpotController) SimPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/multi_simulate.html")
	t.Execute(w, "hello!")

}

func (c *StewpotController) Static(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./"))
	fs.ServeHTTP(w, r)
}

func (c *StewpotController) Restart(w http.ResponseWriter, r *http.Request) {
	nodeNum, err := urlParamToInt(r, "node_num")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("can't get node_num from params: %v", err))
		return
	}
	maxIn, err := urlParamToInt(r, "max_in")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("can't get max_in from params: %v", err))
		return
	}
	maxOut, err := urlParamToInt(r, "max_out")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("can't get max_out from params: %v", err))
		return
	}
	bandwidth, err := urlParamToInt(r, "bandwidth")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("can't get bandwidth from params: %v", err))
		return
	}
	maxBest := 4

	c.stewpot.RestartNetwork(nodeNum, maxIn, maxOut, maxBest, int64(bandwidth), nil)
	Response(w, nil)
}

func (c *StewpotController) GetNetworkGraph(w http.ResponseWriter, r *http.Request) {
	Response(w, c.stewpot.MarshalNodes())
}

type msgReqBody struct {
	NodeID     string
	Size       int64
	Difficulty int64
	Content    string
}

func (c *StewpotController) SendMsg(w http.ResponseWriter, r *http.Request) {
	var reqBody msgReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}

	// check variables
	if reqBody.Difficulty <= 0 {
		ErrorResponse(w, "difficulty invalid, should larger than 0")
		return
	}
	if reqBody.Size <= 0 {
		ErrorResponse(w, "size invalid, should larger than 0")
		return
	}

	nodes := c.stewpot.Nodes()
	node := nodes[reqBody.NodeID]
	msg := c.stewpot.GenerateMsg(reqBody.Difficulty, reqBody.Size, reqBody.Content)
	timestamp := c.stewpot.SendMsg(node, msg)
	fmt.Println(timestamp)
	Response(w, []byte(strconv.FormatInt(timestamp, 10)))
}

func (c *StewpotController) AllNodeSendMsg(w http.ResponseWriter, r *http.Request) {
	var reqBody msgReqBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}

	// check variables
	if reqBody.Difficulty <= 0 {
		ErrorResponse(w, "difficulty invalid, should larger than 0")
		return
	}
	if reqBody.Size <= 0 {
		ErrorResponse(w, "size invalid, should larger than 0")
		return
	}

	for _, node := range c.stewpot.Nodes() {
		msg := c.stewpot.GenerateMsg(reqBody.Difficulty, reqBody.Size, reqBody.Content)
		c.stewpot.SendMsg(node, msg)
	}
	Response(w, "success")
}

type Tasks []types.Task

func (c *StewpotController) GetTimeUnit(w http.ResponseWriter, r *http.Request) {
	timestamps := r.URL.Query()["time"]
	if len(timestamps) == 0 {
		ErrorResponse(w, "time is missing")
		return
	}
	timestamp, err := strconv.Atoi(timestamps[0])
	if err != nil {
		ErrorResponse(w, "time invalid, must be integer")
		return
	}
	//fmt.Println("GetTimeUnit request: ", timestamp)
	tasks := c.stewpot.GetTimeUnitTasks(int64(timestamp))
	if tasks == nil {
		ErrorResponse(w, fmt.Sprintf("can't find timestamp at time: %d", timestamp))
		return
	}
	data, _ := json.Marshal(tasks)
	Response(w, data)
}

type MultiSimResp struct {
	Xs []string `json:"xs"`
	Ys []int64  `json:"ys"`
}

func (c *StewpotController) MultiSimulate(w http.ResponseWriter, r *http.Request) {
	// load data
	iterNum, err := urlParamToInt(r, "iter_num")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("iter_num param error: %v", err))
		return
	}
	nodeNum, err := urlParamToInt(r, "node_num")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("node_num param error: %v", err))
		return
	}
	maxIn, err := urlParamToInt(r, "max_in")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("max_in param error: %v", err))
		return
	}
	maxOut, err := urlParamToInt(r, "max_out")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("max_out param error: %v", err))
		return
	}
	bandwidth, err := urlParamToInt(r, "bandwidth")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("bandwidth param error: %v", err))
		return
	}
	maxMsgSize, err := urlParamToInt(r, "max_msg_size")
	if err != nil {
		ErrorResponse(w, fmt.Sprintf("max_msg_size param error: %v", err))
		return
	}

	// init config
	conf := stewpot.SimConfig{
		IterNum:   iterNum,
		NodeNum:   nodeNum,
		Bandwidth: int64(bandwidth),
		MaxIn:     maxIn,
		MaxOut:    maxOut,
	}

	// simulate loop
	var xs []string
	var ys []int64
	// TODO min msg size and enlarge number per iter should not be hard coded.
	minMsgSize := 256 * types.SizeByte
	sizeGap := 16 * types.SizeKB
	for i := minMsgSize; i < int64(maxMsgSize); i += sizeGap {
		conf.MsgSize = i
		avgTimeUsage := c.stewpot.MultiSimulate(conf)

		x := strconv.Itoa(int(i / types.SizeByte))
		xs = append(xs, x)

		ys = append(ys, avgTimeUsage)

		fmt.Println(fmt.Sprintf("x: %d SizeKB", i/types.SizeKB))
	}

	// response
	resp := MultiSimResp{
		Xs: xs,
		Ys: ys,
	}
	data, _ := json.Marshal(resp)
	Response(w, data)
}

func (c *StewpotController) GetMessage(w http.ResponseWriter, r *http.Request) {
	msgID, err := urlParamToInt(r, "msg_id")
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}

	msg := c.stewpot.GetMsg(int64(msgID))
	Response(w, msg)
}

func (c *StewpotController) GetNodeStatus(w http.ResponseWriter, r *http.Request) {
	nodeID, err := urlParamToString(r, "node_id")
	if err != nil {
		ErrorResponse(w, err.Error())
		return
	}

	node := c.stewpot.GetNode(nodeID)
	if node == nil {
		ErrorResponse(w, fmt.Sprintf("can't find node: %s", nodeID))
		return
	}
	Response(w, node)
}

func urlParamToInt(r *http.Request, key string) (int, error) {
	paramStr := r.URL.Query()[key]
	if len(paramStr) == 0 {
		return 0, fmt.Errorf("cannot find parameter: %s", key)
	}
	param, err := strconv.Atoi(paramStr[0])
	if err != nil {
		return param, fmt.Errorf("cannot format key %s from string to int: %v", key, err)
	}
	return param, nil
}

func urlParamToString(r *http.Request, key string) (string, error) {
	paramStr := r.URL.Query()[key]
	if len(paramStr) == 0 {
		return "", fmt.Errorf("cannot find parameter: %s", key)
	}
	return paramStr[0], nil
}
