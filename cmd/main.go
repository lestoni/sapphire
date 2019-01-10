package main

import (
	"fmt"
	"time"

	"github.com/bxcodec/faker"
	"github.com/lestoni/sapphire/pkg/block"
	"github.com/lestoni/sapphire/pkg/node"
)

type Data struct {
	Latitude           float32 `faker:"lat"`
	Longitude          float32 `faker:"long"`
	CreditCardNumber   string  `faker:"cc_number"`
	CreditCardType     string  `faker:"cc_type"`
	Email              string  `faker:"email"`
	IPV4               string  `faker:"ipv4"`
	IPV6               string  `faker:"ipv6"`
	Password           string  `faker:"password"`
	PhoneNumber        string  `faker:"phone_number"`
	MacAddress         string  `faker:"mac_address"`
	URL                string  `faker:"url"`
	UserName           string  `faker:"username"`
	ToolFreeNumber     string  `faker:"tool_free_number"`
	E164PhoneNumber    string  `faker:"e_164_phone_number"`
	TitleMale          string  `faker:"title_male"`
	TitleFemale        string  `faker:"title_female"`
	FirstName          string  `faker:"first_name"`
	FirstNameMale      string  `faker:"first_name_male"`
	FirstNameFemale    string  `faker:"first_name_female"`
	LastName           string  `faker:"last_name"`
	Name               string  `faker:"name"`
	UnixTime           int64   `faker:"unix_time"`
	Date               string  `faker:"date"`
	Time               string  `faker:"time"`
	MonthName          string  `faker:"month_name"`
	Year               string  `faker:"year"`
	DayOfWeek          string  `faker:"day_of_week"`
	DayOfMonth         string  `faker:"day_of_month"`
	Timestamp          string  `faker:"timestamp"`
	Century            string  `faker:"century"`
	TimeZone           string  `faker:"timezone"`
	TimePeriod         string  `faker:"time_period"`
	Word               string  `faker:"word"`
	Sentence           string  `faker:"sentence"`
	Paragraph          string  `faker:"paragraph"`
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"amount"`
	AmountWithCurrency string  `faker:"amount_with_currency"`
}

func main() {

	var nodes []*node.Node
	var blocks []*block.Block

	blk, _ := block.NewRoot()

	blocks = append(blocks, blk)

	var prev string

	fmt.Println("start - ", time.Now().String())

	for i := 0; i < 10000; i++ {
		var nd *node.Node
		if i == 0 {
			nd = node.NewRoot()

		} else {
			nd = node.New(prev)
		}

		data := Data{}
		err := faker.FakeData(&data)
		if err != nil {
			panic(err)
		}

		nd.AddContent(data)
		prev = nd.Identity

		nodes = append(nodes, nd)

		err = blk.AddNode(nd)
		if err != nil {
			fmt.Println(err)

			// Get Previous Block
			prevBlock := blocks[len(blocks)-1]

			// Build Previous Block
			err := prevBlock.Build()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(prevBlock.Identity, " built!")

			blk, err = block.New(prevBlock.Identity)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("New Block-", blk.Identity)

			blocks = append(blocks, blk)
		}

	}

	err := blk.Build()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("end - ", time.Now().String())

	for _, bk := range blocks {
		fmt.Println(bk.Weight, bk.Height, bk.MRoot)
		fmt.Println(bk.Verify(bk.MRoot))
	}

}
