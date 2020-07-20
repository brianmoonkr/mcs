
function Util_AjaxGet(url, callback) {
  $.ajax({
    url: url,
    type: "GET",
    success: function (obj) {
      console.log(obj)
      callback(obj);
    }
  });
}

function Util_AjaxPut(url, data, callback) {
  $.ajax({
    url: url,
    type: "PUT",
    dataType: 'json',
    contentType: 'application/json',
    data: JSON.stringify(data),
    success: function (obj) {
      console.log(obj)
      callback(obj);
    }
  });
}

function Util_AjaxPost(url, data, callback) {
  $.ajax({
    url: url,
    type: "POST",
    dataType: 'json',
    contentType: 'application/json',
    data: JSON.stringify(data),
    success: function (obj) {
      console.log(obj)
      callback(obj);
    }
  });
}

function Util_AjaxDelete(url, data, callback) {
  $.ajax({
    url: url,
    type: "DELETE",
    dataType: 'json',
    contentType: 'application/json',
    data: JSON.stringify(data),
    success: function (obj) {
      console.log(obj)
      callback(obj);
    }
  });
}

// 특수문자 검증
// return true : 특수문자 Yes
// return false : 특수문자 No
function Util_RegExp(str) {
  var regExp = /[\{\}\[\]\/?.,;:|\)*~`!^\-_+<>@\#$%&\\\=\(\'\"]/gi
  if (regExp.test(str)) {
    return true
  }
  return false
}