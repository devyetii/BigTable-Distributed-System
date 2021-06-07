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
	let serverURL = getServer(rowKey);
	if (serverURL === '') {
		alert('invalid row key');
		return;
	}
	singleRowRequest('POST', `${serverURL}/row/${rowKey}`, addRowColKeysInput.value, 'row added successfully', {
		port: location.port,
		message: `Add Row to server${serverURL === SERVER1 ? 1 : 2}`,
	});
});

document.getElementById('editCellsButton').addEventListener('click', function () {
	let editCellsRowKeyInput = document.getElementById('editCellsRowKeyInput');
	let editCellsColKeysInput = document.getElementById('editCellsColKeysInput');
	if (editCellsColKeysInput.value === '' || editCellsRowKeyInput.value === '') {
		return;
	}
	let rowKey = +editCellsRowKeyInput.value;
	let serverURL = getServer(rowKey);
	if (serverURL === '') {
		alert('invalid row key');
		return;
	}
	singleRowRequest(
		'PUT',
		`${serverURL}/row/${rowKey}/cells`,
		editCellsColKeysInput.value,
		'row cells modified successfully',
		{
			port: location.port,
			message: `Edit Cells from Row:${rowKey} on server${serverURL === SERVER1 ? 1 : 2}`,
		}
	);
});

document.getElementById('deleteCellsButton').addEventListener('click', function () {
	let deleteCellsRowKeyInput = document.getElementById('deleteCellsRowKeyInput');
	let deleteCellsColKeysInput = document.getElementById('deleteCellsColKeysInput');
	if (deleteCellsColKeysInput.value === '' || deleteCellsRowKeyInput.value === '') {
		return;
	}
	let rowKey = +deleteCellsRowKeyInput.value;
	let serverURL = getServer(rowKey);
	if (serverURL === '') {
		alert('invalid row key');
		return;
	}
	let body = deleteCellsColKeysInput.value.split(',');
	fetch(`${serverURL}/row/${rowKey}/cells/delete`, {
		method: 'PUT',
		body: JSON.stringify(body),
	})
		.then((res) => {
			if (res.ok) {
				writeLogs(location.port, `Delete Cell from Row:${rowKey} on server:${serverURL === SERVER1 ? 1 : 2}`);
				alert('row cells deleted successfully');
			} else {
				alert(`${method} Request Failed`);
			}
			console.log(res);
		})
		.catch((err) => {
			alert(err);
		});
});

document.getElementById('deleteRowsButton').addEventListener('click', function () {
	let deleteRowsInput = document.getElementById('deleteRowsInput');
	if (deleteRowsInput.value === '') {
		return;
	}
	let rowsList = deleteRowsInput.value.split(',').map((elm) => +elm);
	let urls = getRowsCorrectedURLs(rowsList);

	fetch(urls[0], { method: 'DELETE' })
		.then(async (res) => {
			writeLogs(location.port, `DELETE rows from server${urls[0].includes(SERVER1) ? 1 : 2}`);
			alert('rows deleted from server 1');
			if (urls.length > 1) {
				fetch(urls[1], { method: 'DELETE' })
					.then(async (res2) => {
						alert('rows deleted from server 2');
						writeLogs(location.port, `DELETE rows from server 2`);
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
