package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "time"
)

const (
    ENDPOINT     = "https://api.tokyometroapp.jp/api/v2/datapoints"
    CONSUMER_KEY = "取得したアクセストークン"
)

// レスポンスJSONデータ用構造体
type TrainInfomation struct {
    Context                string    `json:"@context"`
    Id                     string    `json:"@id"`
    Type                   string    `json:"@type"`
    Date                   time.Time `json:"dc:date"`
    Valid                  time.Time `json:"dct:valid"`
    Operator               string    `json:"odpt:operator"`
    TimeOfOrigin           time.Time `json:"odpt:timeOfOrigin"`
    Railway                string    `json:"odpt:railway"`
    TrainInformationStatus string    `json:"odpt:trainInformationStatus"`
    TrainInformationText   string    `json:"odpt:trainInformationText"`
}

// レスポンスJSONデータ配列
type TrainInformations []TrainInfomation

func main() {
  // URL生成
  q := map[string]string{
    "rdf:type": "odpt:TrainInformation", // TypeにTrainInformation(運行情報)を設定
    "acl:consumerKey": CONSUMER_KEY, // アクセストークン
  }
  url := fmt.Sprintf("%s?%s", ENDPOINT, buildQuery(q))
  // URLを叩いてデータを取得
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }
  // 取得したデータをJSONデコード
  var trains TrainInformations
  err = json.Unmarshal(body, &trains)
  if err != nil {
    log.Fatal(err)
  }
  // 取得したデータを整形して出力する
  for _, train := range trains {
    // 路線名
    railway := strings.Replace(train.Railway, "odpt.Railway:TokyoMetro.", "", -1)
    // 運行情報を組み立てる
    text := train.TrainInformationText
    if len(train.TrainInformationStatus) > 0 {
      text = fmt.Sprintf("%s (%s)", train.TrainInformationStatus, train.TrainInformationText)
    }
    fmt.Printf("%s: %s [%s]\n", railway, text, train.Date)
  }
}
func buildQuery(q map[string]string) string {
  queries := make([]string, 0)
  for k, v := range q {
    qq := fmt.Sprintf("%s=%s", k, v)
      queries = append(queries, qq)
  }
  return strings.Join(queries, "&")
}
