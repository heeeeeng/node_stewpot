package api

import "net/http"

func (c *StewpotController) RegisterRouters() {
	http.HandleFunc("/", c.MainPage)
	http.HandleFunc("/static/", c.Static)
	http.HandleFunc("/restart", c.Restart)
	http.HandleFunc("/graph", c.GetNetworkGraph)
	http.HandleFunc("/send_msg", c.SendMsg)
	http.HandleFunc("/all_node_send_msg", c.AllNodeSendMsg)
	http.HandleFunc("/time_unit", c.GetTimeUnit)

	http.HandleFunc("/multi_sim", c.SimPage)
	http.HandleFunc("/multi_sim/simulate", c.MultiSimulate)
}
