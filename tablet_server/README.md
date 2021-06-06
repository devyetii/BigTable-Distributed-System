## Client API Documentation

### Get Rows (`GET /rows`)
- Query Parameters
    - `list` : 
        - syntax `<row_key_type>,<row_key_type>,<row_key_type>,...` i.e. comma separated list of `row_key`'s
        - Return entries with the given `row_key` values in the list. If a given `row_key` is not found; it's ignored.

- Success Response (`200 OK`)
```json
{
    row_key : {
        col_key: col_val,
        col_key: col_val,
        .
        .
        .
    },
    .
    .
    .
}
```
- Error Responses
    - `400 Invalid Keys` in case one of the given list values is invalid i.e. not not castable to `row_key_type`


### Add a row (`POST /row/:key`)
- Route Parameters
    - `key row_key_type` : The key of the row to add
- Body: the cells to add to the row
```json
{
    col_key: col_val,
    col_key: col_val,
    .
    .
    .
}
```
- Success Response (`200 OK`) : Should be similar to the given Body
- Error Responses
    - `400 Invalid row key` : The given `key` route param is not castable to `row_key_type`
    - `400 Row already exists` : The `key` route param given is already an existing `row_key` within the data

### Add/Edit Row Cells (`PUT /row/:key/cells`)
- Route Parameters
    - `key row_key_type` : The key of the row to edit
- Body: the cells to add/edit
```json
{
    col_key: col_val,
    col_key: col_val,
    .
    .
    .
}
```
- Success Response (`200 OK`) : The complete row you wish to add to after adding/editing its cells. Its format is similar to the body
- Error Responses
    - `400 Invalid row key` : The given `key` route param is not castable to `row_key_type`
    - `400 Invalid column key` : one of the given `col_key`'s in the body is not castable to `col_key_type`
    - `404 Row not fouund` : The `key` route param given doesn't map to an existing `row_key` within the data

### Delete cells within a row (`PUT /row/:key/cells/delete`)
- Route Parameters
    - `key row_key_type` : The key of the row to delete cells from
- Body: array `col_key` of the cells to delete
```json
[ col_key, col_key, col_key, ... ]
```
 > Note : If a col_key is not found, it's ignored
- Success Response (`200 OK`) : The complete row you wish to delete from after deleting its required cells.
```json
{
    col_key: col_val,
    col_key: col_val,
    .
    .
    .
}
```
- Error Responses
    - `400 Invalid row key` : The given `key` route param is not castable to `row_key_type`
    - `400 Invalid column key` : one of the given `col_key`'s in the body is not castable to `col_key_type`
    - `404 Row not found` : The `key` route param given doesn't map to an existing `row_key` within the data
    
### Delete rows (DELETE /rows)
- Query Parameters
    - `list` : 
        - syntax `<row_key_type>,<row_key_type>,<row_key_type>,...` i.e. comma separated list of `row_key`'s
        - Deletes entries with the given `row_key` values in the list. If a given `row_key` is not found; it's ignored.
- Success Response (`200 OK`) : `Deleted n rows`
- Error Responses
    - `400 Invalid Keys` in case one of the given list values is invalid i.e. not not castable to `row_key_type`
