package cardano

type PolicyScript struct {
	Type    string `json:"type"`
	Scripts Script `json:"scripts"`
}

type Script struct {
	PolicyScriptDetails []byte
}

type PolicyScriptComponent struct {
	
}

type PolicyScriptSlot struct {
	Type string `json:"type"`
	Slot string `json:"slot"`
}

type PolicyScriptSig struct {
	Type    string `json:"type"`
	KeyHash string `json:"keyHash"`
}
