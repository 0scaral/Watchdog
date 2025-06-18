const apiCurrentURL = 'http://localhost:8080/metrics/current';

const cpuCtx = document.getElementById('cpuChart').getContext('2d');
const ramCtx = document.getElementById('ramChart').getContext('2d');
const diskCtx = document.getElementById('diskChart').getContext('2d');

const chartConfig = (label, borderColor) => ({
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label,
            data: [],
            borderColor,
            fill: false,
            tension: 0.1
        }]
    },
    options: {
        plugins: {
            title: {
                display: true,
                text: label,
                font: {
                    size: 20,
                    weight: 'bold'
                },
                color: '#e3eafc',
                padding: {top: 10, bottom: 10}
            }
        },
        scales: {
            x: { display: false },
            y: { min: 0, max: 100 }
        }
    }
});

const cpuChart = new Chart(cpuCtx, chartConfig('CPU Usage (%)', '#cf008a'));
const ramChart = new Chart(ramCtx, chartConfig('RAM Usage (%)', '#9c00bf'));
const diskChart = new Chart(diskCtx, chartConfig('Disk Usage (%)', '#0072cf'));

function addData(chart, value) {
    const now = new Date().toLocaleTimeString();
    chart.data.labels.push(now);
    chart.data.datasets[0].data.push(value);
    if (chart.data.labels.length > 30) {
        chart.data.labels.shift();
        chart.data.datasets[0].data.shift();
    }
    chart.update();
}

async function fetchAndUpdate() {
    try {
        const res = await fetch(apiCurrentURL);
        const data = await res.json();
        addData(cpuChart, data.cpu_usage);
        addData(ramChart, data.ram_usage);
        addData(diskChart, data.disk_usage);
    } catch (e) {
    }
}

setInterval(fetchAndUpdate, 2000);
fetchAndUpdate();

