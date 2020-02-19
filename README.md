# terraform-provider-n0stack-yamlcompiler

## これはなに

[terraform-provier-n0stack](https://github.com/onokatio/terraform-provider-n0stack)で作成したtfファイルを、n0cli doのyaml形式へ変換するツールです。  
terraformのHCLで使用できる、ループや変数展開に対応しているため、複雑なyamlを生成することができます。

100%の互換性はないので、特に問題ない場合はterraformをそのまま使ったほうが楽です。どうしてもterraformからyamlを生成したくなったときのみ利用してください。

# インストール方法

1. [Github Release](https://github.com/onokatio/terraform-provider-n0stack-yamlcompiler/releases)から実行ファイルをダウンロードし、 `terraform-provider-n0stack-yamlcompiler` という名前へ変更する
2. `terraform-provider-n0stack-yamlcompiler` がカレントディレクトリにある状態で、 `terraform init` を行う
3. それ以降、そのディレクトリ以下で`terraform apply`が使えるようになります。

# 使い方

1. `terraform init`を行ったディレクトリへ、変換したいtfファイルを突っ込みます。
2. そこで`terraform apply` を実行すると、

- `n0cli-yaml/Delete`
- `n0cli-yaml/Generate`

というディレクトリが掘られ、その中に `n0cli do` が認識できるyamlが生成されます。

# 注意

- このツールは、terraformのproviderとして作られています。そのため、特に意味はないですがtfstateファイルを生成します。このtfstateファイルは、[terraform-provier-n0stack](https://github.com/onokatio/terraform-provider-n0stack)が生成するtfstateファイルと互換性がありません。  
  すなわち、既に`terraform-provider-n0stack`を使っているディレクトリでこのツールを使うとtfstateが破壊されます。予め変換したいtfファイルを新しいディレクトリへコピーした上、そこでこのツールを実行してください。
- 本来n0stackのImageは、バージョンという概念で複数のブロックストレージを指し示すことができますが、このツールでは完全に対応していません。  
  具体的には、バージョンを使って既に存在するイメージに新たなブロックストレージを追加しようとするtfファイルを受け取った場合、n0cliでエラーが発生するようなyamlを生成してしまいます。
