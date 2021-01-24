package splitDownload

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type partialContent struct {
	startByte int
	endByte   int
	body      []byte
}

// コンテンツのデータサイズを取得
func getContentSize(url string) (size int, err error) {
	// HEADでサイズを調べる
	res, err := http.Head(url)
	if err != nil {
		return 0, err
	}

	// データサイズを取得
	header := res.Header
	cl, ok := header["Content-Length"]
	if !ok {
		return 0, errors.New("Content-Length couldn't be found")
	}
	size, err = strconv.Atoi(cl[0])
	return
}

//
func fillByteArr(arr []byte, startAt int, partArr []byte) {
	for i := 0; i < len(partArr); i++ {
		globalIndex := i + startAt
		arr[globalIndex] = partArr[i]
	}
}

// 指定範囲のデータを取得する
func getPartialContent(url string, startByte int, endByte int, fileDataCh chan partialContent) error {
	// Rangeヘッダーを作成
	rangeVal := fmt.Sprintf("bytes=%d-%d", startByte, endByte)

	// リクエストとクライアントの作成
	r := bytes.NewReader([]byte{})
	req, err := http.NewRequest("GET", url, r)
	if err != nil {
		return err
	}
	req.Header.Set("Range", rangeVal)
	client := &http.Client{}

	// 3回までリトライする
	res := &http.Response{}
	for i := 0; i < 3; i++ {
		// リクエストの実行
		log.Printf("rangeVal[%d]:%s", i, rangeVal)
		res, err = client.Do(req)
		if err != nil {
			return err
		}

		// ステータスが200系ならループを抜ける
		log.Printf("res.StatusCode:%d\n", res.StatusCode)
		if res.StatusCode >= 200 && res.StatusCode <= 299 {
			break
		}

		// 乱数分スリープ
		rand.Seed(time.Now().UnixNano())
		randFloat := rand.Float64() + 1
		randMs := math.Pow(randFloat, float64(i+1)) * 1000
		sleepTime := time.Duration(randMs) * time.Millisecond
		time.Sleep(sleepTime)
	}

	// 正常系レスポンスでないとき
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return errors.New("status code is not 2xx, got " + res.Status)
	}

	// bodyの取得
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	pc := partialContent{body: body, startByte: startByte, endByte: endByte}
	fileDataCh <- pc
	return nil
}

func Run(url string, splitCount int) {
	log.Println("Run")

	// コンテンツのデータサイズを取得
	size, err := getContentSize(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("size: %d\n", size)

	// ファイルの作成
	_, filename := filepath.Split(url)
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dwFilePath := homedir + "/Downloads/" + filename + ".download"
	log.Println(dwFilePath)
	dwFile, err := os.Create(dwFilePath)
	if err != nil {
		os.Remove(dwFilePath)
		log.Fatal(err)
	}

	// 1MBごとにアクセス
	singleSize := int(math.Ceil(float64(size) / float64(splitCount)))
	count := int(math.Ceil(float64(size) / float64(singleSize)))
	log.Printf("count: %d\n", count)
	pcch := make(chan partialContent, count)
	var eg errgroup.Group
	for i := 0; i < count; i++ {

		// 担当する範囲を決定
		startByte := singleSize * i
		endByte := singleSize*(i+1) - 1
		if endByte > size {
			endByte = size
		}

		// レンジごとにリクエスト
		eg.Go(func() error {
			return getPartialContent(url, startByte, endByte, pcch)
		})
	}

	// １リクエストでも失敗すれば終了
	if err := eg.Wait(); err != nil {
		os.Remove(dwFilePath)
		log.Fatal(err)
	}

	// 一つのバイト列にマージ
	fileData := make([]byte, size)
	for i := 0; i < count; i++ {
		pc := <-pcch
		fillByteArr(fileData[:], pc.startByte, pc.body)
	}

	// データの書き込み
	_, err = dwFile.Write(fileData)
	if err != nil {
		os.Remove(dwFilePath)
		log.Fatal(err)
	}
	os.Rename(dwFilePath, strings.Trim(dwFilePath, ".download"))

	log.Println("download succeeded!")
}
