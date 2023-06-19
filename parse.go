package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"os/exec"
)

func parse(username string, password string) []string {
	log.Println("[x] Parser: starting...")
	checkContainer()
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.bluecosmo.com/customer/account/login/`),
		chromedp.WaitVisible(`#login-form`, chromedp.ByID),
		chromedp.SendKeys(`#email`, username, chromedp.ByID),
		chromedp.SendKeys(`#pass`, password, chromedp.ByID),
		chromedp.Click(`#send2`, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[x] Parser: logged in to BlueCosmo.")
	}
	data := grabAirtimePlansData(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[x] Parser: finished.")
	command := "docker stop headless-shell"
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Printf("[*] Error: %s", err)
	} else {
		log.Printf("[x] Docker: %s", out)
	}
	return data
}

func grabAirtimePlansData(ctx context.Context) []string {
	log.Println("[x] Parser: collecting data...")
	const serviceNumberXPath = `//table[@id='data-airtimeplans']//tbody//tr//td[1]`
	const planNameXPath = `//table[@id='data-airtimeplans']//tbody//tr//td[2]`
	const minutesUsedXPath = `//table[@id='data-airtimeplans']//tbody//tr//td[3]//li`
	const statusXPath = `//table[@id='data-airtimeplans']//tbody//tr//td[5]`
	var serviceNumber, planName, minutesUsed, status string
	_ = chromedp.Run(ctx,
		chromedp.WaitVisible(`.hc-1`, chromedp.ByQuery),
		chromedp.Navigate(`https://www.bluecosmo.com/services/airtimeplans/`),
		chromedp.WaitVisible(`//table[@id='data-airtimeplans']`, chromedp.BySearch),
		chromedp.Text(serviceNumberXPath, &serviceNumber, chromedp.BySearch),
		chromedp.Text(planNameXPath, &planName, chromedp.BySearch),
		chromedp.Text(minutesUsedXPath, &minutesUsed, chromedp.BySearch),
		chromedp.Text(statusXPath, &status, chromedp.BySearch),
	)
	if minutesUsed == "" {
		minutesUsed = "0"
	}
	data := []string{
		serviceNumber,
		planName,
		minutesUsed,
		status,
	}
	return data
}
