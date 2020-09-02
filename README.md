# BloodTemp

- BloodTemp は OCI上で動作する発熱の記録用の Web App
  - いくらコロナ予防っていわれて従業員全員記録紙渡されてもめんどい
  - 最初はシェルスクリプトで記録していた
```
#!/bin/sh

echo "INSERT INTO BLOODTEMP VALUES (localtimestamp,$1);" >/tmp/tmp.sql
echo "QUIT;" >>/tmp/tmp.sql
sqlplus admin/XXXXX@XXX_tp @/tmp/tmp.sql
```
- でもスマホで記録したかった

- go build -o BloodTemp で作成

- まず ADW で table を作成
```
CREATE TABLE "ADMIN"."BLOODTEMP" 
   (	"DATE" DATE NOT NULL ENABLE, 
     	"TEMP" NUMBER NOT NULL ENABLE
   );
```

- port 3001 を開ける（下記参照）

- OCI上で `BloodTemp &` で起動

- PCやスマホからhttp://XXXXXXX:3001/ で接続
  - 体温を入力してリストに追加しよう

- 例：port 3001の開け方
- まず Oracle Linux で下記を実行
```
$ sudo firewall-cmd --list-ports
3000/tcp
$ sudo firewall-cmd --permanent --add-port=3001/tcp
success
$ sudo firewall-cmd --reload
success
sudo firewall-cmd --list-ports
3000/tcp 3001/tcp
```

- つぎにWebのVCN設定を起動
- ネットワーキング> 仮想クラウド・ネットワーク> 仮想クラウド・ネットワークの詳細より 「Security Lists」をクリック
- 右端の...から、「View Security List 詳細」をクリック
- 「Edit All Rules」をクリック
- Ingressルールの一番下にある「+ イングレス・ルールの追加」をクリック
- ソース: 0.0.0.0/0、宛先ポート3001 を追加

- これで外部から 3001 ポートにアクセスできるようになりました
