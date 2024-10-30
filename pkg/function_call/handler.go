package function_call

import (
	"log"
)

func HandleFunctionCall(functionCall FunctionCallResponse) {
	log.Printf("Function call: %s , %v", functionCall.Name, functionCall.Arguments)
}
