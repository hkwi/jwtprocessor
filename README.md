OpenTelemetry collector processor で、resource attribute に内容の jwt を保存します。

現時点では Log について、scopeLogs の部分のみを jwt でsigned.scopeLogs 属性に保存します。
検証時は jwt から取り出した内容と、実際にある内容を比較してください。

# Usage
[opentelemetry collector builder](https://opentelemetry.io/docs/collector/custom-collector/)
で、jwtprocessorを組み込んだotelcolを作ります。
`builder-config.yaml` を参考にビルドしてください。

