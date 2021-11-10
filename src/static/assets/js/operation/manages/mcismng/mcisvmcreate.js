$(document).ready(function () {
    getVmList()
    getCommonNetworkList('vmcreate')
    getCommonVirtualMachineImageList('vmcreate')
    getCommonVirtualMachineSpecList('vmcreate')
    getCommonSecurityGroupList('vmcreate')
    getCommonSshKeyList('vmcreate')
    // e_vNetListTbody

    $('#alertResultArea').on('hidden.bs.modal', function () {// bootstrap 3 또는 4
        //$('#alertResultArea').on('hidden', function () {// bootstrap 2.3 이전
        let targetUrl = "/operation/manages/mcismng/mngform"
        changePage(targetUrl)
    })

    //OS_HW popup table scrollbar
    $('#OS_HW .btn_spec').on('click', function () {
        console.log("os_hw bpn_spec clicked ** ")
        $('#OS_HW_Spec .dtbox.scrollbar-inner').scrollbar();

        // connection 정보 set
        var esSelectedProvider = $("#es_regProvider option:selected").val();
        var esSelectedRegion = $("#es_regRegion option:selected").val();
        var esSelectedConnectionName = $("#es_regConnectionName option:selected").val();

        console.log("OS_HW_Spec_Assist click");
        if (esSelectedProvider) {
            $("#assist_select_provider").val(esSelectedProvider);
        }
        if (esSelectedRegion) {
            $("#assist_select_resion").val(esSelectedRegion);
        }
        if (esSelectedConnectionName) {
            $("#assist_select_connectionName").val(esSelectedConnectionName);
        }

        console.log("esSelectedProvider = " + esSelectedProvider + " : " + $("#assist_select_provider").val());
        console.log("esSelectedRegion = " + esSelectedRegion + " : " + $("#assist_select_resion").val());
        console.log("esSelectedConnectionName = " + esSelectedConnectionName + " : " + $("#assist_select_connectionName").val());
    });
    //Security popup table scrollbar
    $('#Security .btn_edit').on('click', function () {
        $("#security_edit").modal();
        $('#security_edit .dtbox.scrollbar-inner').scrollbar();
    });

    // $("input[name='vmInfoType']:radio").change(function () {
    //     //라디오 버튼 값을 가져온다.
    //     var formType = this.value;

    // });


    // server add 버튼 클릭 시
    // $('.servers_box .server_add').click(function(){	

    //     //<div class="servers_config import_servers_config" id="importServerConfig">
    //     //<div class="servers_config new_servers_config" id="expertServerConfig">
    // });

    //Servers Expert on/off
    //     var check = $(".switch .ch");
    //     var $Servers = $(".servers_config");
    //     var $NewServers = $(".new_servers_config");
    //     var $SimpleServers = $(".simple_servers_config");
    //     var simple_config_cnt = 0;
    //     var expert_config_cnt = 0;

    //     check.click(function(){
    //         $(".switch span.txt_c").toggle();
    //         $NewServers.removeClass("active");
    //     });

    //   //Expert add
    //     $('.servers_box .server_add').click(function(){
    //         $NewServers.toggleClass("active");
    //       if($Servers.hasClass("active")) {
    //         $Servers.toggleClass("active");
    //     } else {
    //         $Servers.toggleClass("active");
    //     }
    //     });
    //     // Simple add
    //   $(".servers_box .switch").change(function() {
    //     if ($(".switch .ch").is(":checked")) {	
    //             $('.servers_box .server_add').click(function(){	

    //                 $NewServers.addClass("active");
    //                 $SimpleServers.removeClass("active");		
    //             });
    //     } else {
    //             $('.servers_box .server_add').click(function(){

    //                 $NewServers.removeClass("active");
    //                 $SimpleServers.addClass("active");


    //             });		
    //     }
    //   });

});


var totalDeployServerCount = 0;
function btn_deploy() {
    var mcis_name = $("#mcis_name").val();
    var mcis_id = $("#mcis_id").val();
    if (!mcis_id) {
        commonAlert("Please Select MCIS !!!!!")
        return;
    }
    totalDeployServerCount = 0;// deploy vm 개수 초기화

    console.log(Simple_Server_Config_Arr);
    if (Simple_Server_Config_Arr) {// mcissimpleconfigure.js 에 const로 정의 됨.
        var vm_len = Simple_Server_Config_Arr.length;
        if (vm_len > 0) {
            totalDeployServerCount += vm_len
            console.log("Simple_Server_Config_Arr length: ", vm_len);
            // var new_obj = {}
            // new_obj['vm'] = Simple_Server_Config_Arr;
            // console.log("new obj is : ",new_obj);
            // var url = "/operation/manages/mcis/:mcisID/vm/reg/proc"
            var url = "/operation/manages/mcismng/" + mcis_id + "/vm/reg/proc"

            // 한개씩 for문으로 추가
            for (var i in Simple_Server_Config_Arr) {
                new_obj = Simple_Server_Config_Arr[i];
                console.log("new obj is : ", new_obj);
                try {
                    resultVmCreateMap.set(new_obj.name, "")
                    axios.post(url, new_obj, {
                        headers: {
                        },
                    }).then(result => {
                        console.log("MCIR VM Register data : ", result);
                        console.log("Result Status : ", result.status);

                        var statusCode = result.data.status;
                        var message = result.data.message;
                        console.log("Result Status : ", statusCode);
                        console.log("Result message : ", message);

                        if (result.status == 201 || result.status == 200) {
                            vmCreateCallback(new_obj.name, "Success")
                            //     commonAlert("Register Success")
                            //     // location.href = "/Manage/MCIS/list";
                            //     // $('#loadingContainer').show();
                            //     // location.href = "/operation/manages/mcis/mngform/"
                            //     var targetUrl = "/operation/manages/mcis/mngform"
                            //     changePage(targetUrl)
                        } else {
                            vmCreateCallback(new_obj.name, message)
                            //     commonAlert("Register Fail")
                            //     //location.reload(true);
                        }
                    }).catch((error) => {
                        // console.warn(error);
                        console.log(error.response)
                        var errorMessage = error.response.data.error;
                        // commonErrorAlert(statusCode, errorMessage) 
                        vmCreateCallback(new_obj.name, errorMessage)
                    })
                } finally {

                }

                // post로 호출을 했으면 해당 VM의 정보는 비활성시킨 후(클릭 Evnet 안먹게)
                // 상태값을 모니터링 하여 결과 return 까지 대기.
                // return 받으면 해당 VM
            }
        }
    }

    ///////// export
    console.log(Expert_Server_Config_Arr);
    if (Expert_Server_Config_Arr) {
        var vm_len = Expert_Server_Config_Arr.length;
        console.log("Expert_Server_Config_Arr length: ", vm_len);
        if (vm_len > 0) {
            totalDeployServerCount += vm_len
            // var new_obj = {}
            // new_obj['vm'] = Simple_Server_Config_Arr;
            // console.log("new obj is : ",new_obj);
            // var url = "/operation/manages/mcis/:mcisID/vm/reg/proc"
            var url = "/operation/manages/mcismng/" + mcis_id + "/vm/reg/proc"

            // 한개씩 for문으로 추가
            for (var i in Expert_Server_Config_Arr) {
                new_obj = Expert_Server_Config_Arr[i];
                console.log("new obj is : ", new_obj);
                try {
                    resultVmCreateMap.set("Expert" + i, "")
                    axios.post(url, new_obj, {
                        headers: {
                        },
                    }).then(result => {
                        console.log("MCIR VM Register data : ", result);
                        console.log("Result Status : ", result.status);

                        var statusCode = result.data.status;
                        var message = result.data.message;
                        console.log("Result Status : ", statusCode);
                        console.log("Result message : ", message);

                        if (result.status == 201 || result.status == 200) {
                            vmCreateCallback("Expert" + i, "Success")
                        } else {
                            vmCreateCallback("Expert" + i, message)
                        }
                    }).catch((error) => {
                        // console.warn(error);
                        console.log(error.response)
                        var errorMessage = error.response.data.error;
                        // commonErrorAlert(statusCode, errorMessage) 
                        vmCreateCallback("Expert" + i, errorMessage)
                    })
                } finally {

                }

                // post로 호출을 했으면 해당 VM의 정보는 비활성시킨 후(클릭 Evnet 안먹게)
                // 상태값을 모니터링 하여 결과 return 까지 대기.
                // return 받으면 해당 VM
            }
        }
    }
    ///////// import
    if (Import_Server_Config_Arr) {// mcissimpleconfigure.js 에 const로 정의 됨.
        // TODO : 어차피 simple/expert와 로직이 다른데... 
        // json 그대로 넘기도록
        var vm_len = Import_Server_Config_Arr.length;
        if (vm_len > 0) {
            console.log("Import_Server_Config_Arr length: ", vm_len);
            totalDeployServerCount += vm_len
            // var new_obj = {}
            // new_obj['vm'] = Simple_Server_Config_Arr;
            // console.log("new obj is : ",new_obj);
            // var url = "/operation/manages/mcis/:mcisID/vm/reg/proc"
            var url = "/operation/manages/mcismng/" + mcis_id + "/vm/reg/proc"

            // 한개씩 for문으로 추가
            for (var i in Import_Server_Config_Arr) {
                new_obj = Import_Server_Config_Arr[i];
                console.log("new obj is : ", new_obj);
                try {
                    resultVmCreateMap.set("Import" + i, "")
                    axios.post(url, new_obj, {
                        headers: {
                        },
                    }).then(result => {
                        console.log("MCIR VM Register data : ", result);
                        console.log("Result Status : ", result.status);

                        var statusCode = result.data.status;
                        var message = result.data.message;
                        console.log("Result Status : ", statusCode);
                        console.log("Result message : ", message);

                        if (result.status == 201 || result.status == 200) {
                            vmCreateCallback("Import" + i, "Success")
                            //     commonAlert("Register Success")
                            //     // location.href = "/Manage/MCIS/list";
                            //     // $('#loadingContainer').show();
                            //     // location.href = "/operation/manages/mcis/mngform/"
                            //     var targetUrl = "/operation/manages/mcis/mngform"
                            //     changePage(targetUrl)
                        } else {
                            vmCreateCallback("Import" + i, message)
                            //     commonAlert("Register Fail")
                            //     //location.reload(true);
                        }
                    }).catch((error) => {
                        // console.warn(error);
                        console.log(error.response)
                        var errorMessage = error.response.data.error;
                        // commonErrorAlert(statusCode, errorMessage) 
                        vmCreateCallback("Import" + i, errorMessage)
                    })
                } finally {

                }

                // post로 호출을 했으면 해당 VM의 정보는 비활성시킨 후(클릭 Evnet 안먹게)
                // 상태값을 모니터링 하여 결과 return 까지 대기.
                // return 받으면 해당 VM
            }
        }
    }
}

// Import / Export Modal 표시
function btn_ImportExport() {
    // export할 VM을 선택한 후 export 버튼 누르라고...
    $("#VmImportExport").modal();
    $('#VmImportExport .dtbox.scrollbar-inner').scrollbar();
}

// vm 생성 결과 표시
// 여러개의 vm이 생성될 수 있으므로 각각 결과를 표시
var resultVmCreateMap = new Map();
function vmCreateCallback(resultVmKey, resultStatus) {
    resultVmCreateMap.set(resultVmKey, resultStatus)
    var resultText = "";
    var createdServer = 0;
    for (let key of resultVmCreateMap.keys()) {
        console.log("vmCreateresult " + key + " : " + resultVmCreateMap.get(resultVmKey));
        resultText += key + " = " + resultVmCreateMap.get(resultVmKey) + ","
        //totalDeployServerCount--
        createdServer++;
    }

    $("#serverRegistResult").text(resultText);

    if (resultStatus != "Success") {
        // add된 항목 제거 해야 함.

        // array는 초기화
        Simple_Server_Config_Arr.length = 0;
        simple_data_cnt = 0
        // TODO : expert 추가하면 주석 제거할 것
        Expert_Server_Config_Arr.length = 0;
        expert_data_cnt = 0
        Import_Server_Config_Arr.length = 0;
        import_data_cnt = 0
    }

    if (createdServer === totalDeployServerCount) { //모두 성공
        //getVmList();
        commonResultAlert($("#serverRegistResult").text());

    } else if (createdServer < totalDeployServerCount) { //일부 성공
        commonResultAlert($("#serverRegistResult").text());

    } else if (createdServer = 0) { //모두 실패
        commonResultAlert($("#serverRegistResult").text());

    }
}

// 현재 mcis의 vm 목록 조회 : deploy후 상태볼 때 사용
function getVmList() {

    console.log("getVmList()")
    var mcis_id = $("#mcis_id").val();


    // /operation/manages/mcis/:mcisID
    var url = "/operation/manages/mcismng/" + mcis_id
    axios.get(url, {})
        .then(result => {
            console.log("MCIR VM Register data : ", result);
            console.log("Result Status : ", result.status);

            var statusCode = result.data.status;
            var message = result.data.message;
            //
            console.log("Result Status : ", statusCode);
            console.log("Result message : ", message);


            if (result.status == 201 || result.status == 200) {
                var mcis = result.data.McisInfo
                console.log(mcis)


                var vms = mcis.vm
                if (vms) {
                    vm_len = vms.length

                    $("#mcis_server_list *").remove();
                    var appendLi = "";

                    for (var o in vms) {
                        var vm_status = vms[o].status
                        var vm_name = vms[o].name

                        console.log(o + "번째 " + vm_name + " : " + vm_status)
                        // mcis_server_list 밑의 li들을 1개빼고 삭제. 
                        // 가져온 vm list 를 add? (1개는 더하기 버튼이므로)                    


                        appendLi = appendLi + '<li>';
                        appendLi = appendLi + '<div class="server server_on bgbox_g">';
                        appendLi = appendLi + '<div class="icon"></div>';
                        appendLi = appendLi + '<div class="txt">' + vm_name + '</div>';
                        appendLi = appendLi + '</li>';

                        appendLi = appendLi + '</li>';

                    }
                    appendLi = appendLi + '<li>';
                    appendLi = appendLi + '<div class="server server_add" onClick="displayNewServerForm()">';
                    appendLi = appendLi + '</div>';
                    appendLi = appendLi + '</li>';

                    $("#mcis_server_list").append(appendLi);

                    // commonAlert("VM 목록 조회 완료")
                    //$("#serverRegistResult").text("VM 목록 조회 완료");
                }
            }
        }).catch((error) => {
            // console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.data.error;
        })
}

// 화면 Load시 가져오나 굳이?
function getNetworkListCallbackSuccess(caller, data) {
    console.log(data);
    if (data == null || data == undefined || data == "null") {

    } else {// 아직 data가 1건도 없을 수 있음
        var html = ""
        if (data.length > 0) {
            data.forEach(function (vNetItem, vNetIndex) {
                // TODO : 생성 function으로 뺄 것. vnet에 subnet이 2개 이상 있을 수 있는데 그중 1개의 subnet을 선택해야 함.
                var subnetHtml = ""
                var subnetData = vNetItem.subnetInfoList
                var subnetIds = ""
                console.log(subnetData)
                subnetData.forEach(function (subnetItem, subnetIndex) {
                    // subnetHtml +='<input type="hidden" name="vNet_subnet_' + vNetItem.id + '" id="vNet_subnet_' + vNetItem.id + '_' + subnetIndex + '" value="' + subnetItem.iid.nameId + '"/>'
                    //             + subnetIndex + ' || ' + subnetItem.iid.nameId + ' <p>'
                    // console.log(subnetItem)
                    // console.log(subnetItem.iid)
                    subnetHtml += subnetIndex + ' || ' + subnetItem.id + '<p>'
                    if (subnetIndex > 0) { subnetIds += "," }
                    subnetIds += subnetItem.name

                })
                subnetIds += ""
                subnetHtml += '<input type="hidden" name="vNet_subnet_' + vNetItem.id + '" id="vNet_subnet_' + vNetItem.id + '_' + vNetIndex + '" value="' + subnetIds + '"/>'

                console.log("subnetIds = " + subnetIds)

                console.log(subnetHtml)
                html += '<tr onclick="setVnetValueToFormObj(\'es_vNetList\', \'tab_vNet\', \'vNetItem.ID\',\'vNet\',' + vNetIndex + ', \'e_vNetId\');">'

                    + '        <input type="hidden" id="vNet_id_' + vNetIndex + '" value="' + vNetItem.id + '"/>'
                    + '        <input type="hidden" name="vNet_connectionName" id="vNet_connectionName_' + vNetIndex + '" value="' + vNetItem.connectionName + '"/>'
                    + '        <input type="hidden" name="vNet_name" id="vNet_name_' + vNetIndex + '" value="' + vNetItem.name + '"/>'
                    + '        <input type="hidden" name="vNet_description" id="vNet_description_' + vNetIndex + '" value="' + vNetItem.description + '"/>'
                    + '        <input type="hidden" name="vNet_cidrBlock" id="vNet_cidrBlock_' + vNetIndex + '" value="' + vNetItem.cidrBlock + '"/>'
                    + '        <input type="hidden" name="vNet_cspVnetName" id="vNet_cspVnetName_' + vNetIndex + '" value="' + vNetItem.cspVNetName + '"/>'

                    + '        <input type="hidden" name="vNet_subnetInfos" id="vNet_subnetInfos_' + vNetIndex + '" value="' + subnetIds + '"/>'

                    //    사용하지 않는데 굳이 리스트를 할당할 필요가 있을까?
                    //+'        <input type="hidden" name="vNet_keyValueInfos" id="vNet_keyValueInfos_' + vNetIndex + '" value="' + vNetItem.keyValueInfos + '"/>'

                    + '        <input type="hidden" id="vNet_info_' + vNetIndex + '" value="' + vNetItem.id + '|' + vNetItem.name + ' |' + vNetItem.cspVNetName + '|' + vNetItem.cidrBlock + '|' + subnetIds + '"/>'

                    + '    <td class="overlay hidden" data-th="Name">' + vNetItem.name + '</td>'
                    + '    <td class="btn_mtd ovm td_left" data-th="CidrBlock">'
                    + '        ' + vNetItem.cidrBlock
                    + '    </td>'
                    + '    <td class="btn_mtd ovm td_left" data-th="SubnetInfo">' + subnetHtml
                    // +'        { {range $subnetIndex, $subnetItem := .SubnetInfos + ''
                    // +'        <input type="hidden" name="vNet_subnet_' + vNetItem.ID + '" id="vNet_subnet_' + vNetItem.ID + '_' + subnetIndex + '" value="' + subnetItem.IID.NameId + '"/>'
                    // +'        ' + subnetIndex + ' || ' + subnetItem.IID.NameId + ' <p>'
                    // +'        { { end  + ''
                    + '    </td>'
                    + '    <td class="overlay hidden" data-th="Description">' + vNetItem.description + '</td>'
                    + '</tr>'
            })
            $("#e_vNetListTbody").empty()
            $("#e_vNetListTbody").append(html)
        }
    }

}
function getNetworkListCallbackFail(caller, error) {
    // no data
    var html = ""
    html += '<tr>'
        + '<td class="overlay hidden" data-th="" colspan="4">No Data</td>'
        + '</tr>';
    $("#e_vNetListTbody").empty()
    $("#e_vNetListTbody").append(html)
}

function getSpecListCallbackSuccess(caller, data) {
    console.log(data);
    if (data == null || data == undefined || data == "null") {

    } else {// 아직 data가 1건도 없을 수 있음
        var html = ""
        if (data.length > 0) {
            data.forEach(function (vSpecItem, vSpecIndex) {

                html += '<tr onclick="setValueToFormObj(\'tab_vmSpec\', \'vmSpec\',' + vSpecIndex + ', \'e_specId\');">'
                    + '     <input type="hidden" id="vmSpec_id_' + vSpecIndex + '" value="' + vSpecItem.id + '"/>'
                    + '     <input type="hidden" name="vmSpec_connectionName" id="vmSpec_connectionName_' + vSpecIndex + '" value="' + vSpecItem.connectionName + '"/>'
                    + '     <input type="hidden" name="vmSpec_info" id="vmSpec_info_' + vSpecIndex + '" value="' + vSpecItem.id + '|' + vSpecItem.name + '|' + vSpecItem.connectionName + '|' + vSpecItem.cspImageId + '|' + vSpecItem.cspImageName + '|' + vSpecItem.guestOS + '|' + vSpecItem.description + '"/>'
                    + '<td class="overlay hidden" data-th="Name">' + vSpecItem.name + '</td>'
                    + '<td class="btn_mtd ovm td_left" data-th="ConnectionName">'
                    + vSpecItem.connectionName
                    + '</td>'
                    + '<td class="overlay hidden" data-th="CspSpecName">' + vSpecItem.cspSpecName + '</td>'

                    + '<td class="overlay hidden" data-th="Description">' + vSpecItem.description + '</td>'
                    + '</tr>'

            })
            $("#e_specListTbody").empty()
            $("#e_specListTbody").append(html)
        }
    }
}
function getSpecListCallbackFail(caller, error) {
    // no data
    var html = ""
    html += '<tr>'
        + '<td class="overlay hidden" data-th="" colspan="4">No Data</td>'
        + '</tr>';
    $("#e_specListTbody").empty()
    $("#e_specListTbody").append(html)
}

function getImageListCallbackSuccess(caller, data) {
    console.log(data);
    if (data == null || data == undefined || data == "null") {

    } else {// 아직 data가 1건도 없을 수 있음
        var html = ""
        if (data.length > 0) {
            data.forEach(function (vImageItem, vImageIndex) {

                html += '<tr onclick="setValueToFormObj(\'es_imageList\', \'tab_vmImage\', \'vmImage\',' + vImageIndex + ', \'e_imageId\');">'
                    + '     <input type="hidden" id="vmImage_id_' + vImageIndex + '" value="' + vImageItem.id + '"/>'
                    + '     <input type="hidden" name="vmImage_connectionName" id="vmImage_connectionName_' + vImageIndex + '" value="' + vImageItem.connectionName + '"/>'
                    + '     <input type="hidden" name="vmImage_info" id="vmImage_info_' + vImageIndex + '" value="' + vImageItem.id + '|' + vImageItem.name + '|' + vImageItem.connectionName + '|' + vImageItem.cspImageId + '|' + vImageItem.cspImageName + '|' + vImageItem.guestOS + '|' + vImageItem.description + '"/>'

                    + '<td class="overlay hidden" data-th="Name">' + vImageItem.name + '</td>'
                    + '<td class="btn_mtd ovm td_left" data-th="ConnectionName">'
                    + vImageItem.connectionName
                    + '</td>'
                    + '<td class="overlay hidden" data-th="CspImageId">' + vImageItem.cspImageId + '</td>'
                    + '<td class="overlay hidden" data-th="CspImageName">' + vImageItem.cspImageName + '</td>'
                    + '<td class="overlay hidden" data-th="GuestOS">' + vImageItem.guestOS + '</td>'
                    + '<td class="overlay hidden" data-th="Description">' + vImageItem.description + '</td>'
                    + '</tr>'
            })
            $("#es_imageListTbody").empty()
            $("#es_imageListTbody").append(html)
        }
    }
}
function getImageListCallbackFail(error) {
    // no data
    var html = ""
    html += '<tr>'
        + '<td class="overlay hidden" data-th="" colspan="4">No Data</td>'
        + '</tr>';
    $("#es_imageListTbody").empty()
    $("#es_imageListTbody").append(html)
}

// 등록 된 vm search 결과
function getCommonSearchVmImageListCallbackSuccess(caller, vmImageList){
    console.log("11111");
    console.log(vmImageList);
    var html = ""
    if (vmImageList.length > 0) {
        vmImageList.forEach(function (vImageItem, vImageIndex) {
            //connectionName
            //cspSpecName
            html += '<tr onclick="setAssistValue(' + vImageIndex + ');">'
                + '     <input type="hidden" id="vmImageAssist_name_' + vImageIndex + '" value="' + vImageItem.name + '"/>'
                + '     <input type="hidden" id="vmImageAssist_connectionName_' + vImageIndex + '" value="' + vImageItem.connectionName + '"/>'
                + '     <input type="hidden" id="vmImageAssist_cspImageId_' + vImageIndex + '" value="' + vImageItem.cspImageId + '"/>'
                + '     <input type="hidden" id="vmImageAssist_cspImageName_' + vImageIndex + '" value="' + vImageItem.cspImageName + '"/>'
                + '     <input type="hidden" id="vmImageAssist_guestOS_' + vImageIndex + '" value="' + vImageItem.guestOS + '"/>'
                + '     <input type="hidden" id="vmImageAssist_description_' + vImageIndex + '" value="' + vImageItem.description + '"/>'
                + '<td class="overlay hidden" data-th="Name">' + vImageItem.name + '</td>'
                + '<td class="overlay hidden" data-th="CspImageId">' + vImageItem.cspImageId + '</td>'
                + '<td class="overlay hidden" data-th="CspImageName">' + vImageItem.cspImageName + '</td>'
                + '</td>'
                + '<td class="overlay hidden" data-th="GuestOS">' + vImageItem.guestOS + '</td>'
                + '<td class="overlay hidden" data-th="Description">' + vImageItem.description + '</td>'
                + '</tr>'

        });
        $("#assistVmImageList").empty()
        $("#assistVmImageList").append(html)
    }
}

function getSecurityGroupListCallbackSuccess(caller, data){
    // expert에서 사용할 securityGroup
    if (data == null || data == undefined || data == "null") {

    } else {// 아직 data가 1건도 없을 수 있음
        var html = ""
        if (data.length > 0) {
            data.forEach(function (vSecurityGroupItem, vSecurityGroupIndex) {

                html += '<tr>'

                    + '<td class="overlay hidden column-50px" data-th="">'
                    + '     <input type="checkbox" name="securityGroup_chk" id="securityGroup_Raw_' + vSecurityGroupIndex + '" title="" />'
                    + '     <input type="hidden" id="securityGroup_id_' + vSecurityGroupIndex + '" value="' + vSecurityGroupItem.id + '"/>'
                    + '     <input type="hidden" id="securityGroup_name_' + vSecurityGroupIndex + '" value="' + vSecurityGroupItem.name + '"/>'
                    + '     <input type="hidden" name="securityGroup_connectionName" id="securityGroup_connectionName_' + vSecurityGroupIndex +'" value="' + vSecurityGroupItem.connectionName + '"/>'
                    + '     <input type="hidden" name="securityGroup_info" id="securityGroup_info_' + vSecurityGroupIndex + '" value="'+ vSecurityGroupItem.name +'|' + vSecurityGroupItem.connectionName + '|' + vSecurityGroupItem.description + '"/>'
                    + '     <label for="td_ch1"></label> <span class="ov off"></span>'
                    + '</td>'
                    + '<td class="btn_mtd ovm td_left" data-th="Name">'
                    + vSecurityGroupItem.name
                    + '</td>'
                    + '<td class="btn_mtd ovm td_left" data-th="ConnectionName">'
                    + vSecurityGroupItem.connectionName
                    + '</td>'
                    + '<td class="overlay hidden" data-th="Description">' + vSecurityGroupItem.description + '</td>'

                    + '</tr>'
            })
            $("#e_securityGroupListTbody").empty()
            $("#e_securityGroupListTbody").append(html)

        }
    }
}

function getSecurityGroupListCallbackFail(error){

}

function getSshKeyListCallbackSuccess(caller, data){
    // expert에서 사용할 sshkey
    if (data == null || data == undefined || data == "null") {

    } else {// 아직 data가 1건도 없을 수 있음
        var html = ""
        if (data.length > 0) {
            data.forEach(function (vSshKeyItem, vSshKeyIndex) {

                html += '<tr onclick="setValueToFormObj(\'es_sshKeyList\', \'tab_sshKey\', \'sshKey\', ' + vSshKeyIndex + ', \'e_sshKeyId\');">'

                    + '<td class="overlay hidden" data-th="Name">' + vSshKeyItem.name + '</td>'

                    + '     <input type="hidden" name="sshKey_id" id="sshKey_id_' + vSshKeyIndex + '" value="' + vSshKeyItem.id + '"/>'
                    + '     <input type="hidden" name="sshKey_connectionName" id="sshKey_connectionName_' + vSshKeyIndex + '" value="' + vSshKeyItem.connectionName + '"/>'
                    + '     <input type="hidden" name="sshKey_info" id="sshKey_info_' + vSshKeyIndex + '" value="' + vSshKeyItem.name + '|' + vSshKeyItem.connectionName + '|' + vSshKeyItem.description + '"/>'
                    + '</td>'
                    + '<td class="btn_mtd ovm td_left" data-th="ConnectionName">'
                    + vSshKeyItem.connectionName
                    + '</td>'
                    + '<td class="overlay hidden" data-th="Description">' + vSshKeyItem.description + '</td>'

                    + '</tr>'
            })
            $("#e_sshKeyListTbody").empty()
            $("#e_sshKeyListTbody").append(html)

        }
    }
}

function getSshKeyListCallbackFail(caller, error){

}

// EnterKey입력 시 해당 값, keyword 들이 있는 object id, 구분자(caller)
function searchAssistImageByEnter(event, caller){
    if( event.keyCode === 13) {
        // searchKeyword(keyword, caller);
        searchVmImageByKeyword(caller);
        // searchKeyword($(this).val(), caller)
    }
}

//
function searchVmImageByKeyword(caller){
    var keyword = "";
    var keywordObjId = "";
    if( caller == "searchVmImageAssistAtReg"){
        keyword = $("#keywordAssistImage").val();
        keywordObjId = "searchAssistImageKeywords";
    }

    if( !keyword ){
        commonAlert("At least a keyword required");
        return;
    }
    var addKeyword = '<div class="keyword" name="keyword_' + caller + '">' + keyword.trim() + '<button class="btn_del_image" onclick="delSearchKeyword(event, \'' + caller + '\')"></button></div>';

    $("#" + keywordObjId).append(addKeyword);
    var keywords = new Array();// 기존에 있는 keyword에 받은 keyword 추가하여 filter적용
    $("[name='keyword_" + caller + "']").each(function( idx, ele){
        keywords.push($(this).text());
    });

    getCommonSearchVmImageList(keywords, caller);
}
// Assist Spec filter Search버튼 클릭시
function searchSpecsByRange(caller){
    // var specFilter = new Object();

    var assistSpecConnectionNameVal = $("#assistSpecConnectionName option:selected").val();
    if( caller == 'searchVmSpecAssistAtReg'){

    }
    // if (assistSpecConnectionNameVal) {
    //     specFilter.connectionName = assistSpecConnectionNameVal
    // }

    // storage
    var storageMin = $("#assist_num_storage_min").val();
    var storageMax = $("#assist_num_storage_max").val();

    var storageObj = new Object();
    storageObj.min = Number(storageMin)
    storageObj.max = Number(storageMax)

    // Core
    var coreMin = $("#assist_num_core_min").val();
    var coreMax = $("#assist_num_core_max").val();
    var coreObj = new Object();
    coreObj.min = Number(coreMin)
    coreObj.max = Number(coreMax)

    // specFilter.numCore = { "min": coreMin, "max": coreMax };

    // vCPU
    var vCpuMin = $("#assist_num_vCPU_min").val();
    var vCpuMax = $("#assist_num_vCPU_max").val();
    var vCpuObj = new Object();
    vCpuObj.min = Number(vCpuMin)
    vCpuObj.max = Number(vCpuMax)
    // specFilter.numvCPU = { "min": vCpuMin, "max": vCpuMax };

    // memory

    var specFilter = {
        connectionName: assistSpecConnectionNameVal,
        maxTotalStorageTiB: storageObj,
        numCore: coreObj,
        numvCPU: vCpuObj,
    }
    getCommonFilterVmSpecListByRange(specFilter, caller)

    // ID             string `json:"id"`
    // Name           string `json:"name"`
    // Description    string `json:"description"`
    // ConnectionName string `json:"connectionName"`
    // CspSpecName    string `json:"cspSpecName"`
    // OsType         string `json:"osType"`
    //
    // CostPerHour Range `json:"costPerHour"`
    // EbsBwMbps   Range `json:"ebsBwMbps"`
    //
    // EvaluationScore01 Range  `json:"evaluationScore01"`
    // EvaluationScore02 Range  `json:"evaluationScore02"`
    // EvaluationScore03 Range  `json:"evaluationScore03"`
    // EvaluationScore04 Range  `json:"evaluationScore04"`
    // EvaluationScore05 Range  `json:"evaluationScore05"`
    // EvaluationScore06 Range  `json:"evaluationScore06"`
    // EvaluationScore07 Range  `json:"evaluationScore07"`
    // EvaluationScore08 Range  `json:"evaluationScore08"`
    // EvaluationScore09 Range  `json:"evaluationScore09"`
    // EvaluationScore10 Range  `json:"evaluationScore10"`
    // EvaluationStatus  string `json:"evaluationStatus"`
    //
    // GpuModel string `json:"gpuModel"`
    // GpuP2p   string `json:"gpuP2p"`
    //
    // MaxNumStorage      Range `json:"maxNumStorage"`
    // MaxTotalStorageTiB Range `json:"maxTotalStorageTiB"`
    // MemGiB             Range `json:"memGiB"`
    //
    // NetBwGbps  Range `json:"netBwGbps"`
    // NumCore    Range `json:"numCore"`
    // NumGpu     Range `json:"numGpu"`
    // NumStorage Range `json:"numStorage"`
    // NumVCPU    Range `json:"numvCPU"`
    // StorageGiB Range `json:"storageGiB"`
}

// 등록된 spec조회 성공 시 table에 뿌려주고, 클릭시 spec 내용 set.
function getCommonFilterVmSpecListCallbackSuccess(caller, vmSpecList){
// function getCommonFilterVmImageListCallbackSuccess(caller, vmSpecList){
    console.log(vmSpecList);
    if (vmSpecList == null || vmSpecList == undefined || vmSpecList == "null") {

    } else {// 아직 data가 1건도 없을 수 있음
        var html = ""
        if (vmSpecList.length > 0) {
            vmSpecList.forEach(function (vSpecItem, vSpecIndex) {
                //connectionName
                //cspSpecName
                html += '<tr onclick="setAssistValue(' + vSpecIndex + ');">'
                    + '     <input type="hidden" id="vmSpecAssist_id_' + vSpecIndex + '" value="' + vSpecItem.id + '"/>'
                    + '     <input type="hidden" id="vmSpecAssist_name_' + vSpecIndex + '" value="' + vSpecItem.name + '"/>'
                    + '     <input type="hidden" id="vmSpecAssist_cspSpecName_' + vSpecIndex + '" value="' + vSpecItem.cspSpecName + '"/>'
                    + '     <input type="hidden" id="vmSpecAssist_memGiB_' + vSpecIndex + '" value="' + vSpecItem.memGiB + '"/>'
                    + '     <input type="hidden" id="vmSpecAssist_numvCPU_' + vSpecIndex + '" value="' + vSpecItem.numvCPU + '"/>'
                    + '     <input type="hidden" id="vmSpecAssist_numGpu_' + vSpecIndex + '" value="' + vSpecItem.numGpu + '"/>'
                    + '     <input type="hidden" id="vmSpec_connectionName_' + vSpecIndex + '" value="' + vSpecItem.connectionName + '"/>'
                    + '<td class="overlay hidden" data-th="Name">' + vSpecItem.name + '</td>'
                    + '<td class="overlay hidden" data-th="CspSpecName">' + vSpecItem.cspSpecName + '</td>'
                    + '<td class="btn_mtd ovm td_left" data-th="Memory">'
                    + vSpecItem.memGiB
                    + '</td>'
                    + '<td class="overlay hidden" data-th="VCPU">' + vSpecItem.numvCPU + '</td>'

                    + '<td class="overlay hidden" data-th="GPU">' + vSpecItem.numGpu + '</td>'
                    + '</tr>'
            })
            $("#assistSpecList").empty()
            $("#assistSpecList").append(html)
        }

            // "associatedObjectList": null,
            // "connectionName": "aws-conn-osaka",
            // "costPerHour": 0,
            // "cspSpecName": "t3.small",
            // "description": "",
            // "ebsBwMbps": 0,
            // "evaluationScore01": 0,
            // "evaluationScore02": 0,
            // "evaluationScore03": 0,
            // "evaluationScore04": 0,
            // "evaluationScore05": 0,
            // "evaluationScore06": 0,
            // "evaluationScore07": 0,
            // "evaluationScore08": 0,
            // "evaluationScore09": 0,
            // "evaluationScore10": 0,
            // "evaluationStatus": "",
            // "gpuMemGiB": 0,
            // "gpuModel": "",
            // "gpuP2p": "",
            // "id": "osaka-t3small",
            // "isAutoGenerated": false,
            // "maxNumStorage": 0,
            // "maxTotalStorageTiB": 0,
            // "memGiB": 2,
            // "name": "osaka-t3small",
            // "namespace": "osaka-ns",
            // "netBwGbps": 0,
            // "numCore": 0,
            // "numGpu": 0,
            // "numStorage": 0,
            // "numvCPU": 2,
            // "orderInFilteredResult": 0,
            // "osType": "",
            // "storageGiB": 0


    }
}

function clearAssistSpecList(targetTableList){
    $("#" + targetTableList).empty()
}