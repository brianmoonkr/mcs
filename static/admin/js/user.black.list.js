$(document).ready(function () {
  const $userListTable = $('#user-list-table');
  const $userForm = $('#user-form');
  const $categoryCode = $('#category-code');
  const $statusCode = $('#status-code');
  const $startDate = $('#start-date');
  const $endDate = $('#end-date');
  const $searchText = $('#search-text');
  const $searchBtn = $("#search-btn");
  const $searchDefaultBtn = $("#search-default-btn");

  var SearchData = {
    category: "",
    status: "",
    startDate: null,
    endDate: null,
    searchTextType: "",
    searchText: ""
  };

  // CommonCode-Category
  Util_AjaxGet("/api/v1/common/code/8", function (obj) {
    if (obj.status != 200) {
      alert(obj.message);
      return;
    }
    var codes = obj.data.Codes;
    $categoryCode.append('<option value="" selected>카테고리</option>');
    $.each(codes, function (i, v) {
      $categoryCode.append(`<option value="${v.code}">${v.code_name}</option>`);
    });
  });

  // CommonCode-Status
  Util_AjaxGet("/api/v1/common/code/9", function (obj) {
    if (obj.status != 200) {
      alert(obj.message);
      return;
    }
    var codes = obj.data.Codes;
    $statusCode.append('<option value="" selected>처리</option>');
    $.each(codes, function (i, v) {
      $statusCode.append(`<option value="${v.code}">${v.code_name}</option>`);
    });
  });

  // Datepicker
  $startDate.datepicker({
    uiLibrary: 'bootstrap4',
    format: 'yyyy-mm-dd'
  });
  $endDate.datepicker({
    uiLibrary: 'bootstrap4',
    format: 'yyyy-mm-dd'
  });


  const userListTable = $userListTable.DataTable({
    responsive: true,
    paging: true,
    lengthChange: true,
    searching: false,
    ordering: false,
    info: true,
    autoWidth: true,
    processing: true,
    serverSide: true,
    language: {
      infoFiltered: "",
      info: "Showing _START_ to _END_ of _TOTAL_ entries"
    },
    sPaginationType: "full_numbers",
    ajax: {
      url: "/api/v1/user/blacklist",
      contentType: 'application/json',
      type: "POST",
      dataSrc: "user_blacklists",
      data: function (data) {
        data.category = SearchData.category;
        data.status = SearchData.status;
        data.start_date = SearchData.startDate;
        data.end_date = SearchData.endDate;
        data.date_type = SearchData.dateType;
        data.search_text_type = SearchData.searchTextType;
        data.search_text = SearchData.searchText;
        return JSON.stringify(data);
      }
    },
    columns: [
      { data: 'row_num' },
      { data: 'nick_name' },
      { data: 'status_name' },
      { data: 'category_name' },
      { data: 'description' },
      { data: 'created_at' },
    ],
    columnDefs: [
      {
        render: function (data, type, row) {
          if (data.includes("0001-")) {
            return "-";
          } else {
            return moment(data).format('YYYY-MM-DD HH:mm:ss');
          }
        },
        targets: 5
      }
    ]
  });



  // 검색버튼
  $searchBtn.click(function (e) {
    e.preventDefault();
    SearchData.category = $("#category-code option:selected").val();
    SearchData.status = $("#status-code option:selected").val();
    if ($startDate.val() != "") {
      SearchData.startDate = moment.utc($startDate.val() + " 00:00:00", 'YYYY-MM-DD hh:mm:ss');
    } else {
      SearchData.startDate = null;
    }
    if ($endDate.val() != "") {
      SearchData.endDate = moment.utc($endDate.val() + " 23:59:59", 'YYYY-MM-DD hh:mm:ss');;
    } else {
      SearchData.endDate = null;
    }
    SearchData.searchTextType = $('#search-text-type option:selected').val();
    SearchData.searchText = $searchText.val();
    if (SearchData.searchTextType != "") {
      if (SearchData.searchText == "") {
        alert("검색어 입력!")
        return;
      }
      if (Util_RegExp(SearchData.searchText)) {
        alert("특수문자 사용금지");
        return;
      }
    }

    userListTable.ajax.reload();
  });

  // 검색초기화 버튼
  $searchDefaultBtn.click(function (e) {
    e.preventDefault();
    $userForm.each(function () {
      this.reset();
    });
  });

}); // End Of jQuery