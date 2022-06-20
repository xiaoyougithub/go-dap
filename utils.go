package dap

func NewEvent(event string) Event {
	return Event{
		ProtocolMessage: ProtocolMessage{
			Seq:  0,
			Type: "event",
		},
		Event: event,
	}
}

func NewResponse(request RequestMessage) Response {
	return Response{
		ProtocolMessage: ProtocolMessage{
			Seq:  0,
			Type: "response",
		},
		Command:    request.GetRequest().Command,
		RequestSeq: request.GetRequest().Seq,
		Success:    true,
	}
}

func NewErrorResponse(request RequestMessage, message string) *ErrorResponse {
	er := &ErrorResponse{}
	er.Response = NewResponse(request)
	er.Success = false
	er.Message = message
	er.Body.Error.Format = message
	return er
}
