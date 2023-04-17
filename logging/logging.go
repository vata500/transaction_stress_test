package logging

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func Start(log_path string, startTime time.Time){
		// 파일 열기
		fmt.Println(startTime)
		fmt.Println("시작한 시간입니다.")
		file, err := os.Open(log_path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 파일 스캔
		var batchSentList []string

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if isBatchSentLogLine(line) && isAfterStartTime(line, time.Now()) {
				batchSentList = append(batchSentList, line)
				}
			}

		// 에러 체크
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		// // 결과 출력
		// fmt.Println(batchSentList)

		for _, line := range batchSentList {
			fmt.Println(line)
		}

		// 파일을 생성합니다.
		now := time.Now()
		filename := fmt.Sprintf("batch-%s.log", now.Format("2006-01-02-15-04"))
		logFile, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		
		// 파일에 슬라이스의 모든 값을 쓰기 위한 writer를 생성합니다.
		writer := bufio.NewWriter(logFile)
		fmt.Printf("created %s file\n", filename)
		for _, str := range batchSentList {
			_, err := writer.WriteString(str + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
		// 버퍼를 비워서 파일에 모든 값을 저장합니다.
		writer.Flush()
}

func isBatchSentLogLine(logLine string) bool {
	return strings.Contains(logLine, "BatchPoster: batch sent")
}

func isAfterStartTime(logLine string, startTime time.Time) bool {
	logTime, err := extractLogTime(logLine)
	if err != nil {
		// 시간을 추출할 수 없는 로그 라인인 경우에는 startTime 이후로 기록된 것으로 간주
		return false
	}
	return logTime.After(startTime)
}

func extractLogTime(logLine string) (time.Time, error) {
	const logTimeLayout = "2006-01-02T15:04:05-0700"
	logTimeStr := logLine[0:26] // 로그 라인의 처음 26자리가 시간 정보
	return time.Parse(logTimeLayout, logTimeStr)
}