
window.onload = main;

dom = document.getElementById("container");
myChart = echarts.init(dom);

var graph;
var intervalID;

function main() {
    initChart();

    intervalID = window.setInterval(internalFunc, 500);


}

function initChart() {
    // get graph data
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", "/graph", false ); // false for synchronous request
    xmlHttp.send( null );

    graph = JSON.parse(xmlHttp.responseText);

    option = null;
    myChart.showLoading();
    myChart.hideLoading();

    graph.nodes.forEach(function (node) {
        node.itemStyle = null;
        node.symbolSize = 10;
        // node.value = node.symbolSize;
        // node.category = node.attributes.modularity_class;
        // Use random x, y
        node.x = node.y = null;
        node.draggable = true;
    });

    option = {
        title: {
            text: 'Les Miserables',
            subtext: 'Default layout',
            top: 'bottom',
            left: 'right'
        },
        tooltip: {},
        animation: false,
        series : [
            {
                name: 'Les Miserables',
                type: 'graph',
                layout: 'force',
                data: graph.nodes,
                links: graph.links,
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

    myChart.setOption(option);

    myChart.setOption(option);
    if (option && typeof option === "object") {
        myChart.setOption(option, true);
    }

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
    option = myChart.getOption();

    nodes = myChart.getModel().getSeriesByIndex(0).preservedPoints;
    for (var i=0; i<graph.nodes.length; i++) {
        var nodeName = graph.nodes[i].name;
        graph.nodes[i].x = nodes[nodeName][0];
        graph.nodes[i].y = nodes[nodeName][1];
    }

    option.series[0].layout = "none";
    option.series[0].data = graph.nodes;
    option.series[0].force = null;

    myChart.setOption(option);

    node_0 = myChart.getModel().getSeriesByIndex(0).preservedPoints[graph.nodes[0].name];
    console.log(node_0);
}

