let metadata = [];
let clientLogs = [];

export function writeLogs(port, log) {
	clientLogs.push(`${new Date().toISOString().replace('T', ' ').replace('Z', '')} client at port=${port} ${log}`);
}

document.getElementById('exportLogs').addEventListener('click', function () {
	let tempLink = document.createElement('a');
	let uri = `data:application/octet-stream,${encodeURIComponent(clientLogs)}`;
	tempLink.setAttribute('download', `client-${location.port}.log`);
	tempLink.setAttribute('href', uri);
	console.log(clientLogs);
	tempLink.click();
});

function formatMetaData(tempMetaData) {
	tempMetaData = tempMetaData.map((server) => {
		let [min, max] = [Infinity, 0];
		server.Tablets.forEach((tablet) => {
			min = Math.min(min, +tablet.From);
			max = Math.max(max, +tablet.to);
		});
		return { min, max };
	});
}

function getCrrentMetaData() {
	console.log('fettt', clientLogs);
	fetch('http://localhost:3030/')
		.then((res) => {
			metadata = formatMetaData(res.json);
			writeLogs(location.port, 'fetched METADATA');
		})
		.catch((err) => console.log(err));
}

document.documentElement.addEventListener('load', getCrrentMetaData());

export function getRowsCorrectedURLs(paramValue) {
	let urls = [''];
	paramValue = paramValue.sort();
	let breakPoint = Infinity;
	paramValue.forEach((key, index) => {
		if (key > metadata[0].max) {
			breakPoint = Math.min(breakPoint, index);
		}
	});
	if (breakPoint === 0) {
		urls[0] = `:3032/rows?list=${JSON.stringify(paramValue).replace(/[\[\]']+/g, '')}`;
	} else if (breakPoint === Infinity) {
		urls[0] = `:3031/rows?list=${JSON.stringify(paramValue).replace(/[\[\]']+/g, '')}`;
	} else {
		urls[0] = `:3031/rows?list=${JSON.stringify(paramValue.slice(0, breakPoint)).replace(/[\[\]']+/g, '')}`;
		urls.push(`:3032/rows?list=${JSON.stringify(paramValue.slice(breakPoint, Infinity)).replace(/[\[\]']+/g, '')}`);
	}
	return urls;
}

export function getServerPort(rowKey) {
	if (rowKey >= metadata[0].min && rowKey <= metadata[0].max) {
		return '3031';
	} else if (rowKey >= metadata[1].min && rowKey <= metadata[1].max) {
		return '3032';
	} else {
		return '';
	}
}

// schedule metadata fetching every 2 seconds till infinity
setInterval(getCrrentMetaData(), 2000);