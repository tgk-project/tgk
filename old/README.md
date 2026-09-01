# TGK

TGK is a tinygo Keyboard Firmware Supported by tinygo for nrf52xxx .

## アーキテクチャ図

```mermaid
graph TD
    %% TGK actor layers
    subgraph "TGK actor"
        subgraph "loop"
            A[loop]
        end

        subgraph "handler"
            B[tgk manager]
        end

        subgraph "service"
            C[hid]
            D[config]
            E[keyscan]
            F[layer]
            G[split]
            H[remap]
        end

        subgraph "repository"
            I[USB HID]
            J[config]
            K[GPIO]
            L[keymap]
        end
    end

    %% External components
    M[RAW HID]
    N[config.json]

    %% Custom matrix types
    subgraph "custom matrix"
        O[mx]
        P[ec]
    end

    %% Custom keycode types
    subgraph "custom keycode"
        Q[macro]
        R[taphold]
    end

    %% マイコンのピン操作
    S[マイコンのピン操作]

    %% Connections
    A --> B
    B --> C
    B --> D
    B --> E
    B --> F
    B --> G
    B --> H

    C --> I
    D --> J
    E --> K
    F --> L
    G --> L

    J --> N
    K --> S
    K --> O
    K --> P
    L --> Q
    L --> R

    M --> H
    H --> L
```

## TODO

- [ ] hid パッケージに用意されている関数を動的に読み込んで HIDService を構築させる
- [ ] matrix パッケージに用意されている関数を動的に読み込んで keyscan を動的に実行する
- [ ] ユーザ指定の keyscan を config リポジトリに配置して実行できるようにする
