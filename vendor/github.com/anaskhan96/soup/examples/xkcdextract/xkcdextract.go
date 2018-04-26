package main

import (
	"fmt"

	"github.com/anaskhan96/soup"
)

func main() {
	fmt.Println("Enter the xkcd comic number :")
	var num int
	fmt.Scanf("%d", &num)
	url := fmt.Sprintf("https://xkcd.com/%d", num)
	resp, _ := soup.Get(url)
	doc := soup.HTMLParse(resp)
	title := doc.Find("div", "id", "ctitle").Text()
	fmt.Println("Title of the comic :", title)
	comicImg := doc.Find("div", "id", "comic").Find("img")
	fmt.Println("Source of the image :", comicImg.Attrs()["src"])
	fmt.Println("Underlying text of the image :", comicImg.Attrs()["title"])
}

/* --- Console I/O ---
Enter the xkcd comic number :
353
Title of the comic : Python
Source of the image : //imgs.xkcd.com/comics/python.png
Underlying text of the image : I wrote 20 short programs in Python yesterday.  It was wonderful.  Perl, I'm leaving you.
*/
