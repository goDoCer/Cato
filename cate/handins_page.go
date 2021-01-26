package cate

// Gets the handins page for a task and returns the true deadline of the task
// TODO add functionality to this, since this is quite limited
// func getHandinInfo(url string) (time.Time, error) {
// 	url = cateURL + "/" + url
// 	data, err := getPage(url)
// 	if err != nil {
// 		return time.Time{}, errors.New("Couldn't get handins page")
// 	}
// 	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
// 	if err != nil {
// 		return time.Time{}, errors.New("Couldn't parse handins page")
// 	}
// 	sel := doc.Find(":contains('Due')").Parent().Next()
// 	date := sel.Find("[color='blue']").Text()
// 	due := sel.Find("[color='green']").Text()
// 	fmt.Println(string(due), string(data))
// 	return time.Parse("Mon - 08 Feb 2021 (19:00)", date+" "+due)
// }
