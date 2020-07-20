$(document).ready(function () {
  $("#end-btn").click(function (e) {
    e.preventDefault();
    var seq = $("#live-form").data("seq");

    Util_AjaxPut(`/api/v1/live`, {
      live_seq: seq,
      Status: "1005"
    }, function (res) {
      if (res.status != 200) {
        alert(res.message);
        return;
      }
      alert("저장이 완료되었습니다.");
      location.reload();
    });

  });


  $("#cancel-btn").click(function (e) {
    e.preventDefault();
    location.href = "/admin/live";
  });
});