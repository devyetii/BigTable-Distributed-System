import { getCorrectedURLs } from './src.js';

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
	let rangeInput = document.getElementById('range');
	let urlParam = '?';
	let urls = [];
	let currentRows = [];

	if (rangeInput.nodeValue !== '') {
		urlParam += 'range=' + rangeInput.nodeValue;
		urls = getCorrectedURLs('range');
	} else if (listInput.nodeValue !== '') {
		urlParam += 'list=' + listInput.nodeValue;
		urls = getCorrectedURLs('list');
	} else {
		alert('no field filled correctly');
		return;
	}

	fetch(`http://localhost${urls[0]}`)
		.then((res) => {
			currentRows = res.json;
			if (urls.length > 1) {
				fetch(`http://localhost${urls[1]}`)
					.then((res2) => {
						currentRows = currentRows.concat(res2.json);
						renderRows(currentRows);
					})
					.catch((err2) => console.log(err2));
			} else {
				renderRows(currentRows);
			}
		})
		.catch((err) => console.log(err));
}
