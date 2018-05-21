package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/AnimationMentor/cachedmap"
)

func main() {
	fmt.Printf(`Type things:

- to set: key word word word
- to get: key
- blank line to see stats

`)
	log := logrus.NewEntry(logrus.New())

	cm := cachedmap.NewCachedMap(5*time.Second, 15*time.Second, log)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		words := strings.Fields(s.Text())

		if len(words) == 1 {
			v, ok := cm.Get(words[0])
			log.Infof("v=%#v, ok=%t", v, ok)
		} else if len(words) > 1 {
			cm.Set(words[0], words[1:])
		}
		log.Info(cm)
	}
}
