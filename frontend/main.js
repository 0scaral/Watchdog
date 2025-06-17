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

getFetchLogs();