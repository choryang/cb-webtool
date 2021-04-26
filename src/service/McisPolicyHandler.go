package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// "os"
	model "github.com/cloud-barista/cb-webtool/src/model"
	util "github.com/cloud-barista/cb-webtool/src/util"
)

//
func GetMcisPolicyList(nameSpaceID string) ([]model.McisPolicyInfo, model.WebStatus) {

	url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/policy/mcis"
	///ns/{nsId}/policy/mcis

	// resp, err := util.CommonHttp(url, nil, http.MethodGet)
	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode

	mcisPolicyList := map[string][]model.McisPolicyInfo{}
	json.NewDecoder(respBody).Decode(&mcisPolicyList)
	fmt.Println(mcisPolicyList["mcisPolicy"])
	log.Println(respBody)
	util.DisplayResponse(resp) // 수신내용 확인

	return mcisPolicyList["mcisPolicy"], model.WebStatus{StatusCode: respStatus}
}

// Get McisPolish Data
func GetMcisPolicyInfoData(nameSpaceID string, mcisID string) (*model.McisPolicyInfo, model.WebStatus) {

	url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/policy/mcis/" + mcisID
	// /ns/{nsId}/policy/mcis/{mcisId}

	resp, err := util.CommonHttpWithoutParam(url, http.MethodGet)

	// defer body.Close()
	mcisPolicyInfo := model.McisPolicyInfo{}
	if err != nil {
		fmt.Println(err)
		return &mcisPolicyInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// util.DisplayResponse(resp) // 수신내용 확인

	respBody := resp.Body
	respStatus := resp.StatusCode

	json.NewDecoder(respBody).Decode(&mcisPolicyInfo)
	fmt.Println(mcisPolicyInfo)

	return &mcisPolicyInfo, model.WebStatus{StatusCode: respStatus}
}

//
func RegMcisPolicy(nameSpaceID string, mcisID string, mcisPolicyInfo *model.McisPolicyInfo) (*model.McisPolicyInfo, model.WebStatus) {

	url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/policy/mcis" + mcisID
	// /ns/{nsId}/policy/mcis/{mcisId}

	pbytes, _ := json.Marshal(mcisPolicyInfo)
	resp, err := util.CommonHttp(url, pbytes, http.MethodPost)

	returnMcisPolicyInfo := model.McisPolicyInfo{}
	returnStatus := model.WebStatus{}

	respBody := resp.Body
	respStatus := resp.StatusCode

	if err != nil {
		fmt.Println(err)
		return &returnMcisPolicyInfo, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	if respStatus != 200 && respStatus != 201 { // 호출은 정상이나, 가져온 결과값이 200, 201아닌 경우 message에 담겨있는 것을 WebStatus에 set
		errorInfo := model.ErrorInfo{}
		json.NewDecoder(respBody).Decode(&errorInfo)
		fmt.Println("respStatus != 200 reason ", errorInfo)
		returnStatus.Message = errorInfo.Message
	} else {
		json.NewDecoder(respBody).Decode(&returnMcisPolicyInfo)
		fmt.Println(returnMcisPolicyInfo)
	}
	returnStatus.StatusCode = respStatus

	return &returnMcisPolicyInfo, returnStatus
}

//
func DelAllMcisPolicy(nameSpaceID string) (io.ReadCloser, model.WebStatus) {
	// buff := bytes.NewBuffer(pbytes)
	url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/policy/mcis"
	// /ns/{nsId}/policy/mcis
	fmt.Println("url : ", url)

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}
	// return body, err
	respBody := resp.Body
	respStatus := resp.StatusCode

	return respBody, model.WebStatus{StatusCode: respStatus}
}

func DelMcisPolicy(nameSpaceID string, mcisID string) (io.ReadCloser, model.WebStatus) {
	// buff := bytes.NewBuffer(pbytes)
	url := util.TUMBLEBUG + "/ns/" + nameSpaceID + "/policy/mcis/" + mcisID
	// /ns/{nsId}/policy/mcis/{mcisId}
	fmt.Println("url : ", url)

	if mcisID == "" {
		return nil, model.WebStatus{StatusCode: 500, Message: "MCIS ID is required"}
	}

	// 경로안에 parameter가 있어 추가 param없이 호출 함.
	resp, err := util.CommonHttp(url, nil, http.MethodDelete)
	if err != nil {
		fmt.Println(err)
		return nil, model.WebStatus{StatusCode: 500, Message: err.Error()}
	}

	respBody := resp.Body
	respStatus := resp.StatusCode
	return respBody, model.WebStatus{StatusCode: respStatus}
}
