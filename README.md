
gRPC 環境構築
・Protocol buffer コンパイラのインストール
    ・以下のリンクのAssetsから環境に適したzipファイルをダウンロード（windows 64bitの場合は`protoc-xx.x-win64.zip`）
        ・https://github.com/protocolbuffers/protobuf/releases
・zipファイルを解凍し、`C:\Program Files\protoc`内に解凍したファイルを格納する
・環境変数を設定し、パスを通す
    ・windowsの場合、Pathに`C:\Program Files\protoc\protoc-xx.x-win64\bin`を追加する
    ※xx.xはダウンロードしたファイルのバージョンの値が入る
・以下のコマンドを実行し`libprotoc xx.x`が表示されることを確認する
    `protoc --version`
・以下のコマンドをプロジェクト直下で実行しGoプラグインを導入する
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

・protoファイルをビルドする
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/{protoファイル名}.proto
```

```
protoc --proto_path=./proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative ./proto/message.proto
```

※`./proto/{protoファイル名}.proto`には実行するディレクトリを起点とする相対パスを記載する
