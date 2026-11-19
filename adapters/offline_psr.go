package adapters

type OfflinePSR struct{}

func init() {
	Register(OfflinePSR{})
}

func (OfflinePSR) Code() string {
	return "O"
}

func (OfflinePSR) Name() string {
	return "电子客票二维码"
}

func (OfflinePSR) URL() (pattern string, mockValue interface{}) {
	return "https://p.12306.cn/tservice/qr/offlinePSR/v1?%s==", nil
}

func (OfflinePSR) AlwaysOn() bool {
	return true
}

func (OfflinePSR) ForUnits() bool {
	return false
}

func (OfflinePSR) Info(string) (JSONObject, error) {
	return nil, nil
}

func (OfflinePSR) TrainNo(JSONObject) ([]TrainSchedule, error) {
	return nil, nil
}

func (OfflinePSR) UnitNo(string, JSONObject) (string, error) {
	return "", nil
}
