function renderRow(row) {
	let rowElement = document.createElement('span');
	let rowContent = '';
	rowElement.setAttribute('class', 'p-10 text-white text-1xl');
	Object.keys(row).forEach((key) => {
		rowContent += `${key}: ${row[key]} <br />`;
	});
	rowElement.innerHTML = rowContent;
	return rowElement;
}

function renderRows(rows) {
	let parentArea = document.getElementById('rows');
	parentArea.innerHTML = '';
	Object.keys(rows).forEach((row) => {
		parentArea.appendChild(renderRow(rows[row]));
	});
}

document.getElementById('getRowsButton').addEventListener('click', function () {
	let listInput = document.getElementById('list');

	if (listInput.value === '') {
		alert('invalid data');
		return;
	}

	let rowsList = listInput.value.split(',').map((elm) => +elm);
	let urls = getRowsCorrectedURLs(rowsList);
	let currentRows = [];

	fetch(urls[0])
		.then(async (res) => {
			writeLogs(location.port, `GET from server${urls[0].includes(SERVER1) ? 1 : 2}`);
			currentRows = JSON.parse(await res.text());

			if (urls.length > 1) {
				fetch(urls[1])
					.then(async (res2) => {
						writeLogs(location.port, `GET from server 2`);
						currentRows = currentRows.concat(JSON.parse(await res2.text()));
						renderRows(currentRows);
					})
					.catch((err2) => alert(err2));
			} else {
				renderRows(currentRows);
			}
		})
		.catch((err) => alert(err));
});
