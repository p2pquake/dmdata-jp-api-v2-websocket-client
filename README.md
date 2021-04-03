# dmdata-jp-api-v2-websocket-client

[DMDATA.JP (Project DM-D.S.S)](https://dmdata.jp/) API v2 の非公式クライアント実装です。

現在は、地震情報・津波予報を WebSocket で受信することに特化しています。

## 対応 API

- Socket v2: パラメタは一部のみ対応
- WebSocket v2: 非圧縮か gzip の XML 電文のみ対応

## 実行例

```sh
$ mkdir xml
$ DMDATA_JP_API_KEY=<API_KEY> go run main.go
DMDATA.JP API v2 WebSocket client (unofficial)
2021/04/03 21:40:15 GET https://api.dmdata.jp/v2/socket?status=open
2021/04/03 21:40:15 POST https://api.dmdata.jp/v2/socket with dmdata.StartSocketRequest{Classifications:[]string{"telegram.earthquake"}, Types:[]string(nil), AppName:"P2PQuakeWSV2Client 0.1"}
2021/04/03 21:40:15 connecting to wss://<VARIABLE>.api.dmdata.jp/v2/websocket?ticket=<TICKET>
2021/04/03 21:40:15 connected.
2021/04/03 21:40:15 recv: {"type":"start","socketId":<SOCKET_ID>,"classifications":["telegram.earthquake"],"types":null,"test":"no","formats":["xml","a/n","binary"],"appName":"P2PQuakeWSV2Client 0.1","time":"2021-04-03T12:40:15.471Z"}
2021/04/03 21:40:35 recv: {"type":"ping","pingId":"FpKr"}
2021/04/03 21:40:35 pong: dmdata.Ping{Type:"pong", PingId:"FpKr"}

$ ls xml/
20210403220016_02906503_VXSE51.xml
20210403220110_01161598_VXSE52.xml
20210403220213_09189757_VXSE53.xml
```

## 参考

- [DMDATA.JP API | DMDATA.JP Document](https://dmdata.jp/doc/reference/)
- [p2pquake/jmaxml-seis-parser-go: 気象庁防災情報 XML の パーサー (Go 実装、地震情報・津波予報のみ対応)](https://github.com/p2pquake/jmaxml-seis-parser-go)
