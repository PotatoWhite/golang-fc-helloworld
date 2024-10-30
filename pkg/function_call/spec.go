package function_call

func InitialMessage() []Message {
	return []Message{
		{
			Role:    "user",
			Content: "전문의사. 말이많고 친절하고 환자진단과 치료방안 제시를 목적. 10번 안에 대화를 끝내야 해.",
		},
		{
			Role:    "assistant",
			Content: "성별과 나이를 말씀해 주세요.",
		},
	}
}

func GetFunctionSpec() []FunctionSpec {
	return []FunctionSpec{
		{
			Name:        "call_emergency",
			Description: "위급상황",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"emergency_type": map[string]interface{}{
						"type":        "string",
						"description": "응급 상황의 유형",
					},
					"call_number_kor": map[string]interface{}{
						"type":        "string",
						"description": "긴급전화번호",
					},
					"location": map[string]interface{}{
						"type":        "string",
						"description": "응급 상황 발생 위치",
					},
				},
			},
		},
		{
			Name:        "end_conversation",
			Description: "사용자가 대화 종료를 요청",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"reason": map[string]interface{}{
						"type":        "string",
						"description": "대화를 종료하려는 이유",
					},
				},
			},
		},
	}
}
