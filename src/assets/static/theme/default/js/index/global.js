var state=false
$(function () {
		$('#myModal').modal({
				backdrop:false,
				keyboard:true,
				show:false
			});
		$("#addDir").click(function (){
				dire=$("#dir")
				dir=dire.val()
				if(dir==""){
					dire.parents(".control-group").addClass("error")
					dire.next(".help-inline").html("路径不能为空").show();
				}else{
					$.getJSON("add.html",{"path":dir},function (data){
							if(data.S){
								dire.next(".help-inline").html("增加成功").show();
								state=true
							}else{
								dire.parents(".control-group").addClass("error")
								dire.next(".help-inline").html(data.T).show();
							}
						});
				}
			});
		$("#closeShow").click(function (){
				if(state){
					location.reload()
				}
			});
		$("#dir").focus(function() {
				$("#dir").parents(".control-group").removeClass("error");
				$("#dir").next(".help-inline").hide().html("");
			});
		$("#multiSelect").change(function (){
				$(".alert").addClass("hide");
			});
		$("#delDir").click(function (){
				dir=$("#multiSelect").val()
				if(dir==null||dir==""){
					$(".alert").text("需要选择一个目录.");
					$(".alert").removeClass("hide");
					return
				}
				$.getJSON("del.html",{"path":dir[0]},function (data){
						if(!data.S){
								$(".alert").text("删除失败:"+data.T);
								$(".alert").removeClass("hide");
						}else{
							location.reload();
						}
					});
			});
		$("#scanCompile").click(function (){
			$.getJSON("ScanCompile.html",null,function (data){
					location.reload();
					});
			});
		$("#compile").click(function (){
		$.getJSON("Compile.html",null,function (data){
					location.reload();
					});
			});
		$("#exit").click(function (){
				$.getJSON("Close.html",null,function (data){
					location.reload();
					});
		});
		$("#compress").click(function (){
				b=false;
				if($("#compress").attr("checked")){
					b=true;
				}
				$.getJSON("Set.html",{"compress":b},function (data){
						$("#compress").attr("checked",data.S);
				});
		});
		$('.faild').popover({"placement":"bottom"});
		$('.succ').popover({"placement":"bottom","trigger":"hover"});
	})
