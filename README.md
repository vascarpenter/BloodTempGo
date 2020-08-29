# BloodTemp

- BloodTemp は OCI上で動作する Blood Temp 記録用の Web App

- go build -o BloodTemp で作成

- http://XXXXXXX:3001/ で接続 (firewall 3001を開けてあることを前提)

- 例：port 3001の開け方
``
$ sudo firewall-cmd --list-ports
3000/tcp
$ sudo firewall-cmd --permanent --add-port=3001/tcp
success
$ sudo firewall-cmd --reload
success
sudo firewall-cmd --list-ports
3000/tcp 3001/tcp
``
- 上記をOracle Linux上で設定後、WebのVCNから、    ネットワーキング> 仮想クラウド・ネットワーク> 仮想クラウド・ネットワークの詳細より
「Security Lists」をクリック　→右端の...から、「View Security List 詳細」をクリック
→「Edit All Rules」をクリックし、Ingressルールの一番下にある「+ イングレス・ルールの追加」をクリックし、
ソース: 0.0.0.0/0、宛先ポート3001 を追加
