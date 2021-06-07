function singleRowRequest(method, url, bodyUnParsed, successMessage, logData) {
	let parsedValues = bodyUnParsed.split(',');
	let body = {};

	parsedValues.forEach((element) => {
		let [key, value] = element.split(':');
		body[key] = value;
	});

	fetch(url, {
		method,
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(body),
	})
		.then((res) => {
			writeLogs(logData.port, logData.message);
			if (res.ok) {
				alert(successMessage);
			} else {
				alert(`${method} Request Failed`);
			}
			console.log(res);
		})
		.catch((err) => {
			alert(err);
		});
}

document.getElementById('addRowButton').addEventListener('click', function () {
	let addRowKeyInput = document.getElementById('addRowKeyInput');
	let addRowColKeysInput = document.getElementById('addRowColKeysInput');
	if (addRowColKeysInput.value === '' || addRowKeyInput.value === '') {
		return;
	}
	let rowKey = +addRowKeyInput.value;
	let port = getServerPort(rowKey);
	if (port === '') {
		alert('invalid row key');
		return;
	}
	singleRowRequest(
		'POST',
		`http://localhost:${port}/row/${rowKey}`,
		addRowColKeysInput.value,
		'row added successfully',
		{ port, message: `Add Row to server:${port}` }
	);
});

document.getElementById('editCellsButton').addEventListener('click', function () {
	let editCellsRowKeyInput = document.getElementById('editCellsRowKeyInput');
	let editCellsColKeysInput = document.getElementById('editCellsColKeysInput');
	if (editCellsColKeysInput.value === '' || editCellsRowKeyInput.value === '') {
		return;
	}
	let rowKey = +editCellsRowKeyInput.value;
	let port = getServerPort(rowKey);
	if (port === '') {
		alert('invalid row key');
		return;
	}
	singleRowRequest(
		'PUT',
		`http://localhost:${port}/row/${rowKey}/cells`,
		editCellsColKeysInput.value,
		'row cells modified successfully',
		{ port, message: `Edit Cell of Row:${rowKey} on server:${port}` }
	);
});

document.getElementById('deleteCellsButton').addEventListener('click', function () {
	let deleteCellsRowKeyInput = document.getElementById('deleteCellsRowKeyInput');
	let deleteCellsColKeysInput = document.getElementById('deleteCellsColKeysInput');
	if (deleteCellsColKeysInput.value === '' || deleteCellsRowKeyInput.value === '') {
		return;
	}
	let rowKey = +deleteCellsRowKeyInput.value;
	let port = getServerPort(rowKey);
	if (port === '') {
		alert('invalid row key');
		return;
	}
	singleRowRequest(
		'DELETE',
		`http://localhost:${port}/row/${rowKey}/cells/delete`,
		deleteCellsColKeysInput.value,
		'row cells deleted successfully',
		{ port, message: `Delete Cell of Row:${rowKey} on server:${port}` }
	);
});

document.getElementById('deleteRowsButton').addEventListener('click', function () {
	let deleteRowsInput = document.getElementById('deleteRowsInput');
	if (deleteRowsInput.value === '') {
		return;
	}
	let rowsList = deleteRowsInput.value.split(',').map((elm) => +elm);
	let urls = getRowsCorrectedURLs(rowsList);

	fetch(`http://localhost${urls[0]}`, { method: 'DELETE' })
		.then(async (res) => {
			alert(res);
			writeLogs(location.port, `DELETE rows from server${urls[0].split('?')[0]}`);
			if (urls.length > 1) {
				fetch(`http://localhost${urls[1]}`, { method: 'DELETE' })
					.then(async (res2) => {
						alert(res2);
						writeLogs(location.port, `DELETE rows from server${urls[1].split('?')[0]}`);
					})
					.catch((err2) => {
						alert(err2);
						console.log(err2);
					});
			}
		})
		.catch((err) => {
			alert(err);
			console.log(err);
		});
});
