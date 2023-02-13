package downloader

import (
	"fmt"
	"testing"
)

func TestDownload(t *testing.T) {
	downloader, err := NewDownloader()
	if err != nil {
		t.Fatal(err)
	}
	_, err = downloader.Download("https://www.youtube.com/watch?v=LvyHVgocP_0&ab_channel=Yellowcard-Topic")
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsValidURL(t *testing.T) {
	downloader, err := NewDownloader()
	if err != nil {
		t.Fatal(err)
	}
	if downloader.IsValidURL("https://www.youtube.com/watch?v=SB3wAJjJP6c&ab_channel=blink-182-Topic") != true {
		fmt.Println("failed 1")
		t.Fail()
	}
	if downloader.IsValidURL("https://www.youtube.com/watch?v=a_1tA0bpDQs&ab_channel=MyChemicalRomance-Topic") != true {
		fmt.Println("failed 2")
		t.Fail()
	}
	if downloader.IsValidURL("https://youtu.be/a_1tA0bpDQs") != true {
		fmt.Println("failed 3")
		t.Fail()
	}
	if downloader.IsValidURL("youtu.be/a_1tA0bpDQs") != true {
		fmt.Println("failed 4")
		t.Fail()
	}
	//if downloader.isValidURL("www.youtu.be/a_1tA0bpDQs") == true {
	//	fmt.Println("failed 5")
	//	t.Fail()
	//}
	if downloader.IsValidURL("") == true {
		fmt.Println("failed 6")
		t.Fail()
	}
}

func TestDeleteFile(t *testing.T) {
	downloader, err := NewDownloader()
	if err != nil {
		t.Fatal(err)
	}
	downloader.deleteFile("Dysentery Gary")
}
