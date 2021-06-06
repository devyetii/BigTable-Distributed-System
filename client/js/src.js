let metadata = [];

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
	fetch('http://localhost:3030/')
		.then((res) => {
			metadata = formatMetaData(res.json);
		})
		.catch((err) => console.log(err));
}

document.addEventListener('load', getCrrentMetaData);

export function getCorrectedURLs(paramType, paramValue) {
	let urls = [''];

	if (paramType === 'range') {
		if (paramValue[0] >= metadata[0].min && paramValue[1] <= metadata[0].max) {
			urls[0] = `:3031?range=${paramValue[0]}-${paramValue[1]}`;
		} else if (paramValue[0] >= metadata[1].min && paramValue[1] <= metadata[1].max) {
			urls[0] = `:3032?range=${paramValue[0]}-${paramValue[1]}`;
		} else {
			urls[0] = `:3031?range=${paramValue[0]}-${metadata[metadata[0].max < metadata[1].min ? 0 : 1].max}`;
			urls.push(`:3032?range=${metadata[!(metadata[0].max < metadata[1].min) ? 0 : 1].min}-${paramValue[1]}`);
		}
	} else {
		paramValue = paramValue.sort();
		let breakPoint = Infinity;
		paramValue.forEach((key, index) => {
			if (key > metadata[0].max) {
				breakPoint = Math.min(breakPoint, index);
			}
		});
		if (breakPoint === 0) {
			urls[0] = `:3032?list=${JSON.stringify(paramValue).replace(/[\[\]']+/g, '')}`;
		} else if (breakPoint === Infinity) {
			urls[0] = `:3031?list=${JSON.stringify(paramValue).replace(/[\[\]']+/g, '')}`;
		} else {
			urls[0] = `:3031?list=${JSON.stringify(paramValue.slice(0, breakPoint)).replace(/[\[\]']+/g, '')}`;
			urls.push(`:3032?list=${JSON.stringify(paramValue.slice(breakPoint, Infinity)).replace(/[\[\]']+/g, '')}`);
		}
	}

	return urls;
}

// schedule metadata fetching every 2 seconds till infinity
// setInterval(getCrrentMetaData(), 2000);
