package jepsen

type initRequest struct {
	Type    string   `json:"type"`
	Id      int      `json:"id"`
	NodeId  string   `json:"node_id"`
	NodeIds []string `json:"node_ids"`
}

type initOkResponse struct {
	Type string `json:"type"`
}

func (initRequest) isBody()    {}
func (initOkResponse) isBody() {}

func (i initRequest) intoReply(id int) body {
	return &initOkResponse{
		Type: "init_ok",
	}
}

func (i initOkResponse) intoReply(id int) body {
	return &i
}
