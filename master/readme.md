# How To Run

> go run . port GFSHostAddress TabletServerAddress

## example

> go run . 6700 http://localhost:3033 http://localhost:3500

### this will run the master server on port 6700, and will set GFS server address to http://localhost:3033, and Tablet server addresss to http://localhost:3500

# API

<div style="font-size:16px">
<b style="font-size:18px">GET /metadata</b>

<p>return metadata object</p>

<b  style="font-size:18px" >GET /load-balance-change</b>

<p>notify master server to recompute metadata</p>

<b  style="font-size:18px">GET /server-id</b>

<p>return server id to be assigned to tablet server, and send a serve request to tablet server with all needed metadata</p>

</div>
