package model

import "encoding/json"

type Profile struct {
	//Url        string
	//Id         string
	Name       string
	Gender     string
	Age        string
	Height     string
	Weight     string
	Income     string
	Marriage   string
	Education  string
	Occupation string
	Hukou      string
	Xingzuo    string
	House      string
	Car        string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o) //将字符串转换为json
	if err != nil {
		return profile, err
	}

	err = json.Unmarshal(s, &profile) //json转为模型
	return profile, err
}
