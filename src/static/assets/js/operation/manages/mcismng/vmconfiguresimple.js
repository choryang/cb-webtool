$(document).ready(function () {
	$('.btn_recommend').on('click', function () {
		showRecommendAssistPopup();
	});

})

function showRecommendAssistPopup() {
	console.log("showRecommendAssistPopup")
	$("#recommendVmAssist").modal();

}

function showConnectionAssistPopup() {
	console.log("showConnectionAssistPopup")
	$("#connectionAssist").modal();

}

// getConnectionListForSelectbox 로 변경
// function changeProvider(provider, target){
// }

// Connection 정보가 바뀌면 등록에 필요한 목록들을 다시 가져온다.(config는 ID가아닌 configName을 사용한다.)
function changeConnectionInfo(configName) {
	console.log("config name : ", configName)
	if (configName == "") {
		// 0번째면 selectbox들을 초기화한다.(vmInfo, sshKey, image 등)
	}
	getVmiInfo(configName);
	//getVmMyiInfo(configName);
	getSecurityInfo(configName);
	getSSHKeyInfo(configName);
	getVnetInfo(configName);
	getSpecInfo(configName);
}

function getVmiInfo(configName) {

	//var configName = $("#ss_regConnectionName option:selected").val();

	console.log("2 : ", configName);
	// getCommonVirtualMachineImageList("mcissimpleconfigure", "name"); setCommonVirtualMachineImageList()
	// var url = "/setting/resources" + "/machineimage/lookupimage";//TODO : 조회 오류남... why? connectionName으로 lookup
	var url = "/setting/resources" + "/machineimage/list"
	// var url = "http://54.248.3.145:1323/tumblebug/lookupImage";				 
	//  var url = CommonURL+"/ns/"+NAMESPACE+"/resources/image";
	var html = "";
	//  var apiInfo = 'Basic ZGVmYXVsdDpkZWZhdWx0'
	axios.get(url, {
		// headers:{
		// 	'Authorization': apiInfo
		// },
		params: {
			connectionName: configName
		}
	}).then(result => {
		console.log("Image Info : ", result.data)
		data = result.data.VirtualMachineImageList
		if (!data) {
			alert("등록된 이미지 정보가 없습니다.")
			// 		 location.href = "/Image/list"
			return;
		}

		html += "<option value=''>Select OS Platform</option>"
		for (var i in data) {
			if (data[i].connectionName == configName) {
				html += '<option value="' + data[i].id + '" >' + data[i].name + '(' + data[i].id + ')</option>';
			}
		}
		$("#ss_imageId").empty();
		$("#ss_imageId").append(html);//which OS

		//  }).catch(function(error){
		// 	console.log(error);        
		// });
	}).catch((error) => {
		console.warn(error);
		console.log(error.response)
		var errorMessage = error.response.data.error;
		commonErrorAlert(statusCode, errorMessage)
	});
}

function displayAvailableDisk() {
  
        var configName = $("#ss_regConnectionName").val()
        var url = "/setting/resources/datadisk/list"
		url +="?filterKey=connectionName&filterVal="+configName
        console.log("check disk list : ",url);
            axios.get(url).then(result=>{
            console.log("get result : ",result);
            var data = result.data.dataDiskInfoList;
            var html = "";
            console.log("get available disk : ",data);
            if(data != null || data.length > 0){
                var avDiskCnt = 0
                data.forEach(item=>{
                    //console.log("get available disk : ", item);
                    html +='<tr>'
                    + '<td class="overlay hidden column-50px" data-th="">'
                    + '<input type="checkbox" name="chk_attach" value="' + item.id + '"  title=""  /><label for="td_ch1"></label> <span class="ov off"></span>'
            
                    + '</td>'
                    + '<td class="btn_mtd ovm" data-th="name">' + item.name + '<span class="ov"></span></td>'
                    + '<td class="overlay hidden" data-th="diskType">' 
                    + item.diskType
                    +'</td>'
                    + '<td class="overlay hidden" data-th="diskSize">' 
                    + item.diskSize
                    + '/GB</td>'
                    + '</tr>';
                    avDiskCnt++;
                })
                $("#availableDiskCnt").val(avDiskCnt)
                $("#availableDiskList").empty()
                $("#availableDiskList").append(html)
            }else{
                //commonAlert("해당 VM에 Attach 가능한 DISK가 없습니다");
                addRow();
                $("#availableDiskCnt").val(0)
                return;
            }
            // if(data){
            //     if(data.length >0){
            //         data.forEach(item=>{

            //         })
            //     }
            // }

            $("#availableDiskSelectBox").modal();
            $('.dtbox.scrollbar-inner').scrollbar();
            
        })


    } 

function getVmMyiInfo(configName) {

	//var configName = $("#ss_regConnectionName option:selected").val();

	console.log("2 : ", configName);
	// getCommonVirtualMachineImageList("mcissimpleconfigure", "name"); setCommonVirtualMachineImageList()
	// var url = "/setting/resources" + "/machineimage/lookupimage";//TODO : 조회 오류남... why? connectionName으로 lookup
	var url = "/setting/resources" + "/myimage/list?filterKey=connectionName&filterVal="+configName
	var html = "";
	//  var apiInfo = 'Basic ZGVmYXVsdDpkZWZhdWx0'
	axios.get(url).then(result => {
		console.log("MyImage Info : ", result.data)
		data = result.data.myImageInfoList
		if (!data) {
			alert("등록된 My Image 정보가 없습니다.")
			// 		 location.href = "/Image/list"
			return;
		}

		html += "<option value=''>Select My Snapshot</option>"
		for (var i in data) {
			if (data[i].connectionName == configName) {
				html += '<option value="' + data[i].id + '" >' + data[i].name + '(' + data[i].id + ')</option>';
			}
		}
		$("#ss_myImageId").empty();
		$("#ss_myImageId").append(html);//which OS

		//  }).catch(function(error){
		// 	console.log(error);        
		// });
	}).catch((error) => {
		console.warn(error);
		console.log(error.response)
		var errorMessage = error.response.data.error;
		commonErrorAlert(statusCode, errorMessage)
	});
}


function getSecurityInfo(configName) {
	var configName = configName;
	if (!configName) {
		configName = $("#ss_regConnectionName option:selected").val();
	}
	//  var url = CommonURL+"/ns/"+NAMESPACE+"/resources/securityGroup";
	var url = "/setting/resources" + "/securitygroup/list"
	var html = "";
	//  var apiInfo = ApiInfo
	var default_sg = "";
	axios.get(url, {
		//  headers:{
		// 	 'Authorization': apiInfo
		//  }
	}).then(result => {
		console.log(result)
		data = result.data.SecurityGroupList
		var cnt = 0
		for (var i in data) {
			if (data[i].connectionName == configName) {
				cnt++;
				html += '<option value="' + data[i].id + '" >' + data[i].cspSecurityGroupName + '(' + data[i].id + ')</option>';
				if (cnt == 1) {
					default_sg = data[i].id
				}


			}
		}

		$("#sg").empty();
		$("#sg").append(html);// TODO : 해당 화면에 id=sg 가 없음.
		$("#s_securityGroupIds").val(default_sg)

	})
}
function getSpecInfo(configName) {
	var configName = configName;
	if (!configName) {
		configName = $("#ss_regConnectionName option:selected").val();
	}

	var url = "/setting/resources" + "/vmspec/list"
	// var url = CommonURL+"/ns/"+NAMESPACE+"/resources/spec";
	var html = "";
	// var apiInfo = ApiInfo
	axios.get(url, {
		// headers:{
		// 	'Authorization': apiInfo
		// }
	}).then(result => {
		// console.log(result.data)
		var data = result.data.VmSpecList
		console.log("spec result : ", data)
		if (data) {
			html += "<option value=''>Select Spec</option>"
			data.filter(csp => csp.connectionName === configName).map(item => (
				html += '<option value="' + item.id + '">' + item.name + '(' + item.cspSpecName + ')</option>'
			))

		} else {
			html += ""
		}


		$("#ss_spec").empty();
		$("#ss_spec").append(html);

	})
}
function getSSHKeyInfo(configName) {
	//  var configName = configName;
	if (!configName) {
		configName = $("#ss_regConnectionName option:selected").val();
	}
	//  var url = CommonURL+"/ns/"+NAMESPACE+"/resources/sshKey";
	var url = "/setting/resources" + "/sshkey/list"
	var html = "";
	//  var apiInfo = ApiInfo
	axios.get(url, {
		//  headers:{
		// 	 'Authorization': apiInfo
		//  }
	}).then(result => {
		console.log("sshKeyInfo result :", result)
		data = result.data.SshKeyList
		html += '<option value="">Select SSH Key</option>'
		for (var i in data) {
			if (data[i].connectionName == configName) {
				html += '<option value="' + data[i].id + '" >' + data[i].cspSshKeyName + '(' + data[i].id + ')</option>';
			}
		}
		$("#ss_sshKey").empty();
		$("#ss_sshKey").append(html);

	})
}

// TODO : 화면에 어디에 위치해 있는가?
function getVnetInfo(configName) {
	var configName = configName;
	console.log("get vnet INfo config name : ", configName)
	if (!configName) {
		configName = $("#ss_regConnectionName option:selected").val();
	}
	console.log("get vnet INfo config name : ", configName)
	// var url = CommonURL+"/ns/"+NAMESPACE+"/resources/vNet";
	var url = "/setting/resources/network/list";

	var html = "";
	var html2 = "";
	// var apiInfo = ApiInfo
	axios.get(url, {
		// headers:{
		//     'Authorization': apiInfo
		// }
	}).then(result => {
		data = result.data.VNetList
		// console.log("vNetwork Info : ",result)
		var init_vnet = "";
		var init_subnet = "";
		var v_net_cnt = 0
		var subnet_cnt = 0;
		for (var i in data) {
			if (data[i].connectionName == configName) {
				console.log(data[i].connectionName + " : " + configName);
				console.log(data[i]);
				html += '<option value="' + data[i].id + '" selected>' + data[i].cspVNetName + '(' + data[i].id + ')</option>';
				v_net_cnt++;
				var subnetInfoList = data[i].subnetInfoList
				//if(v_net_cnt == 1){
				init_vnet = data[i].id
				// 	console.log("init_vnet :",init_vnet)
				// }

				for (var k in subnetInfoList) {

					// init_subnet = subnetInfoList[0].IId.NameId
					init_subnet = subnetInfoList[k].id
					// console.log("init_subnet :",init_subnet)

					html2 += '<option value="' + subnetInfoList[k].id + '" >' + subnetInfoList[k].ipv4_CIDR + '</option>';
					// html2 += '<option value="'+subnetInfoList[k].IId.NameId+'" >'+subnetInfoList[k].Ipv4_CIDR+'</option>';
				}
			}
		}
		$("#vnet").empty();
		$("#vnet").append(html);
		$("#subnet").empty();
		$("#subnet").append(html2);

		//setting default
		$("#s_subnetId").val(init_subnet);
		$("#s_vNetId").val(init_vnet);
		console.log("init_vnet=" + init_vnet + ", subnet=" + init_subnet)
	})
}


const Simple_Server_Config_Arr = new Array();
var simple_data_cnt = 0
const cloneObj = obj => JSON.parse(JSON.stringify(obj))
function simpleDone_btn() {
	$("#s_regProvider").val($("#ss_regProvider").val())
	$("#s_name").val($("#ss_name").val())
	$("#s_description").val($("#ss_description").val())
	$("#s_regConnectionName").val($("#ss_regConnectionName").val())
	$("#s_spec").val($("#ss_spec").val())
	//$("#s_imageId").val($("#ss_imageId").val()) 이미지는 내가 처리
	$("#s_sshKey").val($("#ss_sshKey").val())
	$("#s_vm_cnt").val($("#ss_vm_add_cnt").val() + "")
	var select_disk = $("#ss_data_disk").val();
	
	
	console.log($("#s_imageId").val());
	console.log($("#ss_imageId").val());

	var simple_form = $("#simple_form").serializeObject()
	
	if(select_disk){
		var arr_disk = select_disk.split(",");
		simple_form.dataDiskIds = arr_disk;
	}
	console.log("simple form : ",simple_form);
	
	var server_name = simple_form.name
	var server_cnt = parseInt(simple_form.subGroupSize)
	// simple_form.subGroupSize = server_cnt
	console.log('server_cnt : ', server_cnt)
	var add_server_html = "";

	// if (server_cnt > 1) {
	// 	for (var i = 1; i <= server_cnt; i++) {
	// 		var new_vm_name = server_name + "-" + i;
	// 		var object = cloneObj(simple_form)
	// 		object.name = new_vm_name
	// 		add_server_html += '<li onclick="view_simple(\'' + simple_data_cnt + '\')">'
	// 			+ '<div class="server server_on bgbox_b">'
	// 			+ '<div class="icon"></div>'
	// 			+ '<div class="txt">' + new_vm_name + '</div>'
	// 			+ '</div>'
	// 			+ '</li>';
	// 		Simple_Server_Config_Arr.push(object)
	// 		console.log(i + "번째 Simple form data 입니다. : ", object);
	// 	}
	// } else {
	Simple_Server_Config_Arr.push(simple_form)
	displayServerCnt = ""
	if (server_cnt > 1) {
		displayServerCnt = '(' + server_cnt + ')'
	}
	add_server_html += '<li onclick="view_simple(\'' + simple_data_cnt + '\')">'
		+ '<div class="server server_on bgbox_b">'
		+ '<div class="icon"></div>'
		+ '<div class="txt">' + server_name + displayServerCnt + '</div>'
		+ '</div>'
		+ '</li>';

	// }
	$(".simple_servers_config").removeClass("active");
	console.log("add server html");
	$("#mcis_server_list").prepend(add_server_html)
	// $("#mcis_server_list").append(add_server_html)
	$("#plusVmIcon").remove();
	$("#mcis_server_list").prepend(getPlusVm());

	console.log("simple btn click and simple form data : ", simple_form)
	console.log("simple data array : ", Simple_Server_Config_Arr);
	simple_data_cnt++;
	$("#simple_form").each(function () {
		this.reset();
	})
	$("#ss_data_disk").val("");

}
function view_simple(cnt) {
	console.log('view simple cnt : ', cnt);
	var select_form_data = Simple_Server_Config_Arr[cnt]
	console.log('select_form_data : ', select_form_data);
	$(".simple_servers_config").addClass("active")
	$(".expert_servers_config").removeClass("active")
	$(".import_servers_config").removeClass("active")

}

function displayNewServerForm() {
	var $SimpleServers = $("#simpleServerConfig");
	var $ExpertServers = $("#expertServerConfig");
	var $ImportServers = $("#importServerConfig");

	var check = $(".switch .ch").is(":checked");
	console.log("check=" + check);
	if (check) {
		$SimpleServers.removeClass("active");
		$ExpertServers.addClass("active");
		$ImportServers.removeClass("active");
	} else {
		$SimpleServers.addClass("active");
		$ExpertServers.removeClass("active");
		$ImportServers.removeClass("active");
	}

	// var vmFormType = $("input[name='vmInfoType']:checked").val();
	// console.log("vmFormType = " + vmFormType)
	// if( vmFormType == "expert"){
	//     $SimpleServers.removeClass("active");
	//     $ExpertServers.addClass("active");            
	//     $ImportServers.removeClass("active");
	// }else if( vmFormType == "import"){
	//     $SimpleServers.removeClass("active");
	//     $ExpertServers.removeClass("active");            
	//     $ImportServers.addClass("active");
	// }else{// simple
	//     $SimpleServers.addClass("active");
	//     $ExpertServers.removeClass("active");            
	//     $ImportServers.removeClass("active");
	// }
}
// Expert Mode에 Import 버튼 클릭 시 해당 form display  // MCIS Create 와 VM Create의 function이름이 같음
function displayVmImportServerFormByImport() {
	var $SimpleServers = $("#simpleServerConfig");
	var $ExpertServers = $("#expertServerConfig");
	var $ImportServers = $("#importServerConfig");
	var check = $(".switch .ch").is(":checked");
	console.log("check=" + check);
	if (check) {
		$SimpleServers.removeClass("active");
		$ExpertServers.removeClass("active");
		$ImportServers.addClass("active");

		importVmInfoFromFile();// import창 띄우기 
	}
}


function importVmInfoFromFile() {
	var input = document.createElement("input");
	input.type = "file";
	// input.accept = "text/plain"; // 확장자가 xxx, yyy 일때, ".xxx, .yyy"
	input.accept = ".json";
	input.onchange = function (event) {
		importFileProcess(event.target.files[0]);
	};
	input.click();
}

// 선택한 파일을 읽어 화면에 보여줌
function importFileProcess(file) {
	try {
		var reader = new FileReader();
		reader.onload = function () {
			console.log(reader.result);
			console.log("---1")
			// $("#fileContent").val(reader.result);

			var jsonStr = JSON.stringify(reader.result)
			console.log(JSON.stringify(jsonStr));

			// 요거는 string으로만 나오네... 
			// console.log("---2")
			// var jsonObj = JSON.parse(reader.result);
			// var jsonObj = JSON.parse(jsonStr);
			// console.log(jsonObj);
			// console.log(jsonObj[0]);
			// console.log(jsonObj[10]);
			// console.log(jsonObj.name);
			// console.log("---3")

			// 요거 작동 하네.  param, value 모두 따옴표로 묶여진 json 형태여야 함.
			var newJ = $.parseJSON(reader.result);

			console.log(newJ.name);
			console.log(newJ.imageId);
			console.log(newJ.connectionName);
			console.log(newJ.securityGroupIds);
			setVmInfoToForm(newJ);
			jsonFormatter(newJ)
			//securityGroupIds: [ 	"sg-mz-aws-us-east-01"	],
		};
		//reader.readAsText(file, /* optional */ "euc-kr");
		reader.readAsText(file);
	} catch (error) {
		commonAlert("File Load Failed");
		console.log(error);
	}
}

// Provider변경시 connection 정보 filter
// function getConnectionListFilterForSimpleSelectbox(referenceObj, targetConnectionObj) {
// 	var referenceVal = $('#' + referenceObj).val();
// 	//var regionValue = region.substring(region.indexOf("]") ).trim();  
// 	// console.log(region + ", regionValue = " + regionValue);
// 	selectBoxFilterByText(targetConnectionObj, referenceVal)


// 	// $("#" + targetConnectionObj + " option:eq(0)").attr("selected", "selected");
// 	$("#es_regConnectionName").val("");
// 	setConnectionValue("");// val("")을 했을 때 자동으로 설정이 안되어서 setConnectionValue("")으로 값 set.

// }

