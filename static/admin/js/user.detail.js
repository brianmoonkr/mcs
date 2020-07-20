$(document).ready(function () {
  $("#goliveBtn").click(function (e) {
    e.preventDefault();

    window.open(`/live/${$(this).data("live")}`, '_blank');
  });

  $("#blacklistBtn").click(function (e) {
    e.preventDefault();
    Util_AjaxGet("/api/v1/common/code/8", function (obj) {
      if (obj.status != 200) {
        alert(obj.message);
        return;
      }
      var codes = obj.data.Codes;
      $("#category-code-blacklist").append('<option value="" selected>카테고리</option>');
      $.each(codes, function (i, v) {
        $("#category-code-blacklist").append(`<option value="${v.code}">${v.code_name}</option>`);
      });
    });
  });

  $("#saveBlacklistBtn").click(function (e) {
    e.preventDefault();
    var userSeq = $(this).data("userseq")
    var categoryCode = $("#category-code-blacklist option:selected").val();
    var desc = $("#blacklist-description").val();
    if (categoryCode.length == 0) {
      alert("카테고리를 선택하세요");
      return;
    }
    if (desc.length == 0) {
      alert("설명을 입력하세요");
      return;
    }
    Util_AjaxPost("/api/v1/user/blacklist/save", { user_seq: userSeq, category_code: categoryCode, desc: desc }, function (obj) {
      if (obj.status != 200) {
        alert(obj.message);
        return;
      }
      alert('저장완료');
      location.href = `/admin/user/detail/${userSeq}`;
    });
  });

  $("#cancel-blacklist-btn").click(function (e) {
    e.preventDefault();
    var userSeq = $(this).data("userseq");
    Util_AjaxPost("/api/v1/user/blacklist/cancel", { user_seq: userSeq }, function (obj) {
      if (obj.status != 200) {
        alert(obj.message);
        return;
      }
      alert('해제완료');
      location.href = `/admin/user/detail/${userSeq}`;
    });
  });

  $("#cancel-btn").click(function (e) {
    e.preventDefault();
    location.href = "/admin/user";
  });

  $("#withdrawalBtn").click(function (e) {
    e.preventDefault();
    var userID = $(this).data("userid");
    var desc = $("#withdrawal-desc").val();
    if (desc.length == 0) {
      alert("설명을 입력하세요");
      return;
    }
    Util_AjaxDelete("/api/v1/user/withdrawal", { user_id: userID, description: desc }, function (obj) {
      if (obj.status != 200) {
        alert(obj.message);
        return;
      }
      alert('탈퇴완료');
      location.href = `/admin/user`;
    });
  });




}); // End Of jQuery