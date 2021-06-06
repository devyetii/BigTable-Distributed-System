import { getRowsCorrectedURLs, writeLogs } from './src.js';

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
	rows.forEach((row) => {
		parentArea.appendChild(renderRow(row));
		parentArea.appendChild(document.createElement('div').setAttribute('class', 'divider'));
	});
}

export function getRows() {
	let listInput = document.getElementById('list');

	if (listInput.value === '') {
		alert('invalid data');
		return;
	}

	let rowsList = listInput.value.split(',').map((elm) => +elm);
	let urls = getRowsCorrectedURLs(rowsList);
	let currentRows = [];

	fetch(`http://localhost${urls[0]}`)
		.then((res) => {
			writeLogs(location.port, `GET from server${urls[0].split('?')[0]}`);
			currentRows = res.json;

			if (urls.length > 1) {
				fetch(`http://localhost${urls[1]}`)
					.then((res2) => {
						writeLogs(location.port, `GET from server${urls[1].split('?')[0]}`);
						currentRows = currentRows.concat(res2.json);
						renderRows(currentRows);
					})
					.catch((err2) => alert(err2));
			} else {
				renderRows(currentRows);
			}
		})
		.catch((err) => alert(err));
}
