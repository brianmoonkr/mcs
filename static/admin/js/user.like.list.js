$(document).ready(function () {
  userLikeLive();
  userLikeVod();
});

function userLikeLive() {
  const $userListTable = $('#like-live-list-table');
  const $userForm = $('#like-live-form');
  const $startDate = $('#live-start-date');
  const $endDate = $('#live-end-date');
  const $searchText = $('#live-search-text');
  const $searchBtn = $("#live-search-btn");
  const $searchDefaultBtn = $("#live-search-default-btn");

  var SearchData = {
    startDate: null,
    endDate: null,
    searchTextType: "",
    searchText: ""
  };

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
      url: "/api/v1/live/like/list",
      contentType: 'application/json',
      type: "POST",
      dataSrc: "likes",
      data: function (data) {
        data.start_date = SearchData.startDate;
        data.end_date = SearchData.endDate;
        data.search_text_type = SearchData.searchTextType;
        data.search_text = SearchData.searchText;
        return JSON.stringify(data);
      }
    },
    columns: [
      { data: 'row_num' },
      { data: 'nick_name' },
      { data: 'title' },
      { data: 'created_at' },
    ],
    columnDefs: [
      {
        render: function (data, type, row) {
          return `<a href="/admin/user/detail/${row.user_seq}">${data}</a>`
        },
        targets: 1
      },
      {
        render: function (data, type, row) {
          return `<a href="#${row.live_seq}">${data}</a>`
        },
        targets: 2
      },
      {
        render: function (data, type, row) {
          if (data.includes("0001-")) {
            return "-";
          } else {
            return moment(data).format('YYYY-MM-DD HH:mm:ss');
          }
        },
        targets: 3
      }
    ]
  });



  // 검색버튼
  $searchBtn.click(function (e) {
    e.preventDefault();
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
    SearchData.searchTextType = $('#live-search-text-type option:selected').val();
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
} // End Of userLikeLive()




function userLikeVod() {
  const $userListTable = $('#like-vod-list-table');
  const $userForm = $('#like-vod-form');
  const $startDate = $('#vod-start-date');
  const $endDate = $('#vod-end-date');
  const $searchText = $('#vod-search-text');
  const $searchBtn = $("#vod-search-btn");
  const $searchDefaultBtn = $("#vod-search-default-btn");

  var SearchData = {
    startDate: null,
    endDate: null,
    searchTextType: "",
    searchText: ""
  };

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
      url: "/api/v1/vod/like/list",
      contentType: 'application/json',
      type: "POST",
      dataSrc: "vod_likes",
      data: function (data) {
        data.start_date = SearchData.startDate;
        data.end_date = SearchData.endDate;
        data.search_text_type = SearchData.searchTextType;
        data.search_text = SearchData.searchText;
        return JSON.stringify(data);
      }
    },
    columns: [
      { data: 'row_num' },
      { data: 'nick_name' },
      { data: 'title' },
      { data: 'created_at' },
    ],
    columnDefs: [
      {
        render: function (data, type, row) {
          return `<a href="/admin/user/detail/${row.user_seq}">${data}</a>`
        },
        targets: 1
      },
      {
        render: function (data, type, row) {
          return `<a href="#${row.vod_seq}">${data}</a>`
        },
        targets: 2
      },
      {
        render: function (data, type, row) {
          if (data.includes("0001-")) {
            return "-";
          } else {
            return moment(data).format('YYYY-MM-DD HH:mm:ss');
          }
        },
        targets: 3
      }
    ]
  });



  // 검색버튼
  $searchBtn.click(function (e) {
    e.preventDefault();
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
    SearchData.searchTextType = $('#live-search-text-type option:selected').val();
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

} // End Of userLikeVod()