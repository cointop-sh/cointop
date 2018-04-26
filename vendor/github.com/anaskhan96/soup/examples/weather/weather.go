package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

func main() {
	fmt.Printf("Enter the name of the city : ")
	city, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	city = city[:len(city)-1]
	cityInURL := strings.Join(strings.Split(city, " "), "+")
	url := "https://www.bing.com/search?q=weather+" + cityInURL
	resp, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(resp)
	grid := doc.Find("div", "class", "b_antiTopBleed b_antiSideBleed b_antiBottomBleed")
	heading := grid.Find("div", "class", "wtr_titleCtrn").Find("div").Text()
	conditions := grid.Find("div", "class", "wtr_condition")
	primaryCondition := conditions.Find("div")
	secondaryCondition := primaryCondition.FindNextElementSibling()
	temp := primaryCondition.Find("div", "class", "wtr_condiTemp").Find("div").Text()
	others := primaryCondition.Find("div", "class", "wtr_condiAttribs").FindAll("div")
	caption := secondaryCondition.Find("div").Text()
	fmt.Println("City Name : " + heading)
	fmt.Println("Temperature : " + temp + "˚C")
	for _, i := range others {
		fmt.Println(i.Text())
	}
	fmt.Println(caption)
}
