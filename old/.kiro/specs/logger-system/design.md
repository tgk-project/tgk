# 設計書: ロガーシステム

## 概要

ロガーシステムは、シリアルコンソールやその他の出力先にログを出力するための標準的な方法を提供するように設計されています。このシステムはTGKフレームワークに統合され、keyboard.jsonファイルを通じて設定可能です。システムは異なるログレベル（Debug、Info、Warn、Error）をサポートし、開発者がログの詳細度を制御できるようにします。

## アーキテクチャ

ロガーシステムは、既存のTGKフレームワークと一貫性のあるサービス指向アーキテクチャパターンに従います。以下のコンポーネントで構成されます：

1. **LoggerServiceインターフェース**: ログ操作のための契約を定義します。
2. **loggerService実装**: ログ機能を実装します。
3. **設定統合**: ロガー設定をサポートするために既存の設定システムを拡張します。
4. **TGKマネージャー統合**: ロガーが適切に初期化され、他のコンポーネントからアクセス可能であることを保証します。

## コンポーネントとインターフェース

### LoggerServiceインターフェース

```go
// LoggerService はログサービスのインターフェースです
type LoggerService interface {
    Debug(format string, args ...interface{})
    Info(format string, args ...interface{})
    Warn(format string, args ...interface{})
    Error(format string, args ...interface{})
    SetLevel(level LogLevel)
    SetOutput(writer io.Writer)
    AddOutput(writer io.Writer)
}
```

### loggerService実装

```go
// loggerService はログサービスの実装です
type loggerService struct {
    level   LogLevel
    loggers []*log.Logger
    outputs []io.Writer
    mu      sync.Mutex
}
```

### 設定拡張

`KeyboardConfig`構造体にロガー設定を含めるように拡張します：

```go
type KeyboardConfig struct {
    // 既存のフィールド...
    LogLevel string `json:"log_level,omitempty"`
}
```

## データモデル

### LogLevelタイプ

```go
// LogLevel はログレベルを表す型です
type LogLevel int

const (
    // LogLevelDebug はデバッグレベルです
    LogLevelDebug LogLevel = iota
    // LogLevelInfo は情報レベルです
    LogLevelInfo
    // LogLevelWarn は警告レベルです
    LogLevelWarn
    // LogLevelError はエラーレベルです
    LogLevelError
)

// String はログレベルの文字列表現を返します
func (l LogLevel) String() string {
    switch l {
    case LogLevelDebug:
        return "DEBUG"
    case LogLevelInfo:
        return "INFO"
    case LogLevelWarn:
        return "WARN"
    case LogLevelError:
        return "ERROR"
    default:
        return "UNKNOWN"
    }
}

// ParseLogLevel は文字列をLogLevelに解析します
func ParseLogLevel(level string) LogLevel {
    switch strings.ToLower(level) {
    case "debug":
        return LogLevelDebug
    case "info":
        return LogLevelInfo
    case "warn", "warning":
        return LogLevelWarn
    case "error":
        return LogLevelError
    default:
        return LogLevelInfo // デフォルトはInfoレベル
    }
}
```

## エラー処理

ロガーシステムは以下の方法でエラーを処理します：

1. **設定エラー**: 設定内のログレベルが無効な場合、システムはデフォルトでInfoレベルを使用し、警告をログに記録します。
2. **出力エラー**: 出力への書き込みが失敗した場合、カスケード障害を防ぐためにエラーは静かに無視されますが、可能であれば他の出力にエラーをログとして記録しようとします。
3. **初期化エラー**: ロガーの初期化に失敗した場合、システムはデフォルト設定で標準のGoログパッケージを使用するフォールバックを行います。

## テスト戦略

ロガーシステムは以下のアプローチでテストされます：

1. **ユニットテスト**: ロガーシステムの個々のコンポーネントを分離してテストします。
   - ログレベルの解析をテスト
   - レベルに基づくログフィルタリングをテスト
   - 複数出力の処理をテスト

2. **統合テスト**: ロガーシステムと他のコンポーネントとの統合をテストします。
   - keyboard.jsonからの設定読み込みをテスト
   - TGKマネージャーとの統合をテスト

3. **モックテスト**: モック出力を使用して、ログメッセージが正しくフォーマットされ送信されることを検証します。
   - ログ出力をキャプチャして検証するためのモックio.Writerを作成
   - ログのフォーマットと内容を検証

## 実装詳細

### ロガー初期化

ロガーはTGKマネージャーの初期化プロセス中に初期化されます：

```go
func (m *Manager) Init() error {
    // 最初にロガーを初期化
    m.loggerService = NewLoggerService()
    
    // 設定を読み込み
    if err := m.configRepository.Load(m.configFileName); err != nil {
        m.loggerService.Error("設定の読み込みに失敗しました: %v", err)
        return err
    }
    
    // 設定に基づいてロガーを構成
    config := m.configRepository.GetConfig()
    if config.LogLevel != "" {
        m.loggerService.SetLevel(ParseLogLevel(config.LogLevel))
    }
    
    // 他のサービスを初期化...
    
    return nil
}
```

### 複数出力のサポート

ロガーはio.Writerのスライスを通じて複数の出力をサポートします：

```go
func (l *loggerService) AddOutput(writer io.Writer) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    l.outputs = append(l.outputs, writer)
    l.loggers = append(l.loggers, log.New(writer, "", log.LstdFlags|log.Lshortfile))
}

func (l *loggerService) log(level LogLevel, format string, args ...interface{}) {
    if l.level <= level {
        prefix := fmt.Sprintf("[%s] ", level.String())
        message := fmt.Sprintf(format, args...)
        
        l.mu.Lock()
        defer l.mu.Unlock()
        
        for _, logger := range l.loggers {
            logger.Output(2, prefix+message)
        }
    }
}
```

### デフォルト設定

keyboard.jsonファイルにログレベルが指定されていない場合、システムはデフォルトでInfoレベルを使用します：

```go
func NewLoggerService() LoggerService {
    service := &loggerService{
        level:   LogLevelInfo, // デフォルトはInfoレベル
        outputs: make([]io.Writer, 0),
        loggers: make([]*log.Logger, 0),
    }
    
    // デフォルト出力としてstdoutを追加
    service.AddOutput(os.Stdout)
    
    return service
}
```