package cointop

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (ct *Cointop) readAPIKeyFromStdin(name string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter %s API Key: ", name)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(text)
}
