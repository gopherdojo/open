package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/yuonoda/gopherdojo-studyroom/kadai3-2/yuonoda/utilities"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type partialContent struct {
	StartByte int
	EndByte   int
	Body      []byte
}

type Downloader struct {
	Url              string
	Size             int
	BatchSize        int
	BatchCount       int
	Content          []byte
	PartialContentCh chan partialContent
	Http             http.Client
}

func (d *Downloader) GetSize() error {
	log.Printf("resource.getSize()\n")

	// HEADでサイズを調べる
	res, err := d.Http.Head(d.Url)
	if err != nil {
		return err
	}

	// データサイズを取得
	header := res.Header
	cl, ok := header["Content-Length"]
	if !ok {
		return errors.New("Content-Length couldn't be found")
	}
	d.Size, err = strconv.Atoi(cl[0])
	if err != nil {
		return err
	}
	return nil

}

func (d *Downloader) GetPartialContent(startByte int, endByte int, ctx context.Context) error {
	log.Printf("resource.getPartialContent(%d, %d)\n", startByte, endByte)
	// Rangeヘッダーを作成
	rangeVal := fmt.Sprintf("bytes=%d-%d", startByte, endByte)

	// リクエストとクライアントの作成
	reader := bytes.NewReader([]byte{})
	req, err := http.NewRequest("GET", d.Url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Range", rangeVal)
	client := &http.Client{}

	res := &http.Response{}
	const retryCount = 3
	for i := 0; i < retryCount; i++ {

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
		log.Printf("sleep:%v\n", sleepTime)
		time.Sleep(sleepTime)
	}

	// 正常系レスポンスでないとき
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return errors.New("status code is not 2xx, got " + res.Status)
	}

	// bodyの取得
	log.Println("start reading")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	pc := partialContent{StartByte: startByte, EndByte: endByte, Body: body}
	d.PartialContentCh <- pc
	return nil
}

func (d *Downloader) GetContent(batchCount int, ctx context.Context) error {
	log.Println("resource.getContent()")

	// コンテンツのデータサイズを取得
	err := d.GetSize()
	if err != nil {
		return err
	}
	log.Printf("d.size: %d\n", d.Size)

	// batchCount分リクエスト
	d.BatchCount = batchCount
	d.BatchSize = int(math.Ceil(float64(d.Size) / float64(d.BatchCount)))
	d.Content = make([]byte, d.Size)
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)
	d.PartialContentCh = make(chan partialContent, d.BatchCount)
	for i := 0; i < d.BatchCount; i++ {

		// 担当する範囲を決定
		startByte := d.BatchSize * i
		endByte := d.BatchSize*(i+1) - 1
		if endByte > d.Size {
			endByte = d.Size
		}

		// レンジごとにリクエスト
		go d.GetPartialContent(startByte, endByte, ctx)
	}

	// リクエスト回数分受け付けてマージ
	var mu sync.Mutex
	d.Content = make([]byte, d.Size)
	for i := 0; i < d.BatchCount; i++ {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case pc := <-d.PartialContentCh:
				mu.Lock()
				utilities.FillByteArr(d.Content[:], pc.StartByte, pc.Body)
				mu.Unlock()
			}
			return nil
		})
	}

	// １リクエストでも失敗すれば終了
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func (d Downloader) Download(ctx context.Context, batchCount int, dwDirPath string) error {
	log.Println("Download")

	// 一時ファイルの作成
	_, filename := filepath.Split(d.Url)
	dwFilePath := dwDirPath + "/" + filename + ".download"
	finishedFilePath := dwDirPath + "/" + filename
	dwFile, err := os.Create(dwFilePath)
	if err != nil {
		return err
	}
	defer os.Remove(dwFilePath)

	// ダウンロード実行
	err = d.GetContent(batchCount, ctx)
	if err != nil {
		return err
	}

	// データの書き込み
	_, err = dwFile.Write(d.Content)
	if err != nil {
		return err
	}
	if err = os.Rename(dwFilePath, finishedFilePath); err != nil {
		return err
	}
	log.Println("download succeeded!")
	return nil
}
