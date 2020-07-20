$(document).ready(function () {
  const $faqForm = $("#faq-form")

  $(".save-btn").click(function (e) {
    e.preventDefault()
    var btnType = $(this).data('btn');

    var userSeq = $faqForm.data("userseq");
    var createdAt = $faqForm.data("createdat");
    var faqSeq = $("#faq-seq").text();
    var question = $("#question").text();
    var answer = $("#answer").val();
    var status = $("#status option:selected").val();
    if (status == "") {
      alert("상태를 선택해주세요")
      return;
    }
    if (answer.length == 0) {
      alert("답변을 입력해주세요");
      return;
    }

    if (btnType == "temp") {
      if (status != "1001") {
        alert("임시저장은 미처리 상태에서만 가능합니다.")
        return;
      }
    }

    var data = {
      faq_seq: Number(faqSeq),
      user_seq: Number(userSeq),
      status: status,
      question: question,
      answer: answer
    };

    Util_AjaxPost("/api/v1/custcenter/faq/save", data, function (res) {
      if (res.status != 200) {
        alert(res.message);
        return;
      }
      alert("저장이 완료되었습니다.");
    });
  });

  $("#cancel-btn").click(function (e) {
    e.preventDefault();
    location.href = "/admin/customer/center/faq";
  });

});