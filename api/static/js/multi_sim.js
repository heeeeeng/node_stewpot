window.onload = main;

dom = document.getElementById("container");
dom_iternum = document.getElementById("iter_num");
dom_nodenum = document.getElementById("node_num");
dom_maxin = document.getElementById("max_in");
dom_maxout = document.getElementById("max_out");
dom_bandwidth = document.getElementById("bandwidth");
dom_max_msg_size = document.getElementById("msg_size");

myChart = echarts.init(dom);

function main() {

}

function simulate() {
    let iter_num = dom_iternum.value;
    let node_num = dom_nodenum.value;
    let max_in = dom_maxin.value;
    let max_out = dom_maxout.value;
    let bandwidth = dom_bandwidth.value * SIZE_MB;
    let max_msg_size = dom_max_msg_size.value * SIZE_BYTE;

    let url = "/multi_sim/simulate?" +
        "iter_num=" + iter_num +
        "&node_num=" + node_num +
        "&max_in=" + max_in +
        "&max_out=" + max_out +
        "&bandwidth=" + bandwidth +
        "&max_msg_size=" + max_msg_size;

    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open( "GET", url, false ); // false for synchronous request
    xmlHttp.send( null );

    data = JSON.parse(xmlHttp.responseText);
    initChart(data);
}

function initChart(data) {
    option = {
        tooltip: {
            trigger: 'axis'
        },
        toolbox: {
            show: true,
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                dataView: {readOnly: false},
                magicType: {type: ['line', 'bar']},
                restore: {},
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            data: data.xs
        },
        yAxis: {
            type: 'value'
        },
        series: [{
            type: 'line',
            data: data.ys
        }]
    };

    myChart.setOption(option);
}

const SIZE_BIT = 1;
const SIZE_BYTE = 8*SIZE_BIT;
const SIZE_KB = 1024*SIZE_BYTE;
const SIZE_MB = 1024*SIZE_KB;
const SIZE_GB = 1024*SIZE_MB;
const SIZE_TB = 1024*SIZE_GB;