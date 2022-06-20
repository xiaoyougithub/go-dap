package dap

func newEvent(event string) Event {
	return Event{
		ProtocolMessage: ProtocolMessage{
			Seq:  0,
			Type: "event",
		},
		Event: event,
	}
}

func newResponse(request RequestMessage) Response {
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

func newErrorResponse(request RequestMessage, message string) *ErrorResponse {
	er := &ErrorResponse{}
	er.Response = newResponse(request)
	er.Success = false
	er.Message = message
	er.Body.Error.Format = message
	return er
}
