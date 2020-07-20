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

    var category = $("#category-code option:selected").val();
    var lang = $('#lang option:selected').val();
    var status = $("#status-code option:selected").val();
    var question = $("#question").val();
    var answer = $("#summernote").summernote('code');

    //console.log("category : ", category, ", status : ", status, ", question: ", question, ", answer: ", answer);

    if (category.length == 0) {
      alert("category!")
      return;
    }
    if (lang.length == 0) {
      alert("lang!")
      return;
    }
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

    Util_AjaxPost("/api/v1/custcenter/help/save", {
      category: category,
      lang: lang,
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
}); // End Of jQuery