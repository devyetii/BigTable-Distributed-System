let metadata = [];
let clientLogs = '';
const MASTER = 'http://localhost:3030',
	SERVER1 = 'http://localhost:3031',
	SERVER2 = 'http://localhost:3032',
	infinity = 10000000;

function writeLogs(port, log) {
	clientLogs +=
		`${new Date().toISOString().replace('T', ' ').replace('Z', '')} client at port=${port} ${log}` + '\r\n';
}

document.getElementById('exportLogs').addEventListener('click', function () {
	let tempLink = document.createElement('a');
	let uri = `data:application/octet-stream,${encodeURIComponent(clientLogs)}`;
	tempLink.setAttribute('download', `client-${location.port}.log`);
	tempLink.setAttribute('href', uri);
	tempLink.click();
});

function formatMetaData(tempMetaData) {
	tempMetaData = tempMetaData.map((server) => {
		let [min, max] = [infinity, 0];
		server.Tablets.forEach((tablet) => {
			if (tablet.From !== tablet.To) {
				min = Math.min(min, tablet.From);
				max = Math.max(max, tablet.To);
			}
		});
		return { min, max };
	});
	return tempMetaData;
}

function getCrrentMetaData() {
	fetch(`${MASTER}/metadata`)
		.then(async (res) => {
			metadata = formatMetaData(JSON.parse(await res.text()));
			console.log(metadata);
			writeLogs(location.port, 'fetched METADATA');
		})
		.catch((err) => console.log(err));
}

getCrrentMetaData();

function getRowsCorrectedURLs(paramValue) {
	let urls = [''];
	paramValue = paramValue.sort();
	let breakPoint = infinity;
	paramValue.forEach((key, index) => {
		if (key > metadata[0].max) {
			breakPoint = Math.min(breakPoint, index);
		}
	});
	if (breakPoint === 0) {
		urls[0] = `${SERVER2}/row?list=${JSON.stringify(paramValue).replace(/[\[\]']+/g, '')}`;
	} else if (breakPoint === infinity) {
		urls[0] = `${SERVER1}/row?list=${JSON.stringify(paramValue).replace(/[\[\]']+/g, '')}`;
	} else {
		urls[0] = `${SERVER1}/row?list=${JSON.stringify(paramValue.slice(0, breakPoint)).replace(/[\[\]']+/g, '')}`;
		urls.push(
			`${SERVER2}/row?list=${JSON.stringify(paramValue.slice(breakPoint, infinity)).replace(/[\[\]']+/g, '')}`
		);
	}
	return urls;
}

function getServer(rowKey) {
	if (rowKey >= metadata[0].min && rowKey <= metadata[0].max) {
		return SERVER1;
	} else if (rowKey >= metadata[1].min && rowKey <= metadata[1].max) {
		return SERVER2;
	} else {
		return '';
	}
}

// schedule metadata fetching every 2 seconds till infinity
setInterval(getCrrentMetaData, 2000);
