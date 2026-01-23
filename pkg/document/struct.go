package document

type PropertyDoc struct {
	Name, Type, Description, DefaultValue, ExecSpace string
}
type ParamInfo struct {
	Name, Type, Description string // 설명 필드 추가
}
type MethodDoc struct {
	Name, ReturnType, Description, ExecSpace string
	Params                                   []ParamInfo
}
type HandlerDoc struct {
	Name, EventType, EventVar, Description, ExecSpace, ReturnType string
	EventSenderType                                               string // Type of EventSender (Entity, LocalPlayer, Logic, Self, Model, Service)
	EventSenderValue                                              string // Additional value for Logic and Service types
	Params                                                        []ParamInfo // 핸들러도 파라미터를 가질 수 있으므로 추가
}
type Documentation struct {
	DocType    string // 파일의 종류 (@Logic, @Component 등)
	Properties []PropertyDoc
	Methods    []MethodDoc
	Handlers   []HandlerDoc
}
