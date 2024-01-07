package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aclindsa/ofxgo"
	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
)

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func ParseOFX() error {
	f, err := os.Open("./downloads/ned.ofx")
	if err != nil {
		fmt.Printf("could not open OFX file: %v", err) 
		return err
	}
	defer f.Close()

	resp, err := ofxgo.ParseResponse(f)
	if err != nil {
		fmt.Printf("could not parse OFX file: %v", err)
		return err
	}


	// Log the account number and balance for each bank account
	if len(resp.Bank) > 0 {
		bankMessage := resp.Bank[0]
		if stmt, ok := bankMessage.(*ofxgo.StatementResponse); ok {
			// Access the TransactionList
			transactions := stmt.BankTranList
			for _, transaction := range transactions.Transactions {
				fmt.Printf("Transaction type: %s\n", transaction.TrnType)
				fmt.Printf("Transaction date: %s\n", transaction.DtPosted.Format("2006-01-02"))
				fmt.Printf("Transaction amount: %v\n", transaction.TrnAmt)
				fmt.Printf("Transaction ID: %s\n", transaction.FiTID)
				fmt.Printf("Transaction name: %s\n", transaction.Name)
				fmt.Printf("Transaction memo: %s\n", transaction.Memo)
				fmt.Printf("Transaction check number: %s\n", transaction.CheckNum)
				fmt.Printf("Transaction ref number: %s\n", transaction.RefNum)
				fmt.Printf("Transaction sic: %v\n", transaction.SIC)
				fmt.Printf("Transaction mcc: %v\n", transaction.DtPosted.Time)
				fmt.Printf("=======================================================================================\n")
			}
		}
	}

	return nil
}

func cleanUp() error {

	// Remove downloads directory
	// err := os.RemoveAll("./downloads")
	// if err != nil {
	// 	fmt.Printf("could not remove downloads directory: %v", err)
	// 	return err
	// }

	return nil
}

func main() {
	err := godotenv.Load("./.env")
	assertErrorToNilf("Error loading .env file", err)

	var usern string
	var pass string
	var website string
	var waitForLogin string
	var waitForLogout string
	usern = os.Getenv("USERN")
	pass = os.Getenv("PASSWORD")
	website = os.Getenv("WEBSITE")
	waitForLogin = os.Getenv("WEBSITE_LOGIN_WAIT")
	waitForLogout = os.Getenv("WEBSITE_LOGOUT_WAIT")
	
	// Launch Playwright
	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)

	// Launch Browser with UI
	browser, err := pw.Chromium.Launch()
	// browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
	// 	Headless: playwright.Bool(false),
	// })
	assertErrorToNilf("could not launch Chromium: %w", err)

	// Create New Page
	page, err := browser.NewPage()
	assertErrorToNilf("could not create page: %w", err)

	// Goto Website
	_, err = page.Goto(website)
	assertErrorToNilf("could not goto: %w", err)

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select Use Nedbank ID to log in: %v", page.Locator(`[aria-label="Use Nedbank ID to log in"]`).Click())  

	time.Sleep(3 * time.Second) // Wait for 3 seconds


	// Fill in Username
	assertErrorToNilf("could not type: %v", page.Locator("input#username").Fill(usern))

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Fill in Password
	assertErrorToNilf("could not type: %v", page.Locator("input#password").Fill(pass))

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click Login
	assertErrorToNilf("could not press: %v", page.Locator("#log_in").Click())


	//WaitForLogin to complete
	frame := page.MainFrame()
	_ = frame.WaitForURL(waitForLogin)

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	assertErrorToNilf("could not select statement-position: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-landing/div[1]/div/div[2]/div/div/div/div[1]/a`).Click())

	time.Sleep(3 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select statement enquiry tab: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[1]/app-toggle-tab-group/div/div[2]/label`).Click())

	time.Sleep(3 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not open Enquire by dropdown: %v", page.Locator(`.enquireby-options app-enquiry-dropdown .gd-dropdown .discrp-block`).Click())

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select Enquire by dropdown: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[2]/div/app-statements-enquiry/form/div/div[1]/app-enquiry-dropdown/div/ul/li[2]`).Click())

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not open Format dropdown: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[2]/div/app-statements-enquiry/form/div/app-enquiry-dropdown/div/div`).Click())

	time.Sleep(2 * time.Second) // Wait for 2 seconds

	assertErrorToNilf("could not select OFX option: %v", page.Locator(`//*[@id="scroll-page"]/div/div[1]/div/app-landing/section/app-statement-documents-global/div/section/section[2]/div/app-statements-enquiry/form/div/app-enquiry-dropdown/div/ul/li[3]`).Click())

	time.Sleep(2 * time.Second) // Wait for 2 seconds


	//Download
	download, err := page.ExpectDownload(func() error {
		return page.Locator(`#download`).Click()
	})
	assertErrorToNilf("could not download file:  %w", err)

	// Save download to file
	err = download.SaveAs("./downloads/ned.ofx") // Save to current directory
	assertErrorToNilf("could not save download to file: %w", err)

	// Logout
	assertErrorToNilf("could not logout: %v", page.Locator("header .shiftHeader li.logout a").Click())


	// WaitForLogout to complete
	_ = frame.WaitForURL(waitForLogout)

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Parse OFX
	err = ParseOFX()
	assertErrorToNilf("could not parse OFX: %w", err)

	// Clean Up
	err = cleanUp()
	assertErrorToNilf("could not clean up: %w", err)


	assertErrorToNilf("could not close browser: %w", browser.Close())
	assertErrorToNilf("could not stop Playwright: %w", pw.Stop())

	log.Printf("Completed!")
}