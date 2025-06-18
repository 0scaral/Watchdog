const apiURL = 'http://localhost:8080/logs/stored';

async function fetchStoredLogs() {
    try {
        const response = await fetch(apiURL);
        if (!response.ok) throw new Error('Error fetching stored logs');
        return await response.json();
    } catch (e) {
        return [];
    }
}

async function deleteStoredLog(id) {
    try {
        const res = await fetch(`http://localhost:8080/logs/stored/id/${id}`, {
            method: 'DELETE'
        });
        if (res.ok) {
            updateStoredLogs();
        } else {
            alert('Error deleting log');
        }
    } catch (e) {
        alert('Error deleting log');
    }
}

function renderStoredLogs(logs) {
    const logsList = document.getElementById('stored-logs-list');
    logsList.innerHTML = '';
    if (!Array.isArray(logs) || logs.length === 0) {
        logsList.innerHTML = '<li>No stored logs found.</li>';
        return;
    }
    logs.forEach(log => {
        const li = document.createElement('li');
        li.innerHTML = `<strong>ID:</strong> ${log.id} - <strong>TIMESTAMP:</strong> [${log.timeCreated}] - <strong>TYPE:</strong> (${log.levelDisplayName}) - <strong>MESSAGE:</strong> ${log.message} `;
        const delBtn = document.createElement('button');
        delBtn.textContent = 'Delete';
        delBtn.style.marginLeft = '10px';
        delBtn.onclick = () => deleteStoredLog(log.id);
        li.appendChild(delBtn);
        logsList.appendChild(li);
    });
}

async function updateStoredLogs() {
    const logs = await fetchStoredLogs();
    renderStoredLogs(logs);
}

setInterval(updateStoredLogs, 5000);
updateStoredLogs();

// Buscar log almacenado por ID
document.getElementById('search-stored-log-btn').addEventListener('click', async () => {
    const logId = document.getElementById('stored-log-search').value.trim();
    if (!logId) return;
    try {
        const res = await fetch(`http://localhost:8080/logs/stored/id/${logId}`);
        if (!res.ok) throw new Error('Not found');
        const log = await res.json();
        renderStoredLogs([log]);
    } catch (e) {
        renderStoredLogs([]);
    }
});

// Buscar logs almacenados por tipo
document.getElementById('search-stored-log-btn-type').addEventListener('click', async () => {
    const logType = document.getElementById('stored-log-search-type').value.trim();
    if (!logType) return;
    try {
        const res = await fetch(`http://localhost:8080/logs/stored/type/${logType}`);
        if (!res.ok) throw new Error('Not found');
        const logs = await res.json();
        renderStoredLogs(Array.isArray(logs) ? logs : []);
    } catch (e) {
        renderStoredLogs([]);
    }
});

