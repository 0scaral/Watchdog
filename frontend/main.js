const apiURL = 'http://localhost:8080/logs/';

async function getFetchLogs() {
    try {
        const response = await fetch(apiURL);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const logs = await response.json();
        console.log('Logs fetched successfully:', logs);
    } catch (error) {
        console.error('Error fetching logs:', error);
    }
}

async function fetchLogs() {
    try {
        const response = await fetch(apiURL);
        if (!response.ok) throw new Error('Error fetching logs');
        return await response.json();
    } catch (e) {
        return [];
    }
}

// Función para renderizar los logs en la lista
function renderLogs(logs) {
    const logsList = document.getElementById('logs-list');
    logsList.innerHTML = '';
    if (logs.length === 0) {
        logsList.innerHTML = '<li>No logs found.</li>';
        return;
    }
    logs.forEach(log => {
        const li = document.createElement('li');
        li.textContent = `[${log.timeCreated}] (${log.levelDisplayName}) ${log.message}`;
        logsList.appendChild(li);
    });
}

async function updateLogs() {
    const logs = await fetchLogs();
    renderLogs(logs);
}
setInterval(updateLogs, 5000);

updateLogs();
getFetchLogs();