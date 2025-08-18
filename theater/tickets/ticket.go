package tickets

type Ticket struct {
	Id                  string `json:"id"`
	MobileNumber        string `json:"mobileNumber"`
	NumberOfTickets     int    `json:"numberOfTickets"`
	Name                string `json:"name"`
	PaymentMobileNumber string `json:"paymentMobileNumber"`
	DateTime            string `json:"dateTime"`
	Item                string `json:"item"`
	UniqueId            string `json:"uniqueId"`
	BookedBy            string `json:"bookedBy"`
}
