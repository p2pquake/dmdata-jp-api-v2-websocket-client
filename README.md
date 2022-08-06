# dmdata-jp-api-v2-websocket-client

[DMDATA.JP (Project DM-D.S.S)](https://dmdata.jp/) API v2 の非公式クライアント実装です。

現在は、 WebSocket での電文受信に特化しています。

## 対応 API

- Socket v2: パラメタは一部のみ対応
- WebSocket v2: 非圧縮か gzip の XML 電文のみ対応

## ダウンロード

実行可能なバイナリは [Releases](https://github.com/p2pquake/dmdata-jp-api-v2-websocket-client/releases) にあります。

## 使用方法

### コマンド

`--help` オプションで説明が確認できます。

```sh
$ ./dmdata-jp-api-v2-websocket-client --help
DMDATA.JP (Project DM-D.S.S) API v2 の非公式クライアント実装

Usage:
  dmdata-jp-api-v2-websocket-client [flags]

Examples:
DMDATA_JP_API_KEY=<API_KEY> ./dmdata-jp-api-v2-websocket-client -k -c telegram.earthquake -c eew.warning

Flags:
  -c, --classifications strings     Retreive classifications (default [telegram.earthquake])
  -h, --help                        help for dmdata-jp-api-v2-websocket-client
  -k, --keep-existing-connections   Keep existing connections
  -v, --version                     version for dmdata-jp-api-v2-websocket-client
```

オプション|説明
--|--
-c|[Socket Start v2](https://dmdata.jp/docs/reference/api/v2/socket.start) のリクエストパラメータ `classifications` の値（複数回指定可）。デフォルトは `telegram.earthquake` です。
-k|既存の WebSocket 接続を切断しません。

### 電文受信時の挙動

電文は xml ディレクトリに出力します。

なお、電文出力時は`.xml.tmp` に書き出した上で `.xml` にアトミックにリネームします。そのため、 `.xml` ファイルはいつ読み取っても完全な状態です（書き込み途中の状態を読み取ることはありません）。

### 実行例

```sh
$ mkdir xml
$ DMDATA_JP_API_KEY=<API_KEY> ./dmdata-jp-api-v2-websocket-client
2021/04/05 00:03:46.685373 DMDATA.JP API v2 WebSocket client (unofficial)
2021/04/05 00:03:46.685425 GET https://api.dmdata.jp/v2/socket?status=open
2021/04/05 00:03:46.900731 POST https://api.dmdata.jp/v2/socket with dmdata.StartSocketRequest{Classifications:[]string{"telegram.earthquake"}, Types:[]string(nil), AppName:"P2PQuakeWSV2Client 0.1"}
2021/04/05 00:03:46.924852 connecting to wss://<VARIABLE>.api.dmdata.jp/v2/websocket?ticket=<TICKET>
2021/04/05 00:03:47.003456 connected.
2021/04/05 00:03:47.264121 recv: {"type":"start","socketId":<SOCKET_ID>,"classifications":["telegram.earthquake"],"types":null,"test":"no","formats":["xml","a/n","binary"],"appName":"P2PQuakeWSV2Client 0.1","time":"2021-04-04T15:03:47.255Z"  }
2021/04/05 00:04:07.620397 recv: {"type":"ping","pingId":"zmiL"}
2021/04/05 00:04:07.620897 pong: dmdata.Ping{Type:"pong", PingId:"zmiL"}
2021/04/05 02:21:44.001913 recv: {"type":"data","version":"2.0","id":"<ID>","classification":"telegram.earthquake", *snip*}
2021/04/05 02:21:44.001913 recv: {"type":"data","version":"2.0","id":"<ID>","classification":"telegram.earthquake","head":{"type":"VXSE51", *snip*}, *snip*}
2021/04/05 02:22:26.996647 recv: {"type":"data","version":"2.0","id":"<ID>","classification":"telegram.earthquake","head":{"type":"VXSE52", *snip*}, *snip*}

$ ls -l xml/
-rw-r--r-- 1 user user  1817  4月  5 02:21 20210405022144_02967774_VXSE51.xml
-rw-r--r-- 1 user user  1659  4月  5 02:22 20210405022226_06028089_VXSE52.xml
```

## 参考

- [DMDATA.JP API | DMDATA.JP Document](https://dmdata.jp/doc/reference/)
- [p2pquake/jmaxml-seis-parser-go: 気象庁防災情報 XML の パーサー (Go 実装、地震情報・津波予報のみ対応)](https://github.com/p2pquake/jmaxml-seis-parser-go)
