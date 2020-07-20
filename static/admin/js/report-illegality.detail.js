$(document).ready(function () {
  const $reportIllegalityForm = $("#report-illegality-form")

  $(".save-btn").click(function (e) {
    e.preventDefault()
    var seq = $reportIllegalityForm.data('seq');
    var handleDetail = $('#handle-detail').val();
    var btnType = $(this).data('btn');
    var status = "";

    if (handleDetail.length == 0) {
      alert("처리 내용을 입력하세요");
      return;
    }

    if (btnType == "temp") {
      status = '1001';
    } else {
      status = '1002'
    }

    var data = {
      report_illegality_seq: Number(seq),
      status: status,
      handle_detail: handleDetail
    };

    Util_AjaxPost("/api/v1/custcenter/illegality/save", data, function (res) {
      if (res.status != 200) {
        alert(res.message);
        return;
      }
      alert("저장이 완료되었습니다.");
    });

  });

  $("#cancel-btn").click(function (e) {
    e.preventDefault();
    location.href = "/admin/customer/center/report-illegality";
  });
}); // End Of jQuery