package function_call

type CallEmergencyArguments struct {
	EmergencyType string `json:"emergency_type"`
	Location      string `json:"location"`
}

type EndConversationArguments struct {
	Reason string `json:"reason"`
}

type GetPatientInfoArguments struct {
	PatientName string `json:"patient_name"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Symptoms    string `json:"symptoms"`
	Diagnosis   string `json:"diagnosis"`
	Treatment   string `json:"treatment"`
}
