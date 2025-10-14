package models

import "encoding/json"

// StreamPayload คือ interface ที่บังคับว่า model ที่จะส่งเข้า Stream
// ต้องมี method สำหรับแปลงเป็น map
type StreamPayload interface {
	ToMap() (map[string]interface{}, error)
}

type ApprovalJob struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Amount  int    `json:"amount"`
	Details string `json:"details"`
	Status  string `json:"status"`
}

// ToMap คือ method ที่ทำการแปลง ApprovalJob struct ให้เป็น map[string]interface{}
// เพื่อให้เข้ากันได้กับ `go-redis` XAdd command
func (aj *ApprovalJob) ToMap() (map[string]interface{}, error) {
	var inInterface map[string]interface{}

	//แปลง struct เป็น []byte (JSON)
	inrec, err := json.Marshal(aj)
	if err != nil {
		return nil, err
	}
	//แปลง []byte (JSON) กลับมาเป็น map[string]interface{}
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}
	return inInterface, nil
}
