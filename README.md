
# 概要
Discord風のチャットアプリでフレンドとなったユーザーと1対1でのチャットができます。また、グループを作成し、フレンドを招待することでグループチャットが行えます。  
アプリはSPAを採用し、メッセージのリアルタイム更新機能にはgRPCのサーバーストリーミングを使用しました。

アプリURL: https://schaty.net

# 動作イメージ
<img src="https://raw.github.com/wiki/yoshinori0811/chat_app_backend/images/demo.gif" style="width: 100%;" >

# システム構成図
<img src="https://raw.github.com/wiki/yoshinori0811/chat_app_backend/images/system_architecture.png" style="width: 75%;">

# 技術スタック
- フロントエンド
    - React
    - grpc-web
    - tailwind
    - react-redux
- バックエンド
    - Go
    - Gorm
    - grpc
- インフラ
    - MySQL
    - Docker
    - enginx
    - envoy
    - aws (EC2, Route53, VPC)

# ディレクトリ構成
- フロントエンド
```
chat_app_frontend
├─public  // 静的な公開リソースを格納するディレクトリ
└─src  // ソースコードを格納するディレクトリ
    ├─app  // react-reduxのストアやグローバルな設定を格納するディレクトリ
    ├─assets  // 静的なリソースを格納するディレクトリ
    │  └─images
    ├─components  // 関数コンポーネントを格納するディレクトリ
    │  ├─AddRoomMember
    │  ├─dropDownMenu
    │  │  └─messageMenu
    │  ├─home
    │  │  ├─addFriend
    │  │  ├─dm
    │  │  ├─friend
    │  │  ├─friendList
    │  │  ├─friendRequestList
    │  │  └─navbar
    │  ├─message
    │  ├─messageContent
    │  ├─room
    │  └─sidebar
    │      ├─dm
    │      └─room
    │          └─settingsIcon
    │              └─modal
    ├─features  // Redux Toolkitのスライスを格納するディレクトリ
    ├─hooks  // バックエンドとAPI通信を行うロジックを格納するディレクトリ
    ├─pb  // gRPCのスキーマから自動生成されたコードを格納するディレクトリ
    │  └─web
    │      └─src
    │          └─proto
    ├─proto  // gRPCのスキーマを格納するディレクトリ
    └─types  // ユーザー定義型を格納するディレクトリ
```
- バックエンド
```
chat_app_backend
├─config  // 設定ファイルを格納するディレクトリ
├─controller  // 各エンドポイントで呼び出される処理を格納するディレクトリ
├─db  // データベースとの接続に関する処理を格納するディレクトリ
├─middleware  // http通信に関する共通処理を格納するディレクトリ
├─migrate  // データベースのテーブルを作成処理を格納するディレクトリ
├─model  // ユーザー定義型を格納するディレクトリ
│  └─enum  // バックエンドの処理で使用する列挙型を格納するディレクトリ
├─pb  // gRPCのスキーマから自動生成されたコードを格納するディレクトリ
├─proto  // gRPCのスキーマを格納するディレクトリ
├─repository  // データベースのテーブルをCRUD操作する処理を格納するディレクトリ
├─router  // エンドポイントを記述したファイルを格納するディレクトリ
├─server  // gRPC通信に関する処理を格納するディレクトリ
│  ├─interceptor  // gRPC通信に関する共通処理を格納するディレクトリ
│  └─service  // gRPC通信を行うエンドポイントで呼び出される処理を格納するディレクトリ
└─usecase  // ビジネスロジックを格納するディレクトリ
```
# 工夫した点
- フロントエンド
    - チャット画面を表示する際、画面最下部が表示される様に実装しました。
    - チャット画面上部をスクロールした際、25個目のメッセージが画面に表示された時にメッセージを読み込むように実装しました。

- バックエンド
    - httpに関する共通処理をミドルウェアとして実装しました。
    - grpcに関する共通処理をインターセプターとして実装しました。

- インフラ
    - フロントエンド～バックエンド間の通信をリバースプロキシとすることでクライアントPCから直接バックエンドと通信できないようにしました。
# 反省点
- フロントエンド
    - ディレクトリ構成がうまく設計できなかった。
    - サーバーストリーミングのロジックを関数コンポーネントに実装している。
    - エラー処理を実装できていない。
    - バリデーションを実装できていない。

- バックエンド
    - クリーンアーキテクチャを意識して実装したが依存性逆転の法則ができていない。
    - トランザクション処理がusecase層に依存している。（テーブルごとにrepositoryを別けているのだがテーブルを跨いだトランザクション処理がusecase層に依存している）
    - アクティブなグループチャットとグループチャットを開いているメンバーを1つの構造体で管理しているためスケールアウトができない。
    - バリデーションを実装できていない。