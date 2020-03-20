window.onload = main;

dom = document.getElementById("container");
dom_nodenum = document.getElementById("node_num");
dom_maxin = document.getElementById("max_in");
dom_maxout = document.getElementById("max_out");
dom_bandwidth = document.getElementById("bandwidth");

myChart = echarts.init(dom);
links = {};
nodesSet = {};

let base_option;
let graph;
let intervalID;
// let draw_switch = 1;

function main() {
    initChart();

    intervalID = window.setInterval(internalFunc, 500);

}

function restart() {
    let node_num = dom_nodenum.value;
    let max_in = dom_maxin.value;
    let max_out = dom_maxout.value;
    let bandwidth = dom_bandwidth.value * SIZE_MB;

    let url = "/restart?" + "node_num=" + node_num + "&max_in=" + max_in + "&max_out=" + max_out + "&bandwidth=" + bandwidth;

    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", url, false ); // false for synchronous request
    xmlHttp.send( null );

    window.location.reload();
}

function initChart() {
    // get graph data
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", "/graph", false ); // false for synchronous request
    xmlHttp.send( null );

    graph = JSON.parse(xmlHttp.responseText);

    myChart.showLoading();
    myChart.hideLoading();

    for (let i in graph.nodes) {
        let node = graph.nodes[i];

        node.itemStyle = null;
        node.symbolSize = 10;
        node.x = node.y = null;
        node.draggable = true;

        graph.nodes[i] = node;
        nodesSet[node.name] = i;
    }

    for (let i in graph.links) {
        let link = graph.links[i];
        let source = link.source;
        let target = link.target;

        if (links[source] == null) {
            links[source] = {};
        }
        links[source][target] = i;

        if (links[target] == null) {
            links[target] = {};
        }
        links[target][source] = i;
    }

    base_option = {
        title: {
            text: 'OG Node Stewpot',
            subtext: 'Default layout',
            top: 'bottom',
            left: 'right'
        },
        tooltip: {},
        animation: false,
        series : [
            {
                name: 'OG Node Stewpot',
                type: 'graph',
                layout: 'force',
                data: graph.nodes,
                links: graph.links,
                // draggable: true,
                focusNodeAdjacency: true,
                // categories: categories,
                roam: true,
                label: {
                    normal: {
                        position: 'right'
                    }
                },
                force: {
                    repulsion: 200
                },
                emphasis: {
                    lineStyle: {
                        width: 10
                    }
                }
            }
        ]
    };

    myChart.setOption(base_option);
    if (base_option && typeof base_option === "object") {
        myChart.setOption(base_option, true);
    }

    console.log("init: ", base_option.series[0]);

}

var cor0_0;
var cor0_1;

function internalFunc() {
    node_0 = myChart.getModel().getSeriesByIndex(0).preservedPoints[graph.nodes[0].name];
    if (cor0_0 != node_0[0] || cor0_1 != node_0[1]) {
        cor0_0 = node_0[0];
        cor0_1 = node_0[1];
        return;
    }
    console.log(node_0);

    window.clearInterval(intervalID);
    switchToNoneLayout();
}

function switchToNoneLayout() {
    nodes = myChart.getModel().getSeriesByIndex(0).preservedPoints;
    for (var i=0; i<graph.nodes.length; i++) {
        var nodeName = graph.nodes[i].name;
        graph.nodes[i].x = nodes[nodeName][0];
        graph.nodes[i].y = nodes[nodeName][1];
    }

    base_option = myChart.getOption();
    base_option.series[0].layout = "none";
    base_option.series[0].data = graph.nodes;
    base_option.series[0].force = null;

    myChart.setOption(base_option);

    // node_0 = myChart.getModel().getSeriesByIndex(0).preservedPoints[graph.nodes[0].name];
    console.log("switched to none");

    // sendMsg();
}

let node_color_old = "red";
let node_color_new = "yellow";
let node_color_set = [ "red", "yellow", "blue", "green", "black" ];

async function sendMsg() {
    // generate new node color
    while (true) {
        if (node_color_new !== node_color_old) {
            break;
        }
        node_color_new = node_color_set[Math.floor(Math.random()*node_color_set.length)];
    }
    node_color_old = node_color_new;

    let xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", "/send_msg", false ); // false for synchronous request
    xmlHttp.send( null );

    console.log(xmlHttp.responseText);
    let t = parseInt(xmlHttp.responseText);

    while (getTimeUnit(t)) {
        await sleep(1);
        t = t+1;
    }
    // console.log(myChart.getOption());
}

function getTimeUnit(t) {

    let xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", "/time_unit?time=" + t, false ); // false for synchronous request
    xmlHttp.send( null );

    let time_unit = JSON.parse(xmlHttp.responseText);
    if (time_unit.length === 0) {
        return false;
    }

    let option = JSON.parse(JSON.stringify(base_option));
    let graph_nodes = base_option.series[0].data;
    let graph_links = option.series[0].links;
    for (let i in time_unit) {
        let task = time_unit[i];
        if (task.type === TASK_TYPE_CONN_RECV) {
            let recver = task.recver;
            graph_nodes[nodesSet[recver]].itemStyle = {
                color: node_color_new
            };
            continue
        }
        if (task.type === TASK_TYPE_MSG_TRANSMIT_DELAY) {
            let linkIndex = links[task.source][task.target];
            graph_links[linkIndex].lineStyle = {
                color: node_color_new,
                width: 2,
                opacity: 0.3
            };
            continue;
        }
    }

    option.series[0].data = graph_nodes;
    option.series[0].links = graph_links;
    myChart.setOption(option);

    return true;
}

function highLightLink(source, target, base_links) {
    let linkIndex = links[source][target];
    let link = base_links[linkIndex];

    link.lineStyle = {
        color: node_color_new
    };
    base_links[linkIndex] = link;
    return base_links;
}

const TASK_TYPE_CONN_RECV = 1;
const TASK_TYPE_MSG_TRANSMIT_DELAY = 6;

function bin2String(array) {
    var result = "";
    for (var i = 0; i < array.length; i++) {
        result += String.fromCharCode(parseInt(array[i], 2));
    }
    return result;
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

const SIZE_BIT = 1;
const SIZE_BYTE = 8*SIZE_BIT;
const SIZE_KB = 1024*SIZE_BYTE;
const SIZE_MB = 1024*SIZE_KB;
const SIZE_GB = 1024*SIZE_MB;
const SIZE_TB = 1024*SIZE_GB;