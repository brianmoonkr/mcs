$(document).ready(function () {
  let $summernote = $('#summernote');
  $summernote.summernote({
    lang: 'ko-kr',
    height: 300,
    placeholder: '즐겁게 작성하세요.',
    tabsize: 2
  });
  $summernote.summernote('code', $summernote.text())


  $(".save-btn").click(function (e) {
    e.preventDefault();

    var helpSeq = $("#help-form").data("seq");
    var status = $("#status-code option:selected").val();
    var question = $("#question").val();
    var answer = $("#summernote").summernote('code');

    console.log("helpSeq: ", helpSeq);


    if (status.length == 0) {
      alert("status!")
      return;
    }
    if (question.length == 0) {
      alert("question!")
      return;
    }
    if (answer.length == 0) {
      alert("answer!")
      return;
    }

    Util_AjaxPost(`/api/v1/custcenter/help/modify`, {
      help_seq: helpSeq,
      status: status,
      question: question,
      answer: answer
    }, function (res) {
      if (res.status != 200) {
        alert(res.message);
        return;
      }
      alert("저장이 완료되었습니다.");
    });
  });

  $("#cancel-btn").click(function (e) {
    e.preventDefault();
    location.href = "/admin/customer/center/help";
  });
}); // End Of jQuery