package controller

import (
	// "encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-webtool/src/db"

	// "github.com/cloud-barista/cb-webtool/src/model/tumblebug"
	"github.com/cloud-barista/cb-webtool/src/model"
	tbcommon "github.com/cloud-barista/cb-webtool/src/model/tumblebug/common"
	tbmcir "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcir"
	webtool "github.com/cloud-barista/cb-webtool/src/model/webtool"

	// tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"

	service "github.com/cloud-barista/cb-webtool/src/service"

	util "github.com/cloud-barista/cb-webtool/src/util"

	"github.com/labstack/echo"
	// "io/ioutil"
	"log"
	"net/http"

	//"github.com/davecgh/go-spew/spew"
	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
)

func ResourceBoard(c echo.Context) error {
	fmt.Println("=========== ResourceBoard start ==============")
	comURL := service.GetCommonURL()
	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
		nameSpace := service.GetNameSpaceToString(c)
		fmt.Println("Namespace : ", nameSpace)
		return c.Render(http.StatusOK, "ResourceBoard.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": nameSpace,
			"comURL":    comURL,
		})
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// func VpcListForm(c echo.Context) error {
func VpcMngForm(c echo.Context) error {
	fmt.Println("VpcMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}
	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")
	vNetInfoList, respStatus := service.GetVnetListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return echotemplate.Render(c, http.StatusOK,
			"setting/resources/NetworkMng", // 파일명
			map[string]interface{}{
				"LoginInfo":     loginInfo,
				"CloudOSList":   cloudOsList,
				"NameSpaceList": nsList,
				"VNetList":      vNetInfoList,
				"status":        respStatus.StatusCode,
				"error":         respStatus.Message,
			})
	}
	log.Println("VNetList", vNetInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/NetworkMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			"VNetList":      vNetInfoList,
			"status":        respStatus.StatusCode,
		})
}

func GetVpcList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")
	if optionParam == "id" {
		vNetInfoList, respStatus := service.GetVnetListByID(defaultNameSpaceID, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"VNetList":           vNetInfoList,
		})
	} else {
		vNetInfoList, respStatus := service.GetVnetListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"VNetList":           vNetInfoList,
		})
	}

}

// Vpc 상세정보
func GetVpcData(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVNetID := c.Param("vNetID")
	vNetInfo, respStatus := service.GetVpcData(defaultNameSpaceID, paramVNetID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "success",
		"status":   respStatus,
		"VNetInfo": vNetInfo,
	})
}

// Vpc 등록 :
func VpcRegProc(c echo.Context) error {
	log.Println("VpcRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	vNetRegInfo := new(tbmcir.TbVNetReq)
	if err := c.Bind(vNetRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	// log.Println(vNetRegInfo)
	resultVNetInfo, respStatus := service.RegVpc(defaultNameSpaceID, vNetRegInfo)
	// respBody, respStatus := service.RegVpc(defaultNameSpaceID, vNetRegInfo)
	// fmt.Println("=============respStatus =============", respStatus)
	// fmt.Println("=============respBody ===============", respBody)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// return c.JSON(http.StatusOK, map[string]interface{}{
		// 	"message": respStatus.Message,
		// 	"status":  respStatus.StatusCode,
		// })

		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "success",
		"status":   respStatus.StatusCode,
		"VNetInfo": resultVNetInfo,
	})
}

// 삭제
func VpcDelProc(c echo.Context) error {
	log.Println("ConnectionConfigDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVNetID := c.Param("vNetID")

	respMessage, respStatus := service.DelVpc(defaultNameSpaceID, paramVNetID)
	fmt.Println("=============respMessage =============", respMessage)
	log.Println("respStatus : ", respStatus)
	log.Println("respMessage : ", respMessage)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respMessage.StatusCode,
	})
}

// SecurityGroup 관리 화면
func SecirityGroupMngForm(c echo.Context) error {
	fmt.Println("SecirityGroupMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}
	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")

	securityGroupInfoList, respStatus := service.GetSecurityGroupListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return echotemplate.Render(c, http.StatusOK,
			"setting/resources/SecurityGroupMng", // 파일명
			map[string]interface{}{
				"LoginInfo":         loginInfo,
				"CloudOSList":       cloudOsList,
				"NameSpaceList":     nsList,
				"SecurityGroupList": securityGroupInfoList,
				"status":            respStatus.StatusCode,
				"error":             respStatus.Message,
			})

	}
	log.Println("securityGroupInfoList", securityGroupInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/SecurityGroupMng", // 파일명
		map[string]interface{}{
			"LoginInfo":         loginInfo,
			"CloudOSList":       cloudOsList,
			"NameSpaceList":     nsList,
			"SecurityGroupList": securityGroupInfoList,
			"status":            respStatus.StatusCode,
		})
}

// SecurityGroup 목록
func GetSecirityGroupList(c echo.Context) error {
	log.Println("GetSecirityGroupList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")

	if optionParam == "id" {
		securityGroupInfoList, respStatus := service.GetSecurityGroupListByOptionID(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"SecurityGroupList":  securityGroupInfoList,
		})
	} else {
		securityGroupInfoList, respStatus := service.GetSecurityGroupListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"SecurityGroupList":  securityGroupInfoList,
		})
	}

}

// 상세정보
func GetSecirityGroupData(c echo.Context) error {
	log.Println("GetSecirityGroupData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSecurityGroupID := c.Param("securityGroupID")
	securityGroupInfo, respStatus := service.GetSecurityGroupData(defaultNameSpaceID, paramSecurityGroupID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "success",
		"status":            respStatus,
		"SecurityGroupInfo": securityGroupInfo,
	})
}

// 등록 :
func SecirityGroupRegProc(c echo.Context) error {
	log.Println("SecirityGroupRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	securityGroupRegInfo := new(tbmcir.TbSecurityGroupReq)
	if err := c.Bind(securityGroupRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultSecurityGroupInfo, respStatus := service.RegSecurityGroup(defaultNameSpaceID, securityGroupRegInfo)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "success",
		"status":            respStatus.StatusCode,
		"SecurityGroupInfo": resultSecurityGroupInfo,
	})
}

// SecurityGroup 삭제
func SecirityGroupDelProc(c echo.Context) error {
	log.Println("SecirityGroupDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSecurityGroupID := c.Param("securityGroupID")

	// respBody, respStatus := service.DelSecurityGroup(defaultNameSpaceID, paramSecurityGroupID)
	respMessage, respStatus := service.DelSecurityGroup(defaultNameSpaceID, paramSecurityGroupID)
	fmt.Println("=============respMessage =============", respMessage)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respStatus.StatusCode,
	})
}

// security group rule 추가
func FirewallRegProc(c echo.Context) error {
	log.Println("SecirityGroupRulesRegProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	firewallRuleReq := new(tbmcir.TbFirewallRulesWrapper)
	if err := c.Bind(firewallRuleReq); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	paramSecurityGroupID := c.Param("securityGroupID")

	resultSecurityGroupInfo, respStatus := service.RegFirewallRules(defaultNameSpaceID, paramSecurityGroupID, firewallRuleReq)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "success",
		"status":            respStatus.StatusCode,
		"securityGroupInfo": resultSecurityGroupInfo,
	})
}

// firewall rule 삭제 :
func FirewallDelProc(c echo.Context) error {
	log.Println("FirewallDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSecurityGroupID := c.Param("securityGroupID")

	firewallRuleReq := new(tbmcir.TbFirewallRulesWrapper)
	if err := c.Bind(firewallRuleReq); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultSecurityGroupInfo, respStatus := service.DelFirewallRules(defaultNameSpaceID, paramSecurityGroupID, firewallRuleReq)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "success",
		"status":            respStatus.StatusCode,
		"securityGroupInfo": resultSecurityGroupInfo,
	})
}

// SshKey 등록 form
func SshKeyMngForm(c echo.Context) error {
	fmt.Println("SshKeyMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	log.Println(" defaultNameSpaceID  ", defaultNameSpaceID)
	if defaultNameSpaceID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/main?namespace=unknown")
	}

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}

	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")
	sshKeyInfoList, respStatus := service.GetSshKeyInfoListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return echotemplate.Render(c, http.StatusOK,
			"setting/resources/SshKeyMng", // 파일명
			map[string]interface{}{
				"LoginInfo":     loginInfo,
				"CloudOSList":   cloudOsList,
				"NameSpaceList": nsList,
				"SshKeyList":    sshKeyInfoList,
				"status":        respStatus.StatusCode,
				"error":         respStatus.Message,
			})
	}
	log.Println("sshKeyInfoList", sshKeyInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/SshKeyMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			"SshKeyList":    sshKeyInfoList,
			"status":        respStatus.StatusCode,
		})
}

func GetSshKeyList(c echo.Context) error {
	log.Println("GetSshKeyList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")
	if optionParam == "id" {
		sshKeyInfoList, respStatus := service.GetSshKeyInfoListByID(defaultNameSpaceID, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"SshKeyList":         sshKeyInfoList,
		})
	} else {
		sshKeyInfoList, respStatus := service.GetSshKeyInfoListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"SshKeyList":         sshKeyInfoList,
		})
	}
}

// SSHKey 상세정보
func GetSshKeyData(c echo.Context) error {
	log.Println("GetSshKeyData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSshKey := c.Param("sshKeyID")
	sshKeyInfo, respStatus := service.GetSshKeyData(defaultNameSpaceID, paramSshKey)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus,
		"SshKeyInfo": sshKeyInfo,
	})
}

// SSHKey 등록 :
func SshKeyRegProc(c echo.Context) error {
	log.Println("SshKeyRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	sshKeyRegInfo := new(tbmcir.TbSshKeyReq)
	if err := c.Bind(sshKeyRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultSshKeyInfo, respStatus := service.RegSshKey(defaultNameSpaceID, sshKeyRegInfo)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"SshKeyInfo": resultSshKeyInfo,
	})
}

// 삭제
func SshKeyDelProc(c echo.Context) error {
	log.Println("SshKeyDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSshKeyID := c.Param("sshKeyID")

	//respBody, respStatus := service.DelSshKey(defaultNameSpaceID, paramSshKeyID)
	respMessage, respStatus := service.DelSshKey(defaultNameSpaceID, paramSshKeyID)
	fmt.Println("=============respMessage =============", respMessage)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respStatus.StatusCode,
	})
}

func SshKeyUpdateProc(c echo.Context) error {
	log.Println("SshKeyUpdateProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSshKeyID := c.Param("sshKeyID")

	sshKeyInfo := new(tbmcir.TbSshKeyInfo)
	if err := c.Bind(sshKeyInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultSshKeyInfo, respStatus := service.UpdateSshKey(defaultNameSpaceID, paramSshKeyID, sshKeyInfo)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"SshKeyInfo": resultSshKeyInfo,
	})
}

// VirtualMachine Image 등록 form
func VirtualMachineImageMngForm(c echo.Context) error {
	fmt.Println("VirtualMachineImageMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}
	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")

	virtualMachineImageInfoList, respStatus := service.GetVirtualMachineImageInfoListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return echotemplate.Render(c, http.StatusOK,
			"setting/resources/VirtualMachineImageMng", // 파일명
			map[string]interface{}{
				"LoginInfo":               loginInfo,
				"CloudOSList":             cloudOsList,
				"NameSpaceList":           nsList,
				"VirtualMachineImageList": virtualMachineImageInfoList,
				"status":                  respStatus.StatusCode,
				"error":                   respStatus.Message,
			})
	}
	log.Println("VirtualMachineImageInfoList", virtualMachineImageInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/VirtualMachineImageMng", // 파일명
		map[string]interface{}{
			"LoginInfo":               loginInfo,
			"CloudOSList":             cloudOsList,
			"NameSpaceList":           nsList,
			"VirtualMachineImageList": virtualMachineImageInfoList,
			"status":                  respStatus.StatusCode,
		})
}

// 해당 namespace에 등록된 Spec목록 조회
func GetVirtualMachineImageList(c echo.Context) error {
	log.Println("GetVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")

	if optionParam == "id" {
		virtualMachineImageInfoList, respStatus := service.GetVirtualMachineImageInfoListByID(defaultNameSpaceID, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":                 "success",
			"status":                  respStatus.StatusCode,
			"DefaultNameSpaceID":      defaultNameSpaceID,
			"VirtualMachineImageList": virtualMachineImageInfoList,
		})
	} else {
		virtualMachineImageInfoList, respStatus := service.GetVirtualMachineImageInfoListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":                 "success",
			"status":                  respStatus.StatusCode,
			"DefaultNameSpaceID":      defaultNameSpaceID,
			"VirtualMachineImageList": virtualMachineImageInfoList,
		})
	}

}

// VirtualMachineImage 상세정보
func GetVirtualMachineImageData(c echo.Context) error {
	log.Println("GetVirtualMachineImageData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVirtualMachineImage := c.Param("imageID")
	virtualMachineImageInfo, respStatus := service.GetVirtualMachineImageData(defaultNameSpaceID, paramVirtualMachineImage)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus,
		"VirtualMachineImageInfo": virtualMachineImageInfo,
	})
}

// VirtualMachineImage 등록 :
func VirtualMachineImageRegProc(c echo.Context) error {
	log.Println("VirtualMachineImageRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// virtualMachineImageRegInfo := new(tbmcir.TbImageReq)

	resultVirtualMachineImageInfo := new(tbmcir.TbImageInfo)
	respStatus := model.WebStatus{}

	paramVirtualMachineImageRegistType := c.Param("registeringMethod") // registeringMethod = registerWithId 또는 registerWithInfo
	// registeringMethod 에 따라 request Object가 달라짐.
	if paramVirtualMachineImageRegistType == "registerWithInfo" {
		virtualMachineImageRegInfo := new(tbmcir.TbImageInfo)
		if err := c.Bind(virtualMachineImageRegInfo); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "fail",
				"status":  "fail",
			})
		}
		resultVirtualMachineImageInfo, respStatus = service.RegVirtualMachineImageWithInfo(defaultNameSpaceID, paramVirtualMachineImageRegistType, virtualMachineImageRegInfo)
	} else {
		virtualMachineImageRegInfo := new(tbmcir.TbImageReq)
		if err := c.Bind(virtualMachineImageRegInfo); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "fail",
				"status":  "fail",
			})
		}
		resultVirtualMachineImageInfo, respStatus = service.RegVirtualMachineImage(defaultNameSpaceID, paramVirtualMachineImageRegistType, virtualMachineImageRegInfo)
	}
	// resultVirtualMachineImageInfo, respStatus := service.RegVirtualMachineImage(defaultNameSpaceID, paramVirtualMachineImageRegistType, virtualMachineImageRegInfo)

	// todo : return message 조치 필요. 중복 등 에러났을 때 message 표시가 제대로 되지 않음
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	// respBody := resp.Body
	// respStatusCode := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus.StatusCode,
		"VirtualMachineImageInfo": resultVirtualMachineImageInfo,
	})
}

// 삭제
func VirtualMachineImageDelProc(c echo.Context) error {
	log.Println("VirtualMachineImageDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVirtualMachineImageID := c.Param("imageID")

	// respBody, respStatus := service.DelVirtualMachineImage(defaultNameSpaceID, paramVirtualMachineImageID)
	respMessage, respStatus := service.DelVirtualMachineImage(defaultNameSpaceID, paramVirtualMachineImageID)
	fmt.Println("=============respMessage =============", respMessage)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respStatus.StatusCode,
	})
}

// 해당 namespace의 모든 VMImage 삭제
func AllVirtualMachineImageDelProc(c echo.Context) error {
	log.Println("AllVirtualMachineImageDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramNameSpaceID := c.Param("nameSpaceID")

	// 해당 Namespace의 모든 Image가 삭제 되므로 선택한 namespace와 defaultNamespace가 같아야 삭제가능
	if defaultNameSpaceID != paramNameSpaceID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "현재 namespace만 삭제 가능합니다.",
			"status":  4040,
		})
	}

	//respBody, respStatus := service.DelAllVirtualMachineImage(defaultNameSpaceID)
	respMessage, respStatus := service.DelAllVirtualMachineImage(defaultNameSpaceID)
	fmt.Println("=============respMessage =============", respMessage)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respStatus.StatusCode,
	})
}

// connection에 해당하는 machine image 목록
// connection이 resion별로 생성되므로 결국 해당 provider의 resion 내 vm목록을 가져옴
// deprecated
// func LookupVirtualMachineImageList(c echo.Context) error {
// 	log.Println("GetVirtualMachineImageList : ")
// 	loginInfo := service.CallLoginInfo(c)
// 	if loginInfo.UserID == "" {
// 		return c.Redirect(http.StatusTemporaryRedirect, "/login")
// 	}

// 	// paramConnectionName := c.Param("connectionName")
// 	paramConnectionName := c.QueryParam("connectionName")

// 	log.Println("paramConnectionName : ", paramConnectionName)
// 	virtualMachineImageInfoList, respStatus := service.LookupVirtualMachineImageList(paramConnectionName)
// 	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
// 		return c.JSON(respStatus.StatusCode, map[string]interface{}{
// 			"error":  respStatus.Message,
// 			"status": respStatus.StatusCode,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"message": "success",
// 		"status":  respStatus.StatusCode,
// 		// "DefaultNameSpaceID": defaultNameSpaceID,
// 		"VirtualMachineImageList": virtualMachineImageInfoList,
// 	})
// }

// 해당 connection( provider, region ) 에서 사용가능한 image목록 조회 : 등록시 사용하므로 오래걸려도 기다려야 함.
func LookupCspVirtualMachineImageList(c echo.Context) error {
	log.Println("LookupCspVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// paramConnectionName := c.Param("connectionName")
	paramConnectionName := c.QueryParam("connectionName")

	log.Println("paramConnectionName : ", paramConnectionName)
	virtualMachineImageInfoList, respStatus := service.LookupVirtualMachineImageList(paramConnectionName)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		// "DefaultNameSpaceID": defaultNameSpaceID,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// lookupImage 상세정보
func LookupVirtualMachineImageData(c echo.Context) error {
	log.Println("LookupVirtualMachineImageData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	restLookupImageRequest := new(tbmcir.RestLookupImageRequest)
	if err := c.Bind(restLookupImageRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	// paramVirtualMachineImage := c.Param("imageID")
	// virtualMachineImageInfo, respStatus := service.LookupVirtualMachineImageData(paramVirtualMachineImage)
	virtualMachineImageInfo, respStatus := service.LookupVirtualMachineImageData(restLookupImageRequest)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus,
		"VirtualMachineImageInfo": virtualMachineImageInfo,
	})
}

// lookupImage 상세정보
func SearchVirtualMachineImageList(c echo.Context) error {
	log.Println("SearchVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	restSearchImageRequest := new(tbmcir.RestSearchImageRequest)
	if err := c.Bind(restSearchImageRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	virtualMachineImageInfoList, respStatus := service.SearchVirtualMachineImageList(defaultNameSpaceID, restSearchImageRequest)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// TODO : Fetch 의 의미 파악
func FetchVirtualMachineImageList(c echo.Context) error {
	log.Println("FetchVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	virtualMachineImageInfoList, respStatus := service.FetchVirtualMachineImageList(defaultNameSpaceID)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		// "DefaultNameSpaceID": defaultNameSpaceID,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// VMSpecMng 등록 form
func VmSpecMngForm(c echo.Context) error {
	fmt.Println("VmSpecMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}
	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")

	vmSpecInfoList, respStatus := service.GetVmSpecInfoListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return echotemplate.Render(c, http.StatusOK,
			"setting/resources/VirtualMachineSpecMng", // 파일명
			map[string]interface{}{
				"LoginInfo":     loginInfo,
				"CloudOSList":   cloudOsList,
				"NameSpaceList": nsList,
				"VmSpecList":    vmSpecInfoList,
				"status":        respStatus.StatusCode,
				"error":         respStatus.Message,
			})
	}
	log.Println("VmSpecInfoList", vmSpecInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/VirtualMachineSpecMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			"VmSpecList":    vmSpecInfoList,
			"status":        respStatus.StatusCode,
		})
}

func GetVmSpecList(c echo.Context) error {
	log.Println("GetVmSpecList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")

	if optionParam == "id" {
		vmSpecInfoList, respStatus := service.GetVmSpecInfoListByID(defaultNameSpaceID, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"VmSpecList":         vmSpecInfoList,
		})
	} else {
		vmSpecInfoList, respStatus := service.GetVmSpecInfoListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"DefaultNameSpaceID": defaultNameSpaceID,
			"VmSpecList":         vmSpecInfoList,
		})
	}
}

// VMSpec 상세정보
func GetVmSpecData(c echo.Context) error {
	log.Println("GetVmSpecData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVMSpecID := c.Param("vmSpecID")
	vmSpecInfo, respStatus := service.GetVmSpecInfoData(defaultNameSpaceID, paramVMSpecID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"VmSpec":  vmSpecInfo,
	})
}

// VMSpec 등록 :
func VmSpecRegProc(c echo.Context) error {
	log.Println("VMSpecRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	resultVirtualMachineSpecInfo := new(tbmcir.TbSpecInfo)
	respStatus := model.WebStatus{}

	paramVMSpecregisteringMethod := c.Param("specregisteringMethod") // registerWithInfo or Else(간단등록인 경우 param설정 필요 X)
	if paramVMSpecregisteringMethod == "registerWithInfo" {
		vmSpecRegInfo := new(tbmcir.TbSpecInfo)
		if err := c.Bind(vmSpecRegInfo); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "fail",
				"status":  "fail",
			})
		}
		resultVirtualMachineSpecInfo, respStatus = service.RegVmSpecWithInfo(defaultNameSpaceID, paramVMSpecregisteringMethod, vmSpecRegInfo)
	} else {
		vmSpecRegInfo := new(tbmcir.TbSpecReq)
		if err := c.Bind(vmSpecRegInfo); err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "fail",
				"status":  "fail",
			})
		}
		resultVirtualMachineSpecInfo, respStatus = service.RegVmSpec(defaultNameSpaceID, paramVMSpecregisteringMethod, vmSpecRegInfo)
	}

	// resultVmSpecInfo, respStatus := service.RegVmSpec(defaultNameSpaceID, paramVMSpecregisteringMethod, vmSpecRegInfo)
	// todo : return message 조치 필요. 중복 등 에러났을 때 message 표시가 제대로 되지 않음
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		// 호출은 정상: http.StatusOK, 결과는 정상이 아님. (statusCode != 200,201)
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	// respBody := resp.Body
	// respStatusCode := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		"VMSpec":  resultVirtualMachineSpecInfo,
	})
}

// 삭제
func VmSpecDelProc(c echo.Context) error {
	log.Println("VmSpecDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVMSpecID := c.Param("vmSpecID")

	//respBody, respStatus := service.DelVMSpec(defaultNameSpaceID, paramVMSpecID)
	respMessage, respStatus := service.DelVMSpec(defaultNameSpaceID, paramVMSpecID)
	fmt.Println("=============respMessage =============", respMessage)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respStatus.StatusCode,
	})
}

// lookupImage 목록
func LookupVmSpecList(c echo.Context) error {
	log.Println("LookupVmSpecList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// paramConnectionName := c.Param("connectionName")
	connectionName := new(tbcommon.TbConnectionName)
	if err := c.Bind(connectionName); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	cspVmSpecInfoList, respStatus := service.LookupVmSpecInfoList(connectionName)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"CspVmSpecList":      cspVmSpecInfoList,
	})
}

// lookupImage 상세정보
func LookupVmSpecData(c echo.Context) error {
	log.Println("LookupVmSpecData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	restLookupSpecRequest := new(tbmcir.RestLookupSpecRequest)
	if err := c.Bind(restLookupSpecRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	vmSpecInfo, respStatus := service.LookupVmSpecInfoData(restLookupSpecRequest)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"VmSpec":  vmSpecInfo,
	})
}

// Fetch는 Spider에 있는 VM Image정보를 Tumblebug으로 가져오는 작업. 시간이 오래걸리므로 이전에는 전체 Image목록을 가져왔으나 결과만 return하는 것으로 변경 됨.
func FetchVmSpecList(c echo.Context) error {
	log.Println("FetchVMSpecList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	taskKey := defaultNameSpaceID + "||" + "VMSpec" + "||" + "Fetch"
	//func StoreWebsocketMessage(taskType string, taskKey string, lifeCycle string, requestStatus string, c echo.Context) {
	service.StoreWebsocketMessage(util.TASK_TYPE_VMSPEC, taskKey, util.VMSPEC_LIFECYCLE_CREATE, util.TASK_STATUS_REQUEST, c) // session에 작업내용 저장

	// vmSpecInfoList, respStatus := service.FetchVmSpecInfoList(defaultNameSpaceID)
	go service.FetchVmSpecInfoListByAsync(defaultNameSpaceID, c)
	// if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
	// 	return c.JSON(respStatus.StatusCode, map[string]interface{}{
	// 		"error":  respStatus.Message,
	// 		"status": respStatus.StatusCode,
	// 	})
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  "200",
		// "DefaultNameSpaceID": defaultNameSpaceID,
		// "VmSpec": vmSpecInfoList,
	})
}

// resourcesGroup.PUT("/vmspec/put/:specID", controller.VMSpecPutProc)	// RegProc _ SshKey 같이 앞으로 넘길까
// resourcesGroup.POST("/vmspec/filterspecs", controller.FilterVMSpecList)
// resourcesGroup.POST("/vmspec/filterspecsbyrange", controller.FilterVMSpecListByRange)

// Spec Range search
func FilterVmSpecListByRange(c echo.Context) error {
	log.Println("FilterVmSpecListByRange ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// params := make(map[string]string)
	// _ := c.Bind(&params)
	// fmt.Println(params["connectionName"])
	// fmt.Println(params)
	// fmt.Println("-----------")
	// vmSpecRange := new(tumblebug.VmSpecRangeReqInfo)
	// connectionName := new(tumblebug.TbConnectionName)
	// if err := c.Bind(connectionName); err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"message": "fail",
	// 		"status":  "fail",
	// 	})
	// }
	// fmt.Println(connectionName)
	// fmt.Println("ConnectionName=", connectionName)
	//message=Unmarshal type error: expected=int, got=string, field=maxTotalStorageTiB.min, offset=65
	vmSpecRange := &tbmcir.FilterSpecsByRangeRequest{}
	if err := c.Bind(vmSpecRange); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	fmt.Println("vmSpecRange.ConnectionName=", vmSpecRange.ConnectionName)
	fmt.Println(vmSpecRange)

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	resultVmSpecInfo, respStatus := service.FilterVmSpecInfoListByRange(defaultNameSpaceID, vmSpecRange)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		// 호출은 정상: http.StatusOK, 결과는 정상이 아님. (statusCode != 200,201)
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}
	// respBody := resp.Body
	// respStatusCode := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"VmSpecList": resultVmSpecInfo,
	})
}

/*
		connection에 대해 resource 목록 CSP와 Tumblebug을 비교
	    비교 가능 resource : vnet/securityGroup/sshKey/vm
*/
func GetInspectResourceList(c echo.Context) error {
	log.Println("GetInspectResourceList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramInspectResource := new(tbcommon.RestInspectResourcesRequest)
	if err := c.Bind(paramInspectResource); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	inspectResource, respStatus := service.GetInspectResourceList(paramInspectResource)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":         "success",
		"status":          respStatus,
		"inspectResource": inspectResource,
	})
}

/*
		connection에 대해 resource 목록 CSP와 Tumblebug을 비교
	    비교 가능 resource : vnet/securityGroup/sshKey/vm
*/
func GetInspectResourcesOverview(c echo.Context) error {
	log.Println("GetInspectResourceList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	inspectResourceAllResult, respStatus := service.GetInspectResourcesOverview()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                  "success",
		"status":                   respStatus,
		"inspectResourceAllResult": inspectResourceAllResult,
	})
}

// Register CSP Native Resources to CB-Tumblebug
func RegisterCspResourcesProc(c echo.Context) error {
	log.Println("RegisterCspResourcesProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	resourcesRequest := new(tbcommon.RestRegisterCspNativeResourcesRequest)
	if err := c.Bind(resourcesRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	// log.Println(vNetRegInfo)
	optionParam := c.QueryParam("option")
	registerResourceResult, respStatus := service.RegCspResources(resourcesRequest, optionParam)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {

		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus.StatusCode,
		"registerResource": registerResourceResult,
	})
}

func RegisterCspResourcesAllProc(c echo.Context) error {
	log.Println("RegisterCspResourcesProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	resourcesRequest := new(tbcommon.RestRegisterCspNativeResourcesRequestAll)
	if err := c.Bind(resourcesRequest); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	// log.Println(vNetRegInfo)
	optionParam := c.QueryParam("option")
	registerResourceResult, respStatus := service.RegCspResourcesAll(resourcesRequest, optionParam)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {

		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "success",
		"status":           respStatus.StatusCode,
		"registerResource": registerResourceResult,
	})
}

// resourcesGroup.GET("/datadisk/mngform", controller.DataDiskMngForm)
func DataDiskMngForm(c echo.Context) error {
	fmt.Println("DataDiskMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	//defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}
	store.Set("namespace", nsList)
	store.Save()
	log.Println(" nsList  ", nsList)

	//optionParam := c.QueryParam("option")
	//filterKeyParam := c.QueryParam("filterKey")
	//filterValParam := c.QueryParam("filterVal")
	//dataDiskList, respStatus := service.GetDataDiskList(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/DataDiskMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			//		"DataDiskList":  dataDiskList,
			//		"status":        respStatus.StatusCode,
		})
}

// e.GET("/setting/resources/datadisk/list", controller.DataDiskList)
func DataDiskList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")
	if optionParam == "id" {
		dataDiskInfoList, respStatus := service.GetDataDiskListByID(defaultNameSpaceID, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"defaultNameSpaceID": defaultNameSpaceID,
			"dataDiskInfoList":   dataDiskInfoList,
		})
	} else {
		dataDiskInfoList, respStatus := service.GetDataDiskListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"defaultNameSpaceID": defaultNameSpaceID,
			"dataDiskInfoList":   dataDiskInfoList,
		})
	}

}

// e.POST("/setting/resources/datadisk/reg", controller.DataDiskRegProc):
func DataDiskRegProc(c echo.Context) error {
	log.Println("DataDiskRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	dataDiskRegInfo := new(tbmcir.TbDataDiskReq)
	if err := c.Bind(dataDiskRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	resultDataDiskInfo, respStatus := service.RegDataDisk(defaultNameSpaceID, dataDiskRegInfo)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "success",
		"status":       respStatus.StatusCode,
		"DataDiskInfo": resultDataDiskInfo,
	})
}

// e.DELETE("/setting/resources/datadisk/del", controller.)
func DataDiskAllDelProc(c echo.Context) error {
	log.Println("ConnectionConfigDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	respMessage, respStatus := service.DelAllDataDisk(defaultNameSpaceID)
	fmt.Println("=============respMessage =============", respMessage)
	log.Println("respStatus : ", respStatus)
	log.Println("respMessage : ", respMessage)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respMessage.StatusCode,
	})
}

// e.GET("/setting/resources/datadisk/:dataDiskID", controller.DataDiskGet)
func DataDiskGet(c echo.Context) error {
	log.Println("DataDiskGet : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	dataDiskID := c.Param("dataDiskID")
	dataDiskInfo, respStatus := service.DataDiskGet(defaultNameSpaceID, dataDiskID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "success",
		"status":       respStatus,
		"dataDiskInfo": dataDiskInfo,
	})
}

// e.PUT("/setting/resources/datadisk/:dataDiskID", controller.DataDiskPutProc)
func DataDiskPutProc(c echo.Context) error {
	log.Println("DataDiskPutProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	dataDiskID := c.Param("dataDiskID")
	dataDiskUpsizeReq := new(tbmcir.TbDataDiskUpsizeReq)
	if err := c.Bind(dataDiskUpsizeReq); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	resultDataDiskInfo, respStatus := service.DataDiskPut(defaultNameSpaceID, dataDiskID, dataDiskUpsizeReq)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "success",
		"status":       respStatus.StatusCode,
		"DataDiskInfo": resultDataDiskInfo,
	})
}

// e.DELETE("/setting/resources/datadisk/del/:dataDiskID", controller.)
func DataDiskDelProc(c echo.Context) error {
	log.Println("DataDiskDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramDataDiskID := c.Param("dataDiskID")

	respMessage, respStatus := service.DelDataDisk(defaultNameSpaceID, paramDataDiskID)
	fmt.Println("=============respMessage =============", respMessage)
	log.Println("respStatus : ", respStatus)
	log.Println("respMessage : ", respMessage)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respMessage.StatusCode,
	})
}

// Create, Update, Delete를 한번에 하는 Controller
// 단, Attach, Detach 는 1개 vm에 대해서만 가능하게?
func DataDiskMngProc(c echo.Context) error {

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// 받아온 param maiing
	dataDiskReq := new(webtool.DataDiskMngReq)
	if err := c.Bind(dataDiskReq); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	mcisID := c.QueryParam("mcisID")
	vmID := c.QueryParam("vmID")

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	// create data disk list
	if dataDiskReq.CreateDataDiskList != nil {
		for _, createDataDisk := range dataDiskReq.CreateDataDiskList {
			// 생성 후 attach 할 vm이 있으면 attach 한다.
			// disk생성의 경우 delay가 충분히 있어야 한다.
			go service.AsyncRegDataDisk(defaultNameSpaceID, &createDataDisk, c)
		}

	}

	// detach data disk list
	if dataDiskReq.DetachDataDiskList != nil {
		if mcisID != "" && vmID != "" {
			for _, detachDataDisk := range dataDiskReq.DetachDataDiskList {

				// 2. vm에 attach
				optionParam := "detach"
				attachDetachDataDiskReq := new(tbmcir.TbAttachDetachDataDiskReq)
				attachDetachDataDiskReq.DataDiskId = detachDataDisk

				go service.AsyncAttachDetachDataDiskToVM(defaultNameSpaceID, mcisID, vmID, optionParam, attachDetachDataDiskReq, c)

			}
		}
	}

	// attach data disk list
	if dataDiskReq.AttachDataDiskList != nil {
		if mcisID != "" && vmID != "" {
			for _, attachDataDisk := range dataDiskReq.AttachDataDiskList {
				optionParam := "attach"
				attachDetachDataDiskReq := new(tbmcir.TbAttachDetachDataDiskReq)
				attachDetachDataDiskReq.DataDiskId = attachDataDisk
				go service.AsyncAttachDetachDataDiskToVM(defaultNameSpaceID, mcisID, vmID, optionParam, attachDetachDataDiskReq, c)
			}
		}
	}

	// delete data disk list
	if dataDiskReq.DetachDataDiskList != nil {
		for _, delDataDisk := range dataDiskReq.DetachDataDiskList {
			go service.AsyncDelDataDisk(defaultNameSpaceID, delDataDisk, c)
		}
	}
	// return result

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  200,
	})
}

// MyImage 관리화면 호출
func MyImageMngForm(c echo.Context) error {
	fmt.Println("MyImage ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	//defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)

	// 최신 namespacelist 가져오기
	//nsList, _ := service.GetNameSpaceList()
	nsList := []tbcommon.TbNsInfo{}
	if loginInfo.UserID == "admin" {
		nsList, _ = service.GetNameSpaceList()
	} else {
		nsList, _ = db.GetUserNamespaceList(loginInfo.UserID)
	}
	store.Set("namespace", nsList)
	store.Save()

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/MyImageMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
		})
}

// MyImage 목록조회
func MyImageList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	optionParam := c.QueryParam("option")
	filterKeyParam := c.QueryParam("filterKey")
	filterValParam := c.QueryParam("filterVal")
	if optionParam == "id" {
		myImageInfoList, respStatus := service.GetMyImageListByID(defaultNameSpaceID, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"defaultNameSpaceID": defaultNameSpaceID,
			"myImageInfoList":    myImageInfoList,
		})
	} else {
		myImageInfoList, respStatus := service.GetMyImageListByOption(defaultNameSpaceID, optionParam, filterKeyParam, filterValParam)
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
			return c.JSON(respStatus.StatusCode, map[string]interface{}{
				"error":  respStatus.Message,
				"status": respStatus.StatusCode,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":            "success",
			"status":             respStatus.StatusCode,
			"defaultNameSpaceID": defaultNameSpaceID,
			"myImageInfoList":    myImageInfoList,
		})
	}

}

// MyImage 등록 : csp에만 등록된 customImage를 cb-tb에 등록
func MyImageRegProc(c echo.Context) error {
	log.Println("MyImageRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	myImageRegInfo := new(tbmcir.TbCustomImageReq)
	if err := c.Bind(myImageRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	resultMyImageInfo, respStatus := service.RegCspCustomImageToMyImage(defaultNameSpaceID, myImageRegInfo)

	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "success",
		"status":      respStatus.StatusCode,
		"MyImageInfo": resultMyImageInfo,
	})
}

// MyImage 모두 제거 : ui에서는 빼는게 나을 지...
func MyImageAllDelProc(c echo.Context) error {
	log.Println("MyImageAllDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	respMessage, respStatus := service.DelAllMyImage(defaultNameSpaceID)
	fmt.Println("=============respMessage =============", respMessage)
	log.Println("respStatus : ", respStatus)
	log.Println("respMessage : ", respMessage)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respMessage.StatusCode,
	})
}

// MyImage상세 조회
func MyImageGet(c echo.Context) error {
	log.Println("MyImageGet : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	myImageID := c.Param("myImageID")
	myImageInfo, respStatus := service.MyImageGet(defaultNameSpaceID, myImageID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "success",
		"status":      respStatus,
		"myImageInfo": myImageInfo,
	})
}

// MyImage삭제
func MyImageDelProc(c echo.Context) error {
	log.Println("MyImageDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	myImageID := c.Param("myImageID")

	respMessage, respStatus := service.DelMyImage(defaultNameSpaceID, myImageID)
	fmt.Println("=============respMessage =============", respMessage)
	log.Println("respStatus : ", respStatus)
	log.Println("respMessage : ", respMessage)
	if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		return c.JSON(respStatus.StatusCode, map[string]interface{}{
			"error":  respStatus.Message,
			"status": respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": respMessage.Message,
		"status":  respMessage.StatusCode,
	})
}

// Provider, connection 에서 사용가능한 DiskType 조회
// 현재 : spider의 cloudos_meta.yaml 값 사용
func DataDiskLookupList(c echo.Context) error {
	log.Println("DataDiskProviderDisList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	//defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	provider := c.QueryParam("provider")
	connectionName := c.QueryParam("connectionName")

	if provider == "" {
		connectionInfo, respStatus := service.GetCloudConnectionConfigData(connectionName) // (spider.CloudConnectionConfigInfo
		if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
		}
		provider = connectionInfo.ProviderName
	}

	diskInfoList, err := service.DiskLookup(provider, connectionName)
	if err != nil {

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "success",
		"status":       "success",
		"DiskInfoList": diskInfoList,
	})
}
