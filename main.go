package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const outputPath = "history.md"
const fancyFormat = "January 2, 2006"
const pattern = `^\s*[0-9]+\s+`
const skips = `^(sudo\s+)?(mkdir|rm|cd|cp|which|touch|where|make|z|history|open|git add|git commit|go run|git tag|chmod|python main\.py)\s`

const goGet = "go get"
const yarnAdd = "yarn add"
const brewInstall = "brew install"
const pipInstall = "pip install"
const poetryAdd = "poetry add"
const npmInstall = "npm install"
const gitClone = "git clone"
const luarocksInstall = "luarocks install"
const kubectlApply = "kubectl apply"
const helmInstall = "helm install"

func getFancyDate() string {
	currentTime := time.Now()

	return currentTime.Format(fancyFormat)
}

func printGroup(groupname string, list []string, writter *bufio.Writer) {
	writter.WriteString(fmt.Sprintf("## %s\n", groupname))
	writter.WriteString("```bash\n")

	for _, item := range list {
		writter.WriteString(fmt.Sprintf("\t%s\n", item))
	}
	writter.WriteString("```\n\n\n")
}

func main() {
	outputFile, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer outputFile.Close()

	nBytes, nChunks := int64(0), int64(0)

	writter := bufio.NewWriter(outputFile)
	reader := bufio.NewReader(os.Stdin)

	cache := make(map[string]bool)
	groups := map[string][]string{
		goGet:           {},
		yarnAdd:         {},
		brewInstall:     {},
		pipInstall:      {},
		poetryAdd:       {},
		npmInstall:      {},
		gitClone:        {},
		luarocksInstall: {},
		kubectlApply:    {},
		helmInstall:     {},
	}
	buf := make([]byte, 4*1024)

	var leftover []byte
	regex := regexp.MustCompile(pattern)
	skipsR := regexp.MustCompile(skips)

	writter.WriteString(fmt.Sprintf("# %s\n\n", getFancyDate()))
	writter.WriteString("```bash\n")

	for {
		n, err := reader.Read(buf[:cap(buf)])
		buf = append(leftover, buf[:n]...)
		leftover = make([]byte, 0)

		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}

		nChunks++
		nBytes += int64(len(buf))

		stringy := strings.Trim(string(buf), " ")
		parts := strings.Split(stringy, "\n")

		last := len(parts) - 1
		if !strings.HasSuffix(stringy, "\n") {
			leftover = append(leftover, []byte(parts[last])...)
			parts = parts[:last]
		}

		for _, part := range parts {
			part = regex.ReplaceAllLiteralString(part, "")
			if len(part) == 0 {
				continue
			}

			if skipsR.MatchString(part) {
				continue
			}

			_, ok := cache[part]
			if ok {
				continue
			}

			if strings.HasPrefix(part, goGet) {
				groups[goGet] = append(groups[goGet], part)
				continue
			}

			if strings.HasPrefix(part, yarnAdd) {
				groups[yarnAdd] = append(groups[yarnAdd], part)
				continue
			}

			if strings.HasPrefix(part, brewInstall) {
				groups[brewInstall] = append(groups[brewInstall], part)
				continue
			}

			if strings.HasPrefix(part, pipInstall) {
				groups[pipInstall] = append(groups[pipInstall], part)
				continue
			}

			if strings.HasPrefix(part, poetryAdd) {
				groups[poetryAdd] = append(groups[poetryAdd], part)
				continue
			}

			if strings.HasPrefix(part, npmInstall) || strings.HasPrefix(part, "npm i") {
				groups[npmInstall] = append(groups[npmInstall], part)
				continue
			}

			if strings.HasPrefix(part, gitClone) {
				groups[gitClone] = append(groups[gitClone], part)
				continue
			}

			if strings.HasPrefix(part, luarocksInstall) {
				groups[luarocksInstall] = append(groups[luarocksInstall], part)
				continue
			}

			if strings.HasPrefix(part, kubectlApply) {
				groups[kubectlApply] = append(groups[kubectlApply], part)
				continue
			}

			if strings.HasPrefix(part, helmInstall) {
				groups[helmInstall] = append(groups[helmInstall], part)
				continue
			}

			writter.WriteString(fmt.Sprintf("%s \n\n", part))
			cache[part] = true
		}

		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
	}

	writter.WriteString("```\n\n")

	printGroup(goGet, groups[goGet], writter)
	printGroup(yarnAdd, groups[yarnAdd], writter)
	printGroup(brewInstall, groups[brewInstall], writter)
	printGroup(pipInstall, groups[pipInstall], writter)
	printGroup(poetryAdd, groups[poetryAdd], writter)
	printGroup(npmInstall, groups[npmInstall], writter)
	printGroup(gitClone, groups[gitClone], writter)
	printGroup(luarocksInstall, groups[luarocksInstall], writter)
	printGroup(kubectlApply, groups[kubectlApply], writter)
	printGroup(helmInstall, groups[helmInstall], writter)

	if err := writter.Flush(); err != nil {
		log.Fatal(err)
	}

}
