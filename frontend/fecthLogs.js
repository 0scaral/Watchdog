const apiURL = 'http://localhost:8080/logs/';

async function fetchLogs() {
    try {
        const response = await fetch(apiURL);
        if (!response.ok) throw new Error('Error fetching logs');
        return await response.json();
    } catch (e) {
        return [];
    }
}

function renderLogs(logs) {
    const logsList = document.getElementById('logs-list');
    logsList.innerHTML = '';
    if (logs.length === 0) {
        logsList.innerHTML = '<li>No logs found.</li>';
        return;
    }
    logs.forEach(log => {
        const li = document.createElement('li');
        li.innerHTML = `<strong>ID:</strong>${log.id} - <strong>TIMESTAMP: </strong>[${log.timeCreated}] - <strong>TYPE: </strong>(${log.levelDisplayName}) - <strong>MESSAGE: </strong> ${log.message}`;
        logsList.appendChild(li);
    });
}

async function updateLogs() {
    const logs = await fetchLogs();
    renderLogs(logs);
}
setInterval(updateLogs, 5000);

updateLogs();
fetchLogs();

document.getElementById('save-log-btn').addEventListener('click', async () => {
    const logId = document.getElementById('log-input').value.trim();
    if (!logId) return;
    try {
        const res = await fetch(`http://localhost:8080/logs/stored/id/${logId}`, {
            method: 'POST'
        });
        if (res.ok) {
            alert('Log saved!');
        } else {
            alert('Error saving log');
        }
    } catch (e) {
        alert('Error saving log');
    }
});

// Guardar logs por tipo usando POST (para coincidir con el backend actual)
document.getElementById('save-log-btn-type').addEventListener('click', async () => {
    const logType = document.getElementById('log-input-type').value.trim();
    if (!logType) return;
    try {
        const res = await fetch(`http://localhost:8080/logs/stored/type/${logType}`, {
            method: 'POST'
        });
        if (res.ok) {
            alert('Logs by type saved!');
        } else {
            alert('Error saving logs by type');
        }
    } catch (e) {
        alert('Error saving logs by type');
    }
});

// Buscar log por ID
document.getElementById('search-log-btn').addEventListener('click', async () => {
    const logId = document.getElementById('log-search').value.trim();
    if (!logId) return;
    try {
        const res = await fetch(`http://localhost:8080/logs/id/${logId}`);
        if (!res.ok) throw new Error('Not found');
        const log = await res.json();
        renderLogs([log]);
    } catch (e) {
        renderLogs([]);
    }
});

// Buscar logs por tipo
document.getElementById('search-log-btn-type').addEventListener('click', async () => {
    const logType = document.getElementById('log-search-type').value.trim();
    if (!logType) return;
    try {
        const res = await fetch(`http://localhost:8080/logs/type/${logType}`);
        if (!res.ok) throw new Error('Not found');
        const logs = await res.json();
        renderLogs(Array.isArray(logs) ? logs : []);
    } catch (e) {
        renderLogs([]);
    }
});